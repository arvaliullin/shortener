package services

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/arvaliullin/shortener/internal/core/ports"
)

const (
	idLength  = 8
	idCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type URLService struct {
	repository ports.URLRepository
}

func NewURLService(repository ports.URLRepository) *URLService {
	return &URLService{
		repository: repository,
	}
}

func (service *URLService) Shorten(ctx context.Context, originalURL string) (string, error) {
	id, err := generateID(idLength)
	if err != nil {
		return "", err
	}

	if err := service.repository.Save(ctx, id, originalURL); err != nil {
		return "", err
	}

	return id, nil
}

func (service *URLService) Resolve(ctx context.Context, id string) (originalURL string, err error) {
	return service.repository.Find(ctx, id)
}

func generateID(length int) (string, error) {
	result := make([]byte, length)
	max := big.NewInt(int64(len(idCharset)))

	for i := range result {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		result[i] = idCharset[n.Int64()]
	}

	return string(result), nil
}
