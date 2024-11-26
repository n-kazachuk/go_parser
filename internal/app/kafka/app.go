package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/n-kazachuk/go_parser/internal/config"
	"github.com/n-kazachuk/go_parser/internal/domain/models"
	"github.com/n-kazachuk/go_parser/internal/lib/logger/sl"
	"github.com/n-kazachuk/go_parser/internal/services/ticket_request"
	"log/slog"
)

const TicketFindRequestTopic = "ticket_find_request"

type App struct {
	log     *slog.Logger
	cfg     *config.Config
	service *ticket_request.TicketRequest

	consumer sarama.Consumer
}

func New(log *slog.Logger, cfg *config.Config, service *ticket_request.TicketRequest) *App {
	consumerCfg := sarama.NewConfig()
	consumerCfg.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(cfg.Kafka.Brokers, consumerCfg)
	if err != nil {
		log.Error("Failed to create consumer: %v", sl.Err(err))
	}

	return &App{
		log,
		cfg,
		service,
		consumer,
	}
}

// MustRun runs parser and panics if any error occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "Kakfa.Run"

	a.log.Info(fmt.Sprintf("Running %s", op))

	go a.listenTicketsRequest()

	return nil
}

func (a *App) Stop() {
	const op = "parser.Stop"
	a.log.Info(fmt.Sprintf("Kafka stopped %s", op))

	_ = a.consumer.Close()
}

func (a *App) listenTicketsRequest() {
	const op = "Kafka.listenTicketsRequest"

	partitionCosumer, err := a.consumer.ConsumePartition(TicketFindRequestTopic, 0, sarama.OffsetNewest)

	if err != nil {
		a.log.Error("Failed to create consumer: %v", op, sl.Err(err))
	}

	a.log.Info(fmt.Sprintf("Consumer started %s", op))

	for {
		select {
		case err := <-partitionCosumer.Errors():
			a.log.Error("Error from partition error chanel: %v", op, sl.Err(err))
		case msg := <-partitionCosumer.Messages():
			var ticketRequest *models.TicketRequest

			err := json.Unmarshal(msg.Value, ticketRequest)
			if err != nil {
				a.log.Error("Error from partition error chanel: %v", op, sl.Err(err))
			}

			err = a.service.PushToQueue(ticketRequest)
			if err != nil {
				a.log.Error("Error from partition error chanel: %v", op, sl.Err(err))
			}

			a.log.Info(fmt.Sprintf("Message from partition %s", msg.Value))
		}
	}
}
