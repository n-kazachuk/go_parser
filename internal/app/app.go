package app

import (
	kafkaApp "github.com/n-kazachuk/go_parser/internal/app/kafka"
	parserApp "github.com/n-kazachuk/go_parser/internal/app/parser"
	"github.com/n-kazachuk/go_parser/internal/storage/pgsql"

	"github.com/n-kazachuk/go_parser/internal/config"
	"github.com/n-kazachuk/go_parser/internal/services/parser"
	"github.com/n-kazachuk/go_parser/internal/services/ticket_request"
	"github.com/n-kazachuk/go_parser/internal/storage/atlas_parser"
	"log/slog"
)

type App struct {
	Kafka  *kafkaApp.App
	Parser *parserApp.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	ticketsStorage, err := pgsql.New(cfg)
	if err != nil {
		panic(err)
	}

	parserStorage := atlas_parser.NewAtlasStorage()
	parserService := parser.NewParser(parserStorage)
	parserApplication := parserApp.New(log, cfg, parserService)

	ticketRequestService := ticket_request.New(ticketsStorage)
	kafkaApplication := kafkaApp.New(log, cfg, ticketRequestService)

	return &App{
		Kafka:  kafkaApplication,
		Parser: parserApplication,
	}
}
