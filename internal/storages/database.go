package storages

import (
	"database/sql"
	"sync"

	_ "github.com/jackc/pgx/stdlib"
)

type DatabaseStorage struct {
	mux  sync.Mutex
	conn *sql.DB
}

func NewDatabaseStorage(dsn string) (*DatabaseStorage, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, ErrDBConnNotEstablished
	}

	return &DatabaseStorage{
		conn: db,
	}, nil
}

func (s *DatabaseStorage) Read(shortURL string) (string, error) {
	return "", nil
}

func (s *DatabaseStorage) Create(shortURL, longURL, userID string) error {
	return nil
}

func (s *DatabaseStorage) GetAll() map[string]URLEntry {
	return make(map[string]URLEntry)
}

func (s *DatabaseStorage) GetAllByUserID(userID string) map[string]URLEntry {
	return make(map[string]URLEntry)
}

func (s *DatabaseStorage) Ping() error {
	err := s.conn.Ping()
	if err != nil {
		return err
	}
	return nil
}
