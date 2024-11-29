package parserApp

import (
	parserService "github.com/n-kazachuk/go_parser/internal/services/parser"
	parserWorker "github.com/n-kazachuk/go_parser/internal/workers/parser"
	"sync"

	"fmt"
	"github.com/n-kazachuk/go_parser/internal/config"
	"github.com/n-kazachuk/go_parser/internal/services/ticket_request"
	"log/slog"
)

type App struct {
	log                  *slog.Logger
	cfg                  *config.Config
	parserService        *parserService.Parse
	ticketRequestService *ticket_request.TicketRequest

	workersWg *sync.WaitGroup
	workers   []*parserWorker.Worker
}

func New(
	log *slog.Logger,
	cfg *config.Config,
	parserService *parserService.Parse,
	ticketRequestService *ticket_request.TicketRequest,
) *App {
	workersWg := &sync.WaitGroup{}
	workers := make([]*parserWorker.Worker, cfg.Parser.Worker.Count)

	return &App{
		log,
		cfg,
		parserService,
		ticketRequestService,
		workersWg,
		workers,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "parser.Run"

	workersCount := a.cfg.Parser.Worker.Count
	if workersCount <= 0 {
		return fmt.Errorf("%s: workers count can't be empty", op)
	}

	for i := 0; i < workersCount; i++ {
		worker := parserWorker.New(
			i,
			a.log,
			a.cfg,
			a.workersWg,
			a.parserService,
			a.ticketRequestService,
		)

		a.workers[i] = worker

		go worker.Start()
	}

	return nil
}

func (a *App) Stop() {
	const op = "parser.Stop"

	for index, worker := range a.workers {
		worker.Stop()
		a.log.Info(fmt.Sprintf("%s: parser worker %v stopped", op, index))
	}

	a.workersWg.Wait()

	a.log.Info(fmt.Sprintf("%s: parser stopped", op))
}
