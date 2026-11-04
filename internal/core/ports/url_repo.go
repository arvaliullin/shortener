package ports

import (
	"context"

	"github.com/arvaliullin/shortener/internal/core/domain"
)

//go:generate mockgen -source=url_repo.go -destination=mocks/url_repo_mock.go -package=mocks

type URLRepository interface {
	Save(ctx context.Context, id, originalURL string) error
	Find(ctx context.Context, id string) (originalURL string, err error)
	// FindByOriginal возвращает идентификатор уже сохраненного адреса.
	FindByOriginal(ctx context.Context, originalURL string) (id string, err error)
	// List возвращает сохраненные ссылки от новых к старым.
	List(ctx context.Context) ([]domain.ShortURL, error)
	// Update меняет исходный адрес существующей ссылки.
	Update(ctx context.Context, id, originalURL string) error
	// Delete удаляет ссылку.
	Delete(ctx context.Context, id string) error
}
