package kafka_adapter_subscriber

import (
	"context"
	"fmt"
	"github.com/n-kazachuk/go_parser/internal/libs/sl"
)

func (a *KafkaAdapterSubscriber) Start(ctx context.Context) error {
	const op = "Kafka.Run"

	a.log.Info(fmt.Sprintf("Running %s", op))

	go a.listenTicketsRequest(ctx)

	return nil
}

func (a *KafkaAdapterSubscriber) Stop() {
	if err := a.consumer.Close(); err != nil {
		a.log.Error("Failed to close consumer: %v", sl.Err(err))
	}

	a.log.Info(fmt.Sprintf("KafkaAdapterSubscriber gracefully stopped"))
}

func (a *KafkaAdapterSubscriber) listenTicketsRequest(ctx context.Context) {
	const op = "Kafka.listenTicketsRequest"

	err := a.consumer.Subscribe(TicketFindRequestTopic, nil)
	if err != nil {
		a.log.Error("Failed to subscribe topic: %v", sl.Err(err))
	}

	a.log.Info(fmt.Sprintf("Consumer started: %s", op))

	for {
		select {
		case <-ctx.Done():
			a.Stop()
			return
		default:
			event := a.consumer.Poll(a.cfg.Interval)
			if event == nil {
				continue
			}

			err := a.handler.HandleTicketFindRequest(event)
			if err != nil {
				a.log.Error("Error with reading message: %v", sl.Err(err))
			}
		}
	}
}
