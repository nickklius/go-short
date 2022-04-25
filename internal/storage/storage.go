package storage

import (
	"context"
	"github.com/nickklius/go-short/internal/utils"
)

type Repository interface {
	Create(ctx context.Context, longURL string) (string, error)
	Read(ctx context.Context, shortURL string) (string, error)
}

type MapURLStorage struct {
	Storage map[string]string
}

func (s *MapURLStorage) Read(_ context.Context, shorURL string) (string, error) {
	return s.Storage[shorURL], nil
}

func (s *MapURLStorage) Create(_ context.Context, longURL string) (string, error) {
	for {
		short := utils.GenerateKey()
		if _, ok := s.Storage[short]; !ok {
			s.Storage[short] = longURL
			return short, nil
		}
	}
}

func CreateShortURL(r Repository, ctx context.Context, longURL string) (string, error) {
	return r.Create(ctx, longURL)
}

func RetrieveURL(r Repository, ctx context.Context, shortURL string) (string, error) {
	return r.Read(ctx, shortURL)
}
