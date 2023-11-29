package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rshelekhov/url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(dsn string) (*Storage, error) {
	const fn = "storage.postgres.New"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const fn = "storage.postgres.SaveURL"

	var lastInsertId int64
	err := s.db.QueryRow("INSERT INTO url(url, alias) VALUES($1, $2) RETURNING id", urlToSave, alias).Scan(&lastInsertId)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, storage.ErrURLExists)
	}

	return lastInsertId, nil
}
