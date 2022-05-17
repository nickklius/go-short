package storages

import "errors"

var (
	ErrNotFound = errors.New("not found")
)

type Repository interface {
	Create(longURL string) (string, error)
	Read(shortURL string) (string, error)
	GetAll() *map[string]string // без указателя разваливается реализация
	Flush() error
}
