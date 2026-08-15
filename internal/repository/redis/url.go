package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
)

var ErrURLNotFound = fmt.Errorf("URL not found")

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

func (r *URLRepository) Save(ctx context.Context, id, originalURL string) error {
	return r.client.Set(ctx, id, originalURL, 0).Err()
}

func (r *URLRepository) Find(ctx context.Context, id string) (originalURL string, err error) {
	originalURL, err = r.client.Get(ctx, id).Result()
	if err == goredis.Nil {
		return "", ErrURLNotFound
	}
	if err != nil {
		return "", err
	}
	return originalURL, nil
}
