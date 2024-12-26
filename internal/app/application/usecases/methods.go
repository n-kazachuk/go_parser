package usecases

import (
	"fmt"
	"github.com/n-kazachuk/go_parser/internal/app/domain/ticket"
	"github.com/n-kazachuk/go_parser/internal/app/domain/tickets-request"
)

func (s *UseCases) PushToQueue(ticketRequest *tickets_request.TicketRequest) error {
	const op = "TicketRequest.PushToQueue"

	err := s.queueStorage.AddTicketRequestToQueue(ticketRequest)
	if err != nil {
		return fmt.Errorf("%s: %s", op, "Error while adding new ticket request to queue")
	}

	return nil
}

func (s *UseCases) GetFreeFromQueue() (*tickets_request.TicketRequest, error) {
	const op = "TicketRequest.GetFreeFromQueue"

	ticketRequest, err := s.queueStorage.GetFreeTicketRequestFromQueue(s.cfg.Parser.Interval)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, "Error while taking free ticket request")
	}

	if ticketRequest == nil {
		return nil, nil
	}

	err = s.queueStorage.SetTicketRequestPicked(ticketRequest)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, "Error while setting ticket request picked")
	}

	return ticketRequest, nil
}

func (s *UseCases) GetTicketsFromSource(ticketRequest *tickets_request.TicketRequest) ([]*ticket.Ticket, error) {
	const op = "Parse.GetOrders"

	tickets, err := s.ticketsGateway.GetTickets(ticketRequest)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, "Error while getting orders from storage")
	}

	return tickets, nil
}

func (s *UseCases) SaveTickets(tickets []*ticket.Ticket) error {
	const op = "Parse.GetOrders"

	if len(tickets) == 0 {
		return nil
	}

	err := s.ticketsStorage.SaveTickets(tickets)
	if err != nil {
		return fmt.Errorf("%s: %s", op, "Error while saving orders from storage")
	}

	return nil
}

func (s *UseCases) SetProcessed(ticketRequest *tickets_request.TicketRequest) error {
	const op = "TicketRequest.SetProcessed"

	err := s.queueStorage.SetTicketRequestProcessed(ticketRequest)
	if err != nil {
		return fmt.Errorf("%s: %s", op, "Error while making ticket request processed")
	}

	return nil
}
