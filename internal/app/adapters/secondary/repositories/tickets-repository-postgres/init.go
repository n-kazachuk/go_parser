package tickets_repository_postgres

import (
	"database/sql"
	"fmt"

	"github.com/n-kazachuk/go_parser/internal/app/config"
)

type TicketsRepositoryPostgres struct {
	db  *sql.DB
	cfg *config.PgsqlConfig
}

func New(cfg *config.PgsqlConfig) *TicketsRepositoryPostgres {
	pgsqlconn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName,
	)

	db, err := sql.Open("postgres", pgsqlconn)
	if err != nil {
		panic(err)
	}

	return &TicketsRepositoryPostgres{db, cfg}
}
