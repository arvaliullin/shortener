package dto

// ShortenInput задает исходный адрес для новой короткой ссылки.
type ShortenInput struct {
	URL string `json:"url" jsonschema:"исходный адрес"`
}

// OriginalURL возвращает исходный адрес без крайних пробелов.
func (in ShortenInput) OriginalURL() (string, error) {
	return requiredURL(in.URL)
}

// ShortenOutput возвращает идентификатор и полную короткую ссылку.
type ShortenOutput struct {
	ID     string `json:"id" jsonschema:"идентификатор короткой ссылки"`
	Result string `json:"result" jsonschema:"короткая ссылка"`
}

// NewShortenOutput собирает ответ из базового URL и идентификатора.
func NewShortenOutput(baseURL, id string) ShortenOutput {
	return ShortenOutput{
		ID:     id,
		Result: shortURL(baseURL, id),
	}
}
