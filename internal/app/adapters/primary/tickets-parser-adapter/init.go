package tickets_parser_adapter

import (
	worker "github.com/n-kazachuk/go_parser/internal/app/adapters/primary/tickets-parser-adapter/tickets-parser-worker"

	"github.com/n-kazachuk/go_parser/internal/app/application/usecases"
	"github.com/n-kazachuk/go_parser/internal/app/config"
	"log/slog"
	"sync"
)

type TicketsParserAdapter struct {
	log     *slog.Logger
	cfg     *config.ParserConfig
	service *usecases.UseCases

	workersWg *sync.WaitGroup
	workersMu *sync.Mutex
	workers   []*worker.Worker
}

func New(
	log *slog.Logger,
	cfg *config.ParserConfig,
	service *usecases.UseCases,
) *TicketsParserAdapter {
	workersWg := &sync.WaitGroup{}
	workersMu := &sync.Mutex{}
	workers := make([]*worker.Worker, cfg.Worker.Count)

	return &TicketsParserAdapter{
		log,
		cfg,
		service,
		workersWg,
		workersMu,
		workers,
	}
}
