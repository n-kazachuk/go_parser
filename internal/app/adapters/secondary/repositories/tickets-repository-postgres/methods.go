package tickets_repository_postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"github.com/n-kazachuk/go_parser/internal/app/domain/ticket"
	"github.com/n-kazachuk/go_parser/internal/app/domain/tickets-request"
	"strings"
	"time"
)

func (s *TicketsRepositoryPostgres) SaveTickets(tickets []*ticket.Ticket) error {
	const op = "pgsql.SaveTickets"

	query := `
		INSERT INTO ticket (from_city, to_city, date, from_time, to_time, price, is_free)
		VALUES %s
		ON CONFLICT (from_city, to_city, date, from_time, to_time)
		DO UPDATE SET price = EXCLUDED.price
	`

	valueStrings := make([]string, 0, len(tickets))
	valueArgs := make([]interface{}, 0, len(tickets)*6)

	for i, ticket := range tickets {
		valueStrings = append(
			valueStrings,
			fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*7+1, i*7+2, i*7+3, i*7+4, i*7+5, i*7+6, i*7+7),
		)

		valueArgs = append(
			valueArgs,
			ticket.FromCity, ticket.ToCity, ticket.Date, ticket.FromTime, ticket.ToTime, ticket.Price, ticket.IsFree,
		)
	}

	finalQuery := fmt.Sprintf(query, strings.Join(valueStrings, ", "))

	_, err := s.db.Exec(finalQuery, valueArgs...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// AddTicketRequestToQueue add new ticketRequest to queue
func (s *TicketsRepositoryPostgres) AddTicketRequestToQueue(ticketRequest *tickets_request.TicketRequest) error {
	const op = "pgsql.AddTicketRequestToQueue"

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

// GetFreeTicketRequestFromQueue retrieves the oldest free ticket request from the queue
func (s *TicketsRepositoryPostgres) GetFreeTicketRequestFromQueue(expiredInterval time.Duration) (*tickets_request.TicketRequest, error) {
	const op = "pgsql.GetFreeTicketRequestFromQueue"

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

// SetTicketRequestPicked add new ticketRequest to queue
func (s *TicketsRepositoryPostgres) SetTicketRequestPicked(ticketRequest *tickets_request.TicketRequest) error {
	const op = "pgsql.SetTicketRequestPicked"

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

// SetTicketRequestProcessed updates ticket request in the queue
func (s *TicketsRepositoryPostgres) SetTicketRequestProcessed(ticketRequest *tickets_request.TicketRequest) error {
	const op = "pgsql.SetTicketRequestProcessed"

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

func (s *TicketsRepositoryPostgres) Stop() error {
	return s.db.Close()
}
