package tickets_parser_worker

import (
	"context"
	"fmt"
	"time"

	"github.com/n-kazachuk/go_parser/internal/app/domain/tickets-request"
	"github.com/n-kazachuk/go_parser/internal/libs/helpers"
	"github.com/n-kazachuk/go_parser/internal/libs/logger/sl"
)

func (w *Worker) Start(ctx context.Context) {
	op := helpers.GetFunctionName()

	w.log.Info("Worker started")

	defer w.wg.Done()

	ticker := time.NewTicker(w.cfg.Worker.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			task, err := w.fetchTask()
			if err != nil {
				w.log.Error(fmt.Sprintf("%s: error fetching worker task \n", op), sl.Err(err))
				continue
			}

			if task == nil {
				continue
			}

			err = w.processTask(task)
			if err != nil {
				w.log.Error(fmt.Sprintf("%s: error processing worker task \n", op), sl.Err(err))
				continue
			}
		case <-ctx.Done():
			w.Stop()
			return
		}
	}
}

func (w *Worker) Stop() {
	w.log.Info("Worker stopped")
}

func (w *Worker) fetchTask() (*tickets_request.TicketRequest, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	ticketRequest, err := w.service.GetFreeFromQueue()

	if err != nil {
		return nil, err
	}

	return ticketRequest, nil
}

func (w *Worker) processTask(ticketRequest *tickets_request.TicketRequest) error {
	w.log.Info("Worker start processing task")

	tickets, err := w.service.GetTicketsFromSource(ticketRequest)
	if err != nil {
		return err
	}

	err = w.service.SaveTicketsToStorage(tickets)
	if err != nil {
		return err
	}

	err = w.service.SetProcessed(ticketRequest)
	if err != nil {
		return err
	}

	w.log.Info("Worker processed his task")

	return nil
}
