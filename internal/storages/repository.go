package storages

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("key already exist")
)

type Repository interface {
	Create(shortURL, longURL string) error
	Read(shortURL string) (string, error)
	GetAll() map[string]string
}
