package domain

import "errors"

// ErrURLNotFound возвращается, когда короткой ссылки нет в хранилище.
var ErrURLNotFound = errors.New("URL not found")

// ShortURL хранит идентификатор короткой ссылки и исходный адрес.
type ShortURL struct {
	ID       string
	Original string
}
