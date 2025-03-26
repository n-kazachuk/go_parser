package app

import (
	"log/slog"

	"github.com/n-kazachuk/go_parser/internal/app/application/usecases"
	"github.com/n-kazachuk/go_parser/internal/app/config"

	kafkaAdapterSubscriber "github.com/n-kazachuk/go_parser/internal/app/adapters/primary/kafka-adapter-subscriber"
	osSignalAdapter "github.com/n-kazachuk/go_parser/internal/app/adapters/primary/os-signal-adapter"
	ticketsParserAdapter "github.com/n-kazachuk/go_parser/internal/app/adapters/primary/tickets-parser-adapter"
	ticketsDummyGateway "github.com/n-kazachuk/go_parser/internal/app/adapters/secondary/gateways/tickets-dummy-gateway"
	ticketsRepositoryPostgres "github.com/n-kazachuk/go_parser/internal/app/adapters/secondary/repositories/tickets-repository-postgres"
	ticketsRequestsRepositoryPostgres "github.com/n-kazachuk/go_parser/internal/app/adapters/secondary/repositories/tickets-requests-repository-postgres"
)

type App struct {
	OsSignalAdapter        *osSignalAdapter.OsSignalAdapter
	KafkaSubscriberAdapter *kafkaAdapterSubscriber.KafkaAdapterSubscriber
	TicketsParserAdapter   *ticketsParserAdapter.TicketsParserAdapter
}

func New(log *slog.Logger, cfg *config.Config) *App {
	//ticketsGateway := ticketsAtlasGateway.New(log, &cfg.Gateway)
	ticketsGateway := ticketsDummyGateway.New(log, &cfg.Gateway)

	ticketsRepository := ticketsRepositoryPostgres.New(&cfg.Pgsql)
	ticketsRequestsRepository := ticketsRequestsRepositoryPostgres.New(&cfg.Pgsql)

	usc := usecases.New(
		log,
		cfg,
		ticketsGateway,
		ticketsRepository,
		ticketsRequestsRepository,
	)

	osAdapter := osSignalAdapter.New(log)
	kafkaSubscriber := kafkaAdapterSubscriber.New(log, &cfg.Kafka, usc)
	parserAdapter := ticketsParserAdapter.New(log, &cfg.Parser, usc)

	return &App{
		OsSignalAdapter:        osAdapter,
		KafkaSubscriberAdapter: kafkaSubscriber,
		TicketsParserAdapter:   parserAdapter,
	}
}
