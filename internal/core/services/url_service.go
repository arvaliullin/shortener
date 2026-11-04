package services

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/arvaliullin/shortener/internal/core/domain"
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

// Shorten возвращает идентификатор уже сохраненного адреса или создает новый.
func (service *URLService) Shorten(ctx context.Context, originalURL string) (string, error) {
	id, err := service.repository.FindByOriginal(ctx, originalURL)
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, domain.ErrURLNotFound) {
		return "", err
	}

	id, err = generateID(idLength)
	if err != nil {
		return "", err
	}

	if err := service.repository.Save(ctx, id, originalURL); err != nil {
		return "", err
	}

	return id, nil
}

// Resolve ищет оригинальный URL в Redis по ID
func (service *URLService) Resolve(ctx context.Context, id string) (originalURL string, err error) {
	return service.repository.Find(ctx, id)
}

// List возвращает сохраненные короткие ссылки.
func (service *URLService) List(ctx context.Context) ([]domain.ShortURL, error) {
	return service.repository.List(ctx)
}

// Update меняет исходный адрес короткой ссылки.
func (service *URLService) Update(ctx context.Context, id, originalURL string) error {
	return service.repository.Update(ctx, id, originalURL)
}

// Delete удаляет короткую ссылку.
func (service *URLService) Delete(ctx context.Context, id string) error {
	return service.repository.Delete(ctx, id)
}

// generateID генерирует случайный ID для оригинального URL
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
