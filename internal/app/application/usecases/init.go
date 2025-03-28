package usecases

import (
	"log/slog"
	"time"

	"github.com/n-kazachuk/go_parser/internal/app/config"
	"github.com/n-kazachuk/go_parser/internal/app/domain/ticket"
	"github.com/n-kazachuk/go_parser/internal/app/domain/tickets-request"
)

type UseCases struct {
	log                    *slog.Logger
	cfg                    *config.Config
	ticketsGateway         ticketsGateway
	ticketsStorage         ticketsStorage
	ticketsRequestsStorage ticketsRequestsStorage
}

type ticketsGateway interface {
	GetTickets(ticketRequest *tickets_request.TicketRequest) ([]*ticket.Ticket, error)
}

type ticketsStorage interface {
	Save(tickets []*ticket.Ticket) error
}

type ticketsRequestsStorage interface {
	Add(request *tickets_request.TicketRequest) error
	GetFree(expiredInterval time.Duration) (*tickets_request.TicketRequest, error)
	SetPicked(request *tickets_request.TicketRequest) error
	SetProcessed(request *tickets_request.TicketRequest) error
}

func New(
	log *slog.Logger,
	cfg *config.Config,
	ticketsGateway ticketsGateway,
	ticketsStorage ticketsStorage,
	ticketsRequestsStorage ticketsRequestsStorage,
) *UseCases {
	return &UseCases{
		log:                    log,
		cfg:                    cfg,
		ticketsGateway:         ticketsGateway,
		ticketsStorage:         ticketsStorage,
		ticketsRequestsStorage: ticketsRequestsStorage,
	}
}
