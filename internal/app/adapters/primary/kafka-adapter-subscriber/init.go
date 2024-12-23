package kafka_adapter_subscriber

import (
	kafkaHandler "github.com/n-kazachuk/go_parser/internal/app/adapters/primary/kafka-adapter-subscriber/kafka-handler"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/n-kazachuk/go_parser/internal/app/application/usecases"
	"github.com/n-kazachuk/go_parser/internal/app/config"
	"log/slog"
	"strings"
)

const (
	ParserGroup            = "parser"
	TicketFindRequestTopic = "ticket_find_request"
)

type KafkaAdapterSubscriber struct {
	log     *slog.Logger
	cfg     *config.KafkaConfig
	service *usecases.UseCases
	handler *kafkaHandler.KafkaHandler

	consumer *kafka.Consumer
}

func New(log *slog.Logger, cfg *config.KafkaConfig, service *usecases.UseCases) *KafkaAdapterSubscriber {
	handler := kafkaHandler.New(log, service)

	consumerConfig := &kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(cfg.Brokers, ","),
		"group.id":                 ParserGroup,
		"enable.auto.offset.store": true,
		"enable.auto.commit":       true,
	}

	consumer, err := kafka.NewConsumer(consumerConfig)
	if err != nil {
		panic(err)
	}

	return &KafkaAdapterSubscriber{
		log,
		cfg,
		service,
		handler,
		consumer,
	}
}
