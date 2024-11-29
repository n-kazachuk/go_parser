package ticket_request

import (
	"github.com/n-kazachuk/go_parser/internal/domain/models"
)

type QueueStorage interface {
	AddTicketRequestToQueue(ticketRequest *models.TicketRequest) error
	GetFreeTicketRequestFromQueue() (*models.TicketRequest, error)
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
	err := s.QueueStorage.AddTicketRequestToQueue(ticketRequest)
	if err != nil {
		return err
	}

	return nil
}

func (s *TicketRequest) GetFreeFromQueue() (*models.TicketRequest, error) {
	ticketRequest, err := s.QueueStorage.GetFreeTicketRequestFromQueue()

	if err != nil {
		return nil, err
	}

	return ticketRequest, nil
}
