package storages

import (
	"sync"
)

type URLEntry struct {
	URL    string `json:"url"`
	UserID string `json:"user_id"`
}

type MemoryStorage struct {
	mux  sync.Mutex
	data map[string]URLEntry
}

func NewMemoryStorage() Repository {
	return &MemoryStorage{
		data: make(map[string]URLEntry),
	}
}

func (s *MemoryStorage) Read(shortURL string) (string, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if v, ok := s.data[shortURL]; ok {
		return v.URL, nil
	}

	return "", ErrNotFound
}

func (s *MemoryStorage) Create(shortURL, longURL, userID string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if _, ok := s.data[shortURL]; ok {
		return ErrAlreadyExists
	}
	s.data[shortURL] = URLEntry{longURL, userID}

	return nil
}

func (s *MemoryStorage) GetAll() map[string]URLEntry {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.data
}

func (s *MemoryStorage) GetAllByUserID(userID string) map[string]URLEntry {
	s.mux.Lock()
	defer s.mux.Unlock()

	userURLs := make(map[string]URLEntry)

	for short, url := range s.data {
		if url.UserID == userID {
			userURLs[short] = url
		}
	}

	return userURLs
}

func (s *MemoryStorage) Ping() error {
	return ErrMethodNotImplemented
}
