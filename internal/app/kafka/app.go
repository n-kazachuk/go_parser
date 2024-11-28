package kafka

import (
	"fmt"
	"github.com/n-kazachuk/go_parser/internal/config"
	"github.com/n-kazachuk/go_parser/internal/lib/logger/sl"
	"log/slog"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	ParserGroup            = "parser"
	TicketFindRequestTopic = "ticket_find_request"
)

type Handler interface {
	HandleTicketFindRequest(event kafka.Event) error
}

type App struct {
	log     *slog.Logger
	cfg     *config.Config
	handler Handler

	consumer *kafka.Consumer
	stopCh   chan struct{}
}

func New(log *slog.Logger, cfg *config.Config, handler Handler) *App {
	consumerConfig := &kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(cfg.Kafka.Brokers, ","),
		"group.id":                 ParserGroup,
		"enable.auto.offset.store": true,
		"enable.auto.commit":       true,
	}

	consumer, err := kafka.NewConsumer(consumerConfig)
	if err != nil {
		log.Error("Failed to create consumer: %v", sl.Err(err))
	}

	return &App{
		log,
		cfg,
		handler,
		consumer,
		make(chan struct{}),
	}
}

// MustRun runs parser and panics if any error occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "Kafka.Run"

	a.log.Info(fmt.Sprintf("Running %s", op))

	go a.listenTicketsRequest()

	return nil
}

func (a *App) Stop() {
	const op = "kafka.Stop"

	close(a.stopCh)

	if err := a.consumer.Close(); err != nil {
		a.log.Error("Failed to close consumer: %v", sl.Err(err))
	}

	a.log.Info(fmt.Sprintf("Kafka stopped %s", op))
}

func (a *App) listenTicketsRequest() {
	const op = "Kafka.listenTicketsRequest"

	err := a.consumer.Subscribe(TicketFindRequestTopic, nil)
	if err != nil {
		a.log.Error("Failed to subscribe topic: %v", sl.Err(err))
	}

	a.log.Info(fmt.Sprintf("Consumer started: %s", op))

	for {
		select {
		case <-a.stopCh:
			a.log.Info(fmt.Sprintf("Exiting listenTicketsRequest %s", op))
			return
		default:
			event := a.consumer.Poll(a.cfg.Kafka.Interval)
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
