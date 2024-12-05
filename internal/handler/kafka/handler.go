package kafka

import (
	"encoding/json"
	"fmt"
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
		h.log.Info(fmt.Sprintf("%s: %v", op, e.Value))

		ticketFindRequest := models.NewTicketRequest()

		err := json.Unmarshal(e.Value, ticketFindRequest)
		if err != nil {
			return err
		}

		h.log.Info(fmt.Sprintf("%s: %v", op, ticketFindRequest))

		return h.service.PushToQueue(ticketFindRequest)
	case kafka.Error:
		h.log.Error("Error with reading message: %v %s", sl.Err(e), op)
	}

	return nil
}
