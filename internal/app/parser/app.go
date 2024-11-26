package parserApp

import (
	"fmt"
	"github.com/n-kazachuk/go_parser/internal/config"
	"github.com/n-kazachuk/go_parser/internal/services/parser"
	"log/slog"
	"time"
)

type App struct {
	log     *slog.Logger
	cfg     *config.Config
	service *parser.Parser

	stopCh chan struct{}
}

func New(log *slog.Logger, cfg *config.Config, service *parser.Parser) *App {
	return &App{
		log,
		cfg,
		service,
		make(chan struct{}),
	}
}

// MustRun runs parser and panics if any error occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "parser.Run"
	return nil
	ticker := time.NewTicker(a.cfg.Parser.Interval)
	defer ticker.Stop()

	fromCity := "Минск"
	toCity := "Житковичи"
	date := "2024-08-10"

	for {
		select {
		case <-ticker.C:
			orders, err := a.service.GetOrders(fromCity, toCity, date)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}

			for _, value := range orders {
				fmt.Println(value)
			}

		case <-a.stopCh:
			fmt.Println("Stopping parser...")
			return nil
		}
	}
}

func (a *App) Stop() {
	const op = "parser.Stop"
	close(a.stopCh)
}
