package ports

import "context"

//go:generate mockgen -source=url_repo.go -destination=mocks/url_repo_mock.go -package=mocks

type URLRepository interface {
	Save(ctx context.Context, id, originalURL string) error
	Find(ctx context.Context, id string) (originalURL string, err error)
}
