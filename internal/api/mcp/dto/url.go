// Package dto описывает вход и выход инструментов MCP для коротких ссылок.
package dto

import (
	"errors"
	"strings"
)

// IDInput задает идентификатор существующей короткой ссылки.
type IDInput struct {
	ID string `json:"id" jsonschema:"идентификатор короткой ссылки"`
}

// URLID возвращает идентификатор без крайних пробелов.
func (in IDInput) URLID() (string, error) {
	return requiredID(in.ID)
}

// URLItem описывает сохраненную короткую ссылку.
type URLItem struct {
	ID     string `json:"id" jsonschema:"идентификатор короткой ссылки"`
	URL    string `json:"url" jsonschema:"исходный адрес"`
	Result string `json:"result" jsonschema:"короткая ссылка"`
}

// NewURLItem собирает описание ссылки из базового URL, идентификатора и исходного адреса.
func NewURLItem(baseURL, id, originalURL string) URLItem {
	return URLItem{
		ID:     id,
		URL:    originalURL,
		Result: shortURL(baseURL, id),
	}
}

// shortURL собирает короткую ссылку так же, как HTTP API.
func shortURL(baseURL, id string) string {
	return strings.TrimRight(baseURL, "/") + "/" + id
}

// requiredURL отклоняет пустой исходный адрес.
func requiredURL(raw string) (string, error) {
	originalURL := strings.TrimSpace(raw)
	if originalURL == "" {
		return "", errors.New("url is required")
	}
	return originalURL, nil
}

// requiredID отклоняет пустой идентификатор.
func requiredID(raw string) (string, error) {
	id := strings.TrimSpace(raw)
	if id == "" {
		return "", errors.New("id is required")
	}
	return id, nil
}
