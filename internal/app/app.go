package app

import (
	kafkaAdapterSubscriber "github.com/n-kazachuk/go_parser/internal/app/adapters/primary/kafka-adapter-subscriber"
	osSignalAdapter "github.com/n-kazachuk/go_parser/internal/app/adapters/primary/os-signal-adapter"
	ticketsParserAdapter "github.com/n-kazachuk/go_parser/internal/app/adapters/primary/tickets-parser-adapter"
	ticketsAtlasGateway "github.com/n-kazachuk/go_parser/internal/app/adapters/secondary/gateways/tickets-atlas-gateway"
	ticketsRepositoryPostgres "github.com/n-kazachuk/go_parser/internal/app/adapters/secondary/repositories/tickets-repository-postgres"

	"github.com/n-kazachuk/go_parser/internal/app/application/usecases"
	"github.com/n-kazachuk/go_parser/internal/app/config"
	"log/slog"
)

type App struct {
	OsSignalAdapter        *osSignalAdapter.OsSignalAdapter
	KafkaSubscriberAdapter *kafkaAdapterSubscriber.KafkaAdapterSubscriber
	TicketsParserAdapter   *ticketsParserAdapter.TicketsParserAdapter
}

func New(log *slog.Logger, cfg *config.Config) *App {
	ticketsRepository := ticketsRepositoryPostgres.New(&cfg.Pgsql)
	ticketsGateway := ticketsAtlasGateway.New(log, &cfg.Gateway)

	usc := usecases.New(
		log,
		cfg,
		ticketsRepository,
		ticketsRepository,
		ticketsGateway,
	)

	osAdapter := osSignalAdapter.New(log)
	kafkaSubscriber := kafkaAdapterSubscriber.New(log, &cfg.Kafka, usc)
	parserAdapter := ticketsParserAdapter.New(log, &cfg.Parser, usc)

	return &App{
		osAdapter,
		kafkaSubscriber,
		parserAdapter,
	}
}
