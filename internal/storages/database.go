package storages

import (
	"context"
	"database/sql"
	"sync"

	_ "github.com/jackc/pgx/stdlib"
)

type DatabaseStorage struct {
	mux  sync.Mutex
	conn *sql.DB
}

func NewDatabaseStorage(ctx context.Context, dsn string) (*DatabaseStorage, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, ErrDBConnNotEstablished
	}

	s := &DatabaseStorage{
		conn: db,
	}

	err = s.createTables(ctx)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *DatabaseStorage) Read(ctx context.Context, shortURL string) (string, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	var longURL string

	readQuery := "SELECT url FROM urls WHERE short = $1"
	row := s.conn.QueryRowContext(ctx, readQuery, shortURL)

	err := row.Scan(&longURL)
	if err != nil {
		return "", err
	}

	if longURL == "" {
		return "", ErrNotFound
	}

	return longURL, nil
}

func (s *DatabaseStorage) Create(ctx context.Context, shortURL, longURL, userID string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	createQuery := "INSERT INTO urls (user_id, short, url) VALUES($1, $2, $3)"
	_, err := s.conn.ExecContext(ctx, createQuery, userID, shortURL, longURL)

	if err != nil {
		return err
	}

	return nil
}

func (s *DatabaseStorage) GetAll() (map[string]URLEntry, error) {
	return nil, ErrMethodNotImplemented
}

func (s *DatabaseStorage) GetAllByUserID(ctx context.Context, userID string) (map[string]string, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	type result struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	userURLs := make(map[string]string)

	getQuery := "SELECT short, url FROM urls WHERE user_id=$1"
	rows, err := s.conn.QueryContext(ctx, getQuery, userID)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()

	for rows.Next() {
		var u result

		err = rows.Scan(&u.ShortURL, &u.OriginalURL)
		if err != nil {
			return userURLs, err
		}

		userURLs[u.ShortURL] = u.OriginalURL
	}

	return userURLs, nil
}

func (s *DatabaseStorage) Ping() error {
	return s.conn.Ping()
}

func (s *DatabaseStorage) createTables(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS urls (
		id bigserial PRIMARY KEY,
		user_id text not null,
		short text not null UNIQUE,
		url text not null 
	);`

	_, err := s.conn.ExecContext(ctx, query)
	return err
}
