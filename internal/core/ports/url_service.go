package ports

import "context"

//go:generate mockgen -source=url_service.go -destination=mocks/url_service_mock.go -package=mocks

type URLService interface {
	Shorten(ctx context.Context, originalURL string) (id string, err error)
	Resolve(ctx context.Context, id string) (originalURL string, err error)
}
