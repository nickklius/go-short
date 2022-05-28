package storages

import (
	"errors"
)

var (
	ErrNotFound               = errors.New("not found")
	ErrAlreadyExists          = errors.New("key already exist")
	ErrDBConnNotEstablished   = errors.New("couldn't create DB connection")
	ErrLocalStorageNotCreated = errors.New("couldn't create local storage file")
	ErrMethodNotImplemented   = errors.New("method not implemented")
)

type Repository interface {
	Create(shortURL, longURL, userID string) error
	Read(shortURL string) (string, error)
	GetAll() map[string]URLEntry
	GetAllByUserID(userID string) map[string]URLEntry
	Ping() error
}
