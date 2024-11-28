package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/n-kazachuk/go_parser/internal/domain/models"
	"github.com/n-kazachuk/go_parser/internal/lib/logger/sl"
	"github.com/n-kazachuk/go_parser/internal/services/ticket_request"
	"log/slog"
)

type Handler struct {
	log     *slog.Logger
	service *ticket_request.TicketRequest
}

func New(log *slog.Logger, service *ticket_request.TicketRequest) *Handler {
	return &Handler{
		log,
		service,
	}
}

func (h *Handler) HandleTicketFindRequest(event kafka.Event) error {
	const op = "kafkaHandler.handleTicketFindRequest"

	switch e := event.(type) {
	case *kafka.Message:
		ticketFindRequest := models.NewTicketRequest()
		ticketFindRequest.FromCity = "Москва"
		ticketFindRequest.ToCity = "Минск"
		ticketFindRequest.Date = "2024-11-28"
		ticketFindRequest.FromTime = "13:00"
		ticketFindRequest.ToTime = "18:00"

		return h.service.PushToQueue(ticketFindRequest)
	case kafka.Error:
		h.log.Error("Error with reading message: %v %s", sl.Err(e), op)
	}

	return nil
}
