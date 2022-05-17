package storages

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/nickklius/go-short/internal/config"
)

type LocalStorage struct {
	mux      sync.Mutex
	fileName string
	data     Repository
}

type fileHandler struct {
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

func NewLocalStorage(c config.Config) (Repository, error) {
	s := NewMemoryStorage(c)

	f, err := NewFileHandler(c.FileStoragePath)
	if err != nil {
		return s, err
	}
	defer f.Close()

	f.Read(s)

	return &LocalStorage{
		fileName: c.FileStoragePath,
		data:     s,
	}, nil
}

func (s *LocalStorage) Read(shortURL string) (string, error) {
	return s.data.Read(shortURL)
}

func (s *LocalStorage) Create(longURL string) (string, error) {
	shortURL, err := s.data.Create(longURL)
	return shortURL, err
}

func (s *LocalStorage) GetAll() *map[string]string {
	return s.data.GetAll()
}

func (s *LocalStorage) Flush() error {
	s.mux.Lock()
	defer s.mux.Unlock()

	f, err := NewFileHandler(s.fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	err = f.Save(s)
	if err != nil {
		return err
	}

	return nil
}

func NewFileHandler(fileName string) (*fileHandler, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &fileHandler{
		file:    file,
		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
	}, nil
}

func (f *fileHandler) Close() {
	f.file.Close()
}

func (f *fileHandler) Save(s Repository) error {
	err := f.encoder.Encode(s.GetAll())
	if err != nil {
		return err
	}
	return nil
}

func (f *fileHandler) Read(s Repository) error {
	err := f.decoder.Decode(s.GetAll())
	if err != nil {
		return err
	}
	return nil
}
