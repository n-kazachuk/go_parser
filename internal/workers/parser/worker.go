package parser

import (
	"fmt"
	"github.com/n-kazachuk/go_parser/internal/config"
	"github.com/n-kazachuk/go_parser/internal/domain/models"
	"github.com/n-kazachuk/go_parser/internal/lib/logger/sl"
	parserService "github.com/n-kazachuk/go_parser/internal/services/parser"
	"github.com/n-kazachuk/go_parser/internal/services/ticket_request"
	"log/slog"
	"sync"
	"time"
)

type Worker struct {
	id                   int
	log                  *slog.Logger
	cfg                  *config.Config
	wg                   *sync.WaitGroup
	mu                   *sync.Mutex
	parserService        *parserService.Parse
	ticketRequestService *ticket_request.TicketRequest

	stopCh chan struct{}
}

func New(
	id int,
	log *slog.Logger,
	cfg *config.Config,
	wg *sync.WaitGroup,
	mu *sync.Mutex,
	parserService *parserService.Parse,
	ticketRequestService *ticket_request.TicketRequest,
) *Worker {
	return &Worker{
		id,
		log,
		cfg,
		wg,
		mu,
		parserService,
		ticketRequestService,
		make(chan struct{}),
	}
}

func (w *Worker) Start() {
	const op = "worker.Start"

	w.log.Info(fmt.Sprintf("%s: worker #%v started", op, w.id))

	w.wg.Add(1)
	defer w.wg.Done()

	ticker := time.NewTicker(w.cfg.Parser.Worker.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			task, err := w.fetchTask()
			if err != nil {
				w.log.Error("%s: worker #%v error fetching task: %v\n", op, w.id, sl.Err(err))
				continue
			}

			if task == nil {
				continue
			}

			err = w.processTask(task)
			if err != nil {
				w.log.Error("%s: worker #%v error processing task: %v\n", op, w.id, sl.Err(err))
				continue
			}
		case <-w.stopCh:
			w.log.Info(fmt.Sprintf("%s: Worker #%v stopped", op, w.id))
			return
		}
	}
}

func (w *Worker) Stop() {
	close(w.stopCh)
}

func (w *Worker) fetchTask() (*models.TicketRequest, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	ticketRequest, err := w.ticketRequestService.GetFreeFromQueue()

	if err != nil {
		return nil, err
	}

	return ticketRequest, nil
}

func (w *Worker) processTask(ticketRequest *models.TicketRequest) error {
	const op = "worker.processTask"

	w.log.Info(fmt.Sprintf("%s: worker #%v processing task", op, w.id))

	tickets, err := w.parserService.GetTickets(ticketRequest.FromCity, ticketRequest.ToCity, ticketRequest.Date.Format("2006-01-02"))
	if err != nil {
		return err
	}

	err = w.parserService.SaveTickets(tickets)
	if err != nil {
		return err
	}

	err = w.ticketRequestService.SetProcessed(ticketRequest)
	if err != nil {
		return err
	}

	w.log.Info(fmt.Sprintf("%s: worker #%v procesed his task", op, w.id))

	return nil
}
