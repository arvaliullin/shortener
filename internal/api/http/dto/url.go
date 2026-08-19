package dto

import "strings"

// ShortenRequest представляет запрос на сокращение URL.
type ShortenRequest struct {
	URL string `json:"url"`
}

// IsValid проверяет корректность данных запроса.
func (r *ShortenRequest) IsValid() bool {
	return strings.TrimSpace(r.URL) != ""
}

// ShortenResponse представляет ответ с короткой ссылкой.
type ShortenResponse struct {
	Result string `json:"result"`
}

// NewShortenResponse собирает DTO ответа из базового URL и идентификатора.
func NewShortenResponse(baseURL, id string) ShortenResponse {
	return ShortenResponse{Result: strings.TrimRight(baseURL, "/") + "/" + id}
}
