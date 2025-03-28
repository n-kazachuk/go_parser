package tickets_parser_adapter

import (
	"context"
	"fmt"

	"github.com/n-kazachuk/go_parser/internal/libs/helpers"

	worker "github.com/n-kazachuk/go_parser/internal/app/adapters/primary/tickets-parser-adapter/tickets-parser-worker"
)

func (a *TicketsParserAdapter) Start(ctx context.Context) error {
	op := helpers.GetFunctionName()

	workersCount := a.cfg.Worker.Count
	if workersCount <= 0 {
		return fmt.Errorf("%s: workers count can't be empty", op)
	}

	for i := 0; i < workersCount; i++ {
		a.workersWg.Add(1)

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

	a.log.Info("TicketsParserAdapter gracefully stopped")

	return nil
}
