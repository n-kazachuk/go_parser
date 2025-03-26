package tickets_requests_repository_postgres

import (
	"database/sql"
	"fmt"
	"github.com/n-kazachuk/go_parser/internal/app/config"
)

type TicketsRequestsRepositoryPostgres struct {
	db  *sql.DB
	cfg *config.PgsqlConfig
}

func New(cfg *config.PgsqlConfig) *TicketsRequestsRepositoryPostgres {
	pgsqlconn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName,
	)

	db, err := sql.Open("postgres", pgsqlconn)
	if err != nil {
		panic(err)
	}

	return &TicketsRequestsRepositoryPostgres{db, cfg}
}
