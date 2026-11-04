package dto

// UpdateInput задает новый исходный адрес существующей короткой ссылки.
type UpdateInput struct {
	ID  string `json:"id" jsonschema:"идентификатор короткой ссылки"`
	URL string `json:"url" jsonschema:"новый исходный адрес"`
}

// URLID возвращает идентификатор без крайних пробелов.
func (in UpdateInput) URLID() (string, error) {
	return requiredID(in.ID)
}

// OriginalURL возвращает новый исходный адрес без крайних пробелов.
func (in UpdateInput) OriginalURL() (string, error) {
	return requiredURL(in.URL)
}
