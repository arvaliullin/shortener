package dto

// ListInput не содержит полей: инструмент читает все сохраненные ссылки.
type ListInput struct{}

// ListOutput возвращает сохраненные короткие ссылки.
type ListOutput struct {
	URLs []URLItem `json:"urls" jsonschema:"сохраненные короткие ссылки"`
}
