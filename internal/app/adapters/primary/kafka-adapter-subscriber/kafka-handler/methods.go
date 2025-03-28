package kafka_handler

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	ticketsRequest "github.com/n-kazachuk/go_parser/internal/app/domain/tickets-request"
)

func (h *KafkaHandler) HandleTicketFindRequest(event kafka.Event) error {
	switch e := event.(type) {
	case *kafka.Message:
		h.log.Info(fmt.Sprintf("KafkaHandler start handle message: %v", e.Value))

		ticketFindRequest := ticketsRequest.New()

		err := json.Unmarshal(e.Value, ticketFindRequest)
		if err != nil {
			return err
		}

		return h.service.PushToQueue(ticketFindRequest)
	case kafka.Error:
		return e
	}

	return nil
}
