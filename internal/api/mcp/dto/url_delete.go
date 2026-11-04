package dto

// DeleteOutput подтверждает удаление короткой ссылки.
type DeleteOutput struct {
	ID string `json:"id" jsonschema:"идентификатор удаленной короткой ссылки"`
}
