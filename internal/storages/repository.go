package storages

import (
	"context"
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
	Create(ctx context.Context, shortURL, longURL, userID string) error
	Read(ctx context.Context, shortURL string) (string, error)
	GetAll() (map[string]URLEntry, error)
	GetAllByUserID(ctx context.Context, userID string) (map[string]string, error)
	UpdateURLInBatchMode(ctx context.Context, urls []string, userID string) error
	Ping() error
}
