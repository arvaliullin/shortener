package redis

import (
	"context"
	"time"

	"github.com/arvaliullin/shortener/internal/core/domain"
	goredis "github.com/redis/go-redis/v9"
)

const idsKey = "shortener:ids"

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
		pipe.ZAdd(ctx, idsKey, goredis.Z{
			Score:  float64(time.Now().UnixMilli()),
			Member: id,
		})
		return nil
	})
	return err
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
	n, err := r.client.Exists(ctx, id).Result()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrURLNotFound
	}
	return r.client.Set(ctx, id, originalURL, 0).Err()
}

// Delete удаляет ссылку и ее идентификатор из индекса.
func (r *URLRepository) Delete(ctx context.Context, id string) error {
	deleted, err := r.client.Del(ctx, id).Result()
	if err != nil {
		return err
	}
	if deleted == 0 {
		return domain.ErrURLNotFound
	}
	return r.client.ZRem(ctx, idsKey, id).Err()
}
