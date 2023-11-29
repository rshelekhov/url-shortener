package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
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

// Close method
func (s *Storage) Close() error {
	return s.db.Close()
}

// SaveURL method
func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const fn = "storage.postgres.SaveURL"

	sqlStatement := `INSERT INTO url(url, alias) VALUES($1, $2) RETURNING id`
	var lastInsertId int64
	err := s.db.QueryRow(sqlStatement, urlToSave, alias).Scan(&lastInsertId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == storage.UniqueConstraintViolation {
				return 0, fmt.Errorf("%s: %w", fn, storage.ErrURLExists)
			}
		}
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return lastInsertId, nil
}

// GetURL method
func (s *Storage) GetURL(alias string) (string, error) {
	const fn = "storage.postgres.GetURL"
	var resURL string

	sqlStatement := `SELECT * FROM url WHERE alias=$1`
	row := s.db.QueryRow(sqlStatement, alias)
	err := row.Scan(&resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return resURL, nil
}

// DeleteURL method
func (s *Storage) DeleteURL(alias string) error {
	const fn = "storage.postgres.DeleteURL"

	sqlStatement := `DELETE FROM url WHERE alias=$1`
	_, err := s.db.Exec(sqlStatement, alias)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
