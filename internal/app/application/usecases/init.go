package usecases

import (
	"github.com/n-kazachuk/go_parser/internal/app/config"
	"github.com/n-kazachuk/go_parser/internal/app/domain/ticket"
	"github.com/n-kazachuk/go_parser/internal/app/domain/tickets-request"
	"log/slog"
	"time"
)

type UseCases struct {
	log            *slog.Logger
	cfg            *config.Config
	queueStorage   queueStorage
	ticketsStorage ticketsStorage
	ticketsGateway ticketsGateway
}

type queueStorage interface {
	AddTicketRequestToQueue(ticketRequest *tickets_request.TicketRequest) error
	GetFreeTicketRequestFromQueue(expiredInterval time.Duration) (*tickets_request.TicketRequest, error)
	SetTicketRequestPicked(ticketRequest *tickets_request.TicketRequest) error
	SetTicketRequestProcessed(ticketRequest *tickets_request.TicketRequest) error
}

type ticketsStorage interface {
	SaveTickets(orders []*ticket.Ticket) error
}

type ticketsGateway interface {
	GetTickets(ticketRequest *tickets_request.TicketRequest) ([]*ticket.Ticket, error)
}

func New(
	log *slog.Logger,
	cfg *config.Config,
	queueStorage queueStorage,
	ticketsStorage ticketsStorage,
	ticketsGateway ticketsGateway,
) *UseCases {
	return &UseCases{
		log,
		cfg,
		queueStorage,
		ticketsStorage,
		ticketsGateway,
	}
}
