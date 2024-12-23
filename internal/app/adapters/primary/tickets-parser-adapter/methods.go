package tickets_parser_adapter

import (
	worker "github.com/n-kazachuk/go_parser/internal/app/adapters/primary/tickets-parser-adapter/tickets-parser-worker"

	"context"
	"fmt"
)

func (a *TicketsParserAdapter) Start(ctx context.Context) error {
	const op = "parser.Run"

	workersCount := a.cfg.Worker.Count
	if workersCount <= 0 {
		return fmt.Errorf("%s: workers count can't be empty", op)
	}

	for i := 0; i < workersCount; i++ {
		wrk := worker.New(
			i,
			a.log,
			a.cfg,
			a.workersWg,
			a.workersMu,
			a.service,
		)

		a.workers[i] = wrk

		go wrk.Start(ctx)
	}

	a.workersWg.Wait()

	return nil
}
