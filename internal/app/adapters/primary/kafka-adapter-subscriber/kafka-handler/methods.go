package kafka_handler

import (
	ticketsRequest "github.com/n-kazachuk/go_parser/internal/app/domain/tickets-request"

	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/n-kazachuk/go_parser/internal/libs/logger/sl"
)

func (h *KafkaHandler) HandleTicketFindRequest(event kafka.Event) error {
	const op = "kafkaHandler.handleTicketFindRequest"

	switch e := event.(type) {
	case *kafka.Message:
		h.log.Info(fmt.Sprintf("%s: %v", op, e.Value))

		ticketFindRequest := ticketsRequest.New()

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
