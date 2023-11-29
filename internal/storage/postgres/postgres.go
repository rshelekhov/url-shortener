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

func (s *Storage) Close() error {
	return s.db.Close()
}

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
