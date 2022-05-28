package storages

import (
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

func NewLocalStorage(p string) (Repository, error) {
	s := NewMemoryStorage()

	f, err := NewFileHandler(p)
	if err != nil {
		return s, err
	}
	defer f.Close()

	err = f.Read(s)
	if (err != io.EOF) && (err != nil) {
		return nil, err
	}

	return &LocalStorage{
		fileName: p,
		data:     s,
	}, nil
}

func (s *LocalStorage) Read(shortURL string) (string, error) {
	return s.data.Read(shortURL)
}

func (s *LocalStorage) Create(shortURL, longURL, userID string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	f, err := NewFileHandler(s.fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	err = s.data.Create(shortURL, longURL, userID)
	if err != nil {
		return err
	}

	err = f.Save(s)
	if err != nil {
		return err
	}
	return nil
}

func (s *LocalStorage) GetAll() map[string]URLEntry {
	return s.data.GetAll()
}

func (s *LocalStorage) GetAllByUserID(userID string) map[string]URLEntry {
	return s.data.GetAllByUserID(userID)
}

func (s *LocalStorage) Ping() error {
	return ErrMethodNotImplemented
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
	return f.file.Sync()
}

func (f *fileHandler) Read(s Repository) error {
	m := make(map[string]URLEntry)
	err := f.decoder.Decode(&m)
	if err != nil {
		return err
	}

	for k, v := range m {
		err = s.Create(k, v.URL, v.UserID)
		if err != nil {
			return err
		}
	}

	return nil
}
