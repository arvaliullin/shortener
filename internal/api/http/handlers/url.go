package handlers

import (
	"github.com/arvaliullin/shortener/internal/core/ports"
)

type URLHandler struct {
	urlService ports.URLService
}

func NewURLHandler(urlService ports.URLService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
	}
}
