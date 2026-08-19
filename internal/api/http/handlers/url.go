package handlers

import (
	"github.com/arvaliullin/shortener/internal/core/ports"
)

type URLHandler struct {
	urlService ports.URLService
	baseURL    string
}

func NewURLHandler(urlService ports.URLService, baseURL string) *URLHandler {
	return &URLHandler{
		urlService: urlService,
		baseURL:    baseURL,
	}
}
