package usecases

import (
	"github.com/n-kazachuk/go_parser/internal/app/config"
	"github.com/n-kazachuk/go_parser/internal/app/domain/model"
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
	AddTicketRequestToQueue(ticketRequest *model.TicketRequest) error
	GetFreeTicketRequestFromQueue(expiredInterval time.Duration) (*model.TicketRequest, error)
	SetTicketRequestPicked(ticketRequest *model.TicketRequest) error
	SetTicketRequestProcessed(ticketRequest *model.TicketRequest) error
}

type ticketsStorage interface {
	SaveTickets(orders []*model.Ticket) error
}

type ticketsGateway interface {
	GetTickets(fromCity, toCity, date string) ([]*model.Ticket, error)
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
