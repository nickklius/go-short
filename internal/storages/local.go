package storages

import (
	"encoding/json"
	"github.com/nickklius/go-short/internal/config"
	"os"
)

type LocalStorage struct {
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
	f, err := NewFileHandler(s.fileName)
	if err != nil {
		return "", err
	}

	shortURL, err := s.data.Create(longURL)

	f.Save(s)

	return shortURL, err
}

func (s *LocalStorage) GetAll() *map[string]string {
	return s.data.GetAll()
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

func (f *fileHandler) Save(s Repository) {
	f.encoder.Encode(s.GetAll())
	f.Close()
}

func (f *fileHandler) Read(s Repository) {
	f.decoder.Decode(s.GetAll())
	f.Close()
}
