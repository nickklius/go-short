package storages

import (
	"sync"
)

type MemoryStorage struct {
	mux  sync.Mutex
	data map[string]string
}

func NewMemoryStorage() Repository {
	return &MemoryStorage{
		data: make(map[string]string),
	}
}

func (s *MemoryStorage) Read(shortURL string) (string, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if v, ok := s.data[shortURL]; ok {
		return v, nil
	}

	return "", ErrNotFound
}

func (s *MemoryStorage) Create(shortURL, longURL string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if _, ok := s.data[shortURL]; ok {
		return ErrAlreadyExists
	}
	s.data[shortURL] = longURL

	return nil
}

func (s *MemoryStorage) GetAll() map[string]string {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.data
}
