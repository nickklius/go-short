package storage

import "github.com/nickklius/go-short/internal/utils"

type Repository interface {
	Read(shortURL string) string
	Create(longURL string) string
}

type MapURLStorage struct {
	Storage map[string]string
}

func (s *MapURLStorage) Read(shorURL string) string {
	return s.Storage[shorURL]
}

func (s *MapURLStorage) Create(longURL string) string {
	for {
		short := utils.GenerateKey()
		if _, ok := s.Storage[short]; !ok {
			s.Storage[short] = longURL
			return short
		}
	}
}

func CreateShortURL(r Repository, longURL string) string {
	return r.Create(longURL)
}

func RetrieveURL(r Repository, shortURL string) string {
	return r.Read(shortURL)
}
