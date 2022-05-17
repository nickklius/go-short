package storages

import (
	"sync"

	"github.com/nickklius/go-short/internal/config"
	"github.com/nickklius/go-short/internal/utils"
)

type MemoryStorage struct {
	mux  sync.Mutex
	conf config.Config
	data map[string]string
}

func NewMemoryStorage(c config.Config) Repository {
	return &MemoryStorage{
		conf: c,
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

func (s *MemoryStorage) Create(longURL string) (string, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	for {
		short := utils.GenerateKey(s.conf.Letters, s.conf.KeyLength)
		if _, ok := s.data[short]; !ok {
			s.data[short] = longURL
			return short, nil
		}
	}
}

func (s *MemoryStorage) GetAll() *map[string]string {
	s.mux.Lock()
	defer s.mux.Unlock()
	return &s.data
}

func (s *MemoryStorage) Flush() error {
	return nil
}
