package storages

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"sync"
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

func NewLocalStorage(ctx context.Context, p string, closeServiceCh chan struct{}) (Repository, error) {
	s := NewMemoryStorage(ctx, closeServiceCh)

	f, err := NewFileHandler(p)
	if err != nil {
		return s, err
	}
	defer f.Close()

	err = f.Read(ctx, s)
	if (err != io.EOF) && (err != nil) {
		return nil, err
	}

	return &LocalStorage{
		fileName: p,
		data:     s,
	}, nil
}

func (s *LocalStorage) Read(ctx context.Context, shortURL string) (string, error) {
	return s.data.Read(ctx, shortURL)
}

func (s *LocalStorage) Create(ctx context.Context, shortURL, longURL, userID string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	f, err := NewFileHandler(s.fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	err = s.data.Create(ctx, shortURL, longURL, userID)
	if err != nil {
		return err
	}

	err = f.Save(s)
	if err != nil {
		return err
	}
	return nil
}

func (s *LocalStorage) GetAll() (map[string]URLEntry, error) {
	return s.data.GetAll()
}

func (s *LocalStorage) GetAllByUserID(ctx context.Context, userID string) (map[string]string, error) {
	return s.data.GetAllByUserID(ctx, userID)
}

func (s *LocalStorage) Ping() error {
	return ErrMethodNotImplemented
}

func (s *LocalStorage) UpdateURLInBatchMode(_ context.Context, _ string, _ []string) {}

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
	data, err := s.GetAll()
	if err != nil {
		return err
	}

	err = f.encoder.Encode(data)
	if err != nil {
		return err
	}
	return f.file.Sync()
}

func (f *fileHandler) Read(ctx context.Context, s Repository) error {
	m := make(map[string]URLEntry)
	err := f.decoder.Decode(&m)
	if err != nil {
		return err
	}

	for k, v := range m {
		err = s.Create(ctx, k, v.URL, v.UserID)
		if err != nil {
			return err
		}
	}

	return nil
}
