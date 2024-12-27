package graceful

import (
	"context"
	"log/slog"

	"golang.org/x/sync/errgroup"
)

type starter interface {
	Start(ctx context.Context) error
}

type Graceful struct {
	processes []Process
	logger    *slog.Logger
}

func New(processes ...Process) *Graceful {
	return &Graceful{
		processes: processes,
	}
}

func (gr *Graceful) Start(ctx context.Context) {
	g, ctx := errgroup.WithContext(ctx)

	for _, process := range gr.processes {
		process := process // TODO remove if go > 1.22

		if process.disabled {
			continue
		}

		f := func() error {
			return process.starter.Start(ctx)
		}

		g.Go(f)

	}

	_ = g.Wait()

	gr.logger.Info("Application gracefully stopped")
}

func (gr *Graceful) SetLogger(l *slog.Logger) {
	gr.logger = l
}
