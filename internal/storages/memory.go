package storages

import (
	"github.com/nickklius/go-short/internal/utils"
)

type MemoryStorage struct {
	data map[string]string
}

func NewMemoryStorage() Repository {
	return &MemoryStorage{
		data: make(map[string]string),
	}
}

func (s *MemoryStorage) Read(shortURL string) (string, error) {
	return s.data[shortURL], nil
}

func (s *MemoryStorage) Create(longURL string) (string, error) {
	for {
		short := utils.GenerateKey()
		if _, ok := s.data[short]; !ok {
			s.data[short] = longURL
			return short, nil
		}
	}
}

func (s *MemoryStorage) GetAll() *map[string]string {
	return &s.data
}
