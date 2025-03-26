package usecases

import (
	"errors"
	"fmt"
	"github.com/n-kazachuk/go_parser/internal/app/domain/ticket"
	"github.com/n-kazachuk/go_parser/internal/app/domain/tickets-request"
	"github.com/n-kazachuk/go_parser/internal/libs/helpers"
	"github.com/n-kazachuk/go_parser/internal/libs/logger/sl"
)

func (s *UseCases) PushToQueue(ticketRequest *tickets_request.TicketRequest) error {
	log := sl.WithTrace(s.log)

	err := s.ticketsRequestsStorage.Add(ticketRequest)
	if err != nil {
		errMsg := "failed adding new ticket request to queue"
		log.Error(errMsg, sl.Err(err))

		return errors.New(errMsg)
	}

	return nil
}

func (s *UseCases) GetFreeFromQueue() (*tickets_request.TicketRequest, error) {
	op := helpers.GetFunctionName()

	ticketRequest, err := s.ticketsRequestsStorage.GetFree(s.cfg.Parser.Interval)
	if err != nil {
		errMsg := fmt.Sprintf("%s: error while taking free ticket request", op)
		s.log.Error(errMsg, sl.Err(err))

		return nil, errors.New(errMsg)
	}

	if ticketRequest == nil {
		return nil, nil
	}

	err = s.ticketsRequestsStorage.SetPicked(ticketRequest)
	if err != nil {
		errMsg := fmt.Sprintf("%s: error while setting ticket request picked", op)
		s.log.Error(errMsg, sl.Err(err))

		return nil, errors.New(errMsg)
	}

	return ticketRequest, nil
}

func (s *UseCases) GetTicketsFromSource(ticketRequest *tickets_request.TicketRequest) ([]*ticket.Ticket, error) {
	op := helpers.GetFunctionName()

	tickets, err := s.ticketsGateway.GetTickets(ticketRequest)
	if err != nil {
		errMsg := fmt.Sprintf("%s: error while getting orders from storage", op)
		s.log.Error(errMsg, sl.Err(err))

		return nil, errors.New(errMsg)
	}

	return tickets, nil
}

func (s *UseCases) SaveTicketsToStorage(tickets []*ticket.Ticket) error {
	op := helpers.GetFunctionName()

	if len(tickets) == 0 {
		return nil
	}

	err := s.ticketsStorage.Save(tickets)
	if err != nil {
		errMsg := fmt.Sprintf("%s: error while saving orders from storage", op)
		s.log.Error(errMsg, sl.Err(err))

		return errors.New(errMsg)
	}

	return nil
}

func (s *UseCases) SetProcessed(ticketRequest *tickets_request.TicketRequest) error {
	op := helpers.GetFunctionName()

	err := s.ticketsRequestsStorage.SetProcessed(ticketRequest)
	if err != nil {
		errMsg := fmt.Sprintf("%s: error while making ticket request processed", op)
		s.log.Error(errMsg, sl.Err(err))

		return errors.New(errMsg)
	}

	return nil
}
