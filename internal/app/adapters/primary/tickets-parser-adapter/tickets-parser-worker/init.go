package tickets_parser_worker

import (
	"github.com/n-kazachuk/go_parser/internal/app/application/usecases"
	"github.com/n-kazachuk/go_parser/internal/app/config"
	"log/slog"
	"sync"
)

type Worker struct {
	id      int
	log     *slog.Logger
	cfg     *config.ParserConfig
	wg      *sync.WaitGroup
	mu      *sync.Mutex
	service *usecases.UseCases
}

func New(
	id int,
	log *slog.Logger,
	cfg *config.ParserConfig,
	wg *sync.WaitGroup,
	mu *sync.Mutex,
	service *usecases.UseCases,
) *Worker {
	return &Worker{
		id,
		log,
		cfg,
		wg,
		mu,
		service,
	}
}