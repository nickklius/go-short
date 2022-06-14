package storages

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

type DatabaseStorage struct {
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
	readQuery := `SELECT url FROM urls WHERE short = $1`
	row := s.conn.QueryRowContext(ctx, readQuery, shortURL)

	var longURL string

	err := row.Scan(&longURL)
	if err == sql.ErrNoRows {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}

	return longURL, nil
}

func (s *DatabaseStorage) Create(ctx context.Context, shortURL, longURL, userID string) error {
	checkShortViolation := shortURL

	createQuery := `INSERT INTO urls 
						(user_id, short, url) 	
						VALUES($1, $2, $3) 
					ON CONFLICT (url) DO UPDATE SET 
					    url = $3
					RETURNING short`
	err := s.conn.QueryRowContext(ctx, createQuery, userID, shortURL, longURL).Scan(&checkShortViolation)
	if err != nil {
		return err
	}

	if checkShortViolation != shortURL {
		return NewInsertURLUniqError(checkShortViolation, errors.New("duplicate url"))
	}

	return nil
}

func (s *DatabaseStorage) GetAll() (map[string]URLEntry, error) {
	return nil, ErrMethodNotImplemented
}

func (s *DatabaseStorage) UpdateURLInBatchMode(ctx context.Context, urls []string, userID string) error {
	tx, err := s.conn.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `UPDATE urls SET deleted = true WHERE user_id = $1
									AND short = $2`)
	if err != nil {
		return err
	}

	for _, u := range urls {
		if _, err = stmt.Exec(userID, u); err != nil {
			if err = tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *DatabaseStorage) GetAllByUserID(ctx context.Context, userID string) (map[string]string, error) {
	type result struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	userURLs := make(map[string]string)

	getQuery := `SELECT short, url FROM urls WHERE user_id=$1`
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
		url text not null UNIQUE,
		deleted boolean not null default false              
	);`

	_, err := s.conn.ExecContext(ctx, query)
	return err
}

type InsertURLUniqError struct {
	ShortURL string
	Err      error
}

func (e *InsertURLUniqError) Error() string {
	return fmt.Sprintf("%v: %v", e.ShortURL, e.Err)
}

func (e *InsertURLUniqError) Unwrap() error {
	return e.Err
}

func NewInsertURLUniqError(shortURL string, err error) error {
	return &InsertURLUniqError{
		ShortURL: shortURL,
		Err:      err,
	}
}
