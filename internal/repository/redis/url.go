package redis

import (
	"context"
	"time"

	"github.com/arvaliullin/shortener/internal/core/domain"
	goredis "github.com/redis/go-redis/v9"
)

const (
	idsKey   = "shortener:ids"
	byURLKey = "shortener:by-url"
)

type URLRepository struct {
	client *goredis.Client
}

func NewURLRepository(addr string) *URLRepository {
	return &URLRepository{
		client: goredis.NewClient(&goredis.Options{
			Addr: addr,
		}),
	}
}

// Save сохраняет оригинальный URL и добавляет идентификатор в индекс.
func (r *URLRepository) Save(ctx context.Context, id, originalURL string) error {
	_, err := r.client.TxPipelined(ctx, func(pipe goredis.Pipeliner) error {
		pipe.Set(ctx, id, originalURL, 0)
		pipe.HSet(ctx, byURLKey, originalURL, id)
		pipe.ZAdd(ctx, idsKey, goredis.Z{
			Score:  float64(time.Now().UnixMilli()),
			Member: id,
		})
		return nil
	})
	return err
}

// FindByOriginal возвращает идентификатор уже сохраненного адреса.
func (r *URLRepository) FindByOriginal(ctx context.Context, originalURL string) (string, error) {
	id, err := r.client.HGet(ctx, byURLKey, originalURL).Result()
	if err == nil {
		return id, nil
	}
	if err != goredis.Nil {
		return "", err
	}

	items, err := r.List(ctx)
	if err != nil {
		return "", err
	}
	for _, item := range items {
		if item.Original != originalURL {
			continue
		}
		if err := r.client.HSet(ctx, byURLKey, originalURL, item.ID).Err(); err != nil {
			return "", err
		}
		return item.ID, nil
	}
	return "", domain.ErrURLNotFound
}

// Find ищет оригинальный URL в Redis по ID.
func (r *URLRepository) Find(ctx context.Context, id string) (originalURL string, err error) {
	originalURL, err = r.client.Get(ctx, id).Result()
	if err == goredis.Nil {
		return "", domain.ErrURLNotFound
	}
	if err != nil {
		return "", err
	}
	return originalURL, nil
}

// List возвращает сохраненные ссылки от новых к старым.
func (r *URLRepository) List(ctx context.Context) ([]domain.ShortURL, error) {
	ids, err := r.client.ZRangeArgs(ctx, goredis.ZRangeArgs{
		Key:   idsKey,
		Start: 0,
		Stop:  -1,
		Rev:   true,
	}).Result()
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return []domain.ShortURL{}, nil
	}

	values, err := r.client.MGet(ctx, ids...).Result()
	if err != nil {
		return nil, err
	}

	items := make([]domain.ShortURL, 0, len(ids))
	for i, id := range ids {
		original, ok := values[i].(string)
		if !ok || original == "" {
			continue
		}
		items = append(items, domain.ShortURL{ID: id, Original: original})
	}
	return items, nil
}

// Update меняет исходный адрес существующей ссылки.
func (r *URLRepository) Update(ctx context.Context, id, originalURL string) error {
	current, err := r.Find(ctx, id)
	if err != nil {
		return err
	}
	if current == originalURL {
		return nil
	}

	ownsIndex, err := r.ownsURL(ctx, current, id)
	if err != nil {
		return err
	}

	_, err = r.client.TxPipelined(ctx, func(pipe goredis.Pipeliner) error {
		pipe.Set(ctx, id, originalURL, 0)
		if ownsIndex {
			pipe.HDel(ctx, byURLKey, current)
		}
		pipe.HSet(ctx, byURLKey, originalURL, id)
		return nil
	})
	return err
}

// Delete удаляет ссылку и ее идентификатор из индекса.
func (r *URLRepository) Delete(ctx context.Context, id string) error {
	originalURL, err := r.Find(ctx, id)
	if err != nil {
		return err
	}

	ownsIndex, err := r.ownsURL(ctx, originalURL, id)
	if err != nil {
		return err
	}

	_, err = r.client.TxPipelined(ctx, func(pipe goredis.Pipeliner) error {
		pipe.Del(ctx, id)
		pipe.ZRem(ctx, idsKey, id)
		if ownsIndex {
			pipe.HDel(ctx, byURLKey, originalURL)
		}
		return nil
	})
	return err
}

// ownsURL сообщает, что обратный индекс адреса указывает на этот идентификатор.
func (r *URLRepository) ownsURL(ctx context.Context, originalURL, id string) (bool, error) {
	owner, err := r.client.HGet(ctx, byURLKey, originalURL).Result()
	if err == goredis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return owner == id, nil
}
