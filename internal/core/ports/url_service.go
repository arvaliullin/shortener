package ports

import (
	"context"

	"github.com/arvaliullin/shortener/internal/core/domain"
)

//go:generate mockgen -source=url_service.go -destination=mocks/url_service_mock.go -package=mocks

type URLService interface {
	Shorten(ctx context.Context, originalURL string) (id string, err error)
	Resolve(ctx context.Context, id string) (originalURL string, err error)
	// List возвращает сохраненные короткие ссылки.
	List(ctx context.Context) ([]domain.ShortURL, error)
	// Update меняет исходный адрес короткой ссылки.
	Update(ctx context.Context, id, originalURL string) error
	// Delete удаляет короткую ссылку.
	Delete(ctx context.Context, id string) error
}
