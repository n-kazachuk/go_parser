package kafka

import (
	"database/sql"
	"fmt"
	"github.com/n-kazachuk/go_parser/internal/config"
)

type Storage struct {
	db *sql.DB
}

func New(cfg *config.Config) (*Storage, error) {
	const op = "storage.pgsql.New"

	pgsqlconn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Pgsql.Host, cfg.Pgsql.Port, cfg.Pgsql.User, cfg.Pgsql.Password, cfg.Pgsql.DbName,
	)

	db, err := sql.Open("postgres", pgsqlconn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}
