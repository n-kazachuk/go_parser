package os_signal_adapter

import (
	"context"
	"fmt"
	"github.com/n-kazachuk/go_parser/internal/libs/helpers"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type OsSignalAdapter struct {
	log *slog.Logger
}

func New(log *slog.Logger) *OsSignalAdapter {
	return &OsSignalAdapter{log}
}

func (a *OsSignalAdapter) Start(ctx context.Context) error {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-ctx.Done():
		a.Stop()
		return ctx.Err()
	case sig := <-ch:
		return fmt.Errorf("%s: system signal getted %s", helpers.GetFunctionName(), sig.String())
	}
}

func (a *OsSignalAdapter) Stop() {
	a.log.Info(fmt.Sprintf("OsSignalAdapter gracefully stopped"))
}
