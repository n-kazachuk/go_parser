package pgsql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/n-kazachuk/go_parser/internal/config"
	"github.com/n-kazachuk/go_parser/internal/domain/models"
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

// SaveOrders add new ticketRequest to queue
func (s *Storage) SaveOrders(orders []*models.Order) error {
	const op = "pgsql.AddTicketRequestToQueue"

	return nil
}

// AddTicketRequestToQueue add new ticketRequest to queue
func (s *Storage) AddTicketRequestToQueue(ticketRequest *models.TicketRequest) error {
	const op = "pgsql.AddTicketRequestToQueue"

	stmt, err := s.db.Prepare("INSERT INTO search_ticket_queue(from_city, to_city, date, from_time, to_time) VALUES($1, $2, $3, $4, $5)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(ticketRequest.FromCity, ticketRequest.ToCity, ticketRequest.Date, ticketRequest.FromTime, ticketRequest.ToTime)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// GetFreeTicketRequestFromQueue add new ticketRequest to queue
func (s *Storage) GetFreeTicketRequestFromQueue() (*models.TicketRequest, error) {
	//const op = "pgsql.GetFreeTicketRequestFromQueue"
	//
	//stmt, err := s.db.Prepare("INSERT INTO search_ticket_queue(from_city, to_city, date, from_time, to_time) VALUES($1, $2, $3, $4, $5)")
	//if err != nil {
	//	return fmt.Errorf("%s: %w", op, err)
	//}
	//
	//_, err = stmt.Exec(ticketRequest.FromCity, ticketRequest.ToCity, ticketRequest.Date, ticketRequest.FromTime, ticketRequest.ToTime)
	//if err != nil {
	//	return fmt.Errorf("%s: %w", op, err)
	//}
	//
	//return nil

	return nil, nil
}

// SetTicketRequestPicked add new ticketRequest to queue
func (s *Storage) SetTicketRequestPicked() error {
	const op = "pgsql.SetTicketRequestPicked"

	//stmt, err := s.db.Prepare("INSERT INTO search_ticket_queue(from_city, to_city, date, from_time, to_time) VALUES($1, $2, $3, $4, $5)")
	//if err != nil {
	//	return fmt.Errorf("%s: %w", op, err)
	//}
	//
	//_, err = stmt.Exec(ticketRequest.FromCity, ticketRequest.ToCity, ticketRequest.Date, ticketRequest.FromTime, ticketRequest.ToTime)
	//if err != nil {
	//	return fmt.Errorf("%s: %w", op, err)
	//}

	return nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}
