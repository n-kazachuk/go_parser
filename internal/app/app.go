package app

import (
	"github.com/n-kazachuk/go_parser/internal/services/parser"
	"github.com/n-kazachuk/go_parser/internal/storage/atlas_parser"
	"log/slog"

	parserapp "github.com/n-kazachuk/go_parser/internal/app/parser"
	"github.com/n-kazachuk/go_parser/internal/config"
)

type App struct {
	Parser *parserapp.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	parserStorage := atlas_parser.NewAtlasStorage()
	parserService := parser.NewParser(parserStorage)
	parserApp := parserapp.New(log, cfg, parserService)

	return &App{
		Parser: parserApp,
	}
}
