package tickets_requests_repository_postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/n-kazachuk/go_parser/internal/app/domain/tickets-request"
	"github.com/n-kazachuk/go_parser/internal/libs/helpers"
)

func (s *TicketsRequestsRepositoryPostgres) Add(ticketRequest *tickets_request.TicketRequest) error {
	op := helpers.GetFunctionName()

	stmt, err := s.db.Prepare(`
		INSERT INTO search_ticket_queue (from_city, to_city, date, from_time, to_time) 
		VALUES($1, $2, $3, $4, $5)
	`)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(ticketRequest.FromCity, ticketRequest.ToCity, ticketRequest.Date, ticketRequest.FromTime, ticketRequest.ToTime)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation" {
			return nil
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *TicketsRequestsRepositoryPostgres) GetFree(expiredInterval time.Duration) (*tickets_request.TicketRequest, error) {
	op := helpers.GetFunctionName()

	timeout := time.Now().Add(-expiredInterval)

	stmt, err := s.db.Prepare(`
		SELECT from_city, to_city, date, from_time, to_time 
		FROM search_ticket_queue 
		WHERE is_picked = FALSE AND (searched_time IS NULL OR searched_time <= $1)
		ORDER BY searched_time ASC 
		LIMIT 1
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(timeout)
	ticketRequest := tickets_request.New()

	err = row.Scan(
		&ticketRequest.FromCity,
		&ticketRequest.ToCity,
		&ticketRequest.Date,
		&ticketRequest.FromTime,
		&ticketRequest.ToTime,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: error scanning row: %w", op, err)
	}

	return ticketRequest, nil
}

func (s *TicketsRequestsRepositoryPostgres) SetPicked(ticketRequest *tickets_request.TicketRequest) error {
	op := helpers.GetFunctionName()

	stmt, err := s.db.Prepare("UPDATE search_ticket_queue SET is_picked = TRUE WHERE from_city = $1 AND to_city = $2 AND date = $3")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(ticketRequest.FromCity, ticketRequest.ToCity, ticketRequest.Date)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *TicketsRequestsRepositoryPostgres) SetProcessed(ticketRequest *tickets_request.TicketRequest) error {
	op := helpers.GetFunctionName()

	stmt, err := s.db.Prepare(`
		UPDATE search_ticket_queue 
		SET is_picked = FALSE, searched_time = NOW() 
		WHERE from_city = $1 AND to_city = $2 AND date = $3
	`)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(ticketRequest.FromCity, ticketRequest.ToCity, ticketRequest.Date)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *TicketsRequestsRepositoryPostgres) Stop() error {
	return s.db.Close()
}
