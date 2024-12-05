package ticket_request

import (
	"fmt"
	"github.com/n-kazachuk/go_parser/internal/domain/models"
)

type QueueStorage interface {
	AddTicketRequestToQueue(ticketRequest *models.TicketRequest) error
	GetFreeTicketRequestFromQueue() (*models.TicketRequest, error)
	SetTicketRequestPicked(ticketRequest *models.TicketRequest) error
	SetTicketRequestProcessed(ticketRequest *models.TicketRequest) error
}

type TicketRequest struct {
	QueueStorage QueueStorage
}

func New(queueStorage QueueStorage) *TicketRequest {
	return &TicketRequest{
		queueStorage,
	}
}

func (s *TicketRequest) PushToQueue(ticketRequest *models.TicketRequest) error {
	const op = "TicketRequest.PushToQueue"

	err := s.QueueStorage.AddTicketRequestToQueue(ticketRequest)
	if err != nil {
		return fmt.Errorf("%s: %s", op, "Error while adding new ticket request to queue")
	}

	return nil
}

func (s *TicketRequest) GetFreeFromQueue() (*models.TicketRequest, error) {
	const op = "TicketRequest.GetFreeFromQueue"

	ticketRequest, err := s.QueueStorage.GetFreeTicketRequestFromQueue()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, "Error while taking free ticket request")
	}

	if ticketRequest == nil {
		return nil, nil
	}

	err = s.QueueStorage.SetTicketRequestPicked(ticketRequest)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, "Error while setting ticket request picked")
	}

	return ticketRequest, nil
}

func (s *TicketRequest) SetProcessed(ticketRequest *models.TicketRequest) error {
	const op = "TicketRequest.SetProcessed"

	err := s.QueueStorage.SetTicketRequestProcessed(ticketRequest)
	if err != nil {
		return fmt.Errorf("%s: %s", op, "Error while making ticket request processed")
	}

	return nil
}
