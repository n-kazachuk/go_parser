package tickets_parser_worker

import (
	"context"
	"fmt"
	"github.com/n-kazachuk/go_parser/internal/app/domain/tickets-request"
	"github.com/n-kazachuk/go_parser/internal/libs/logger/sl"
	"time"
)

func (w *Worker) Start(ctx context.Context) {
	const op = "worker.Start"

	w.log.Info(fmt.Sprintf("%s: worker #%v started", op, w.id))

	defer w.wg.Done()

	ticker := time.NewTicker(w.cfg.Worker.Interval)
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
		case <-ctx.Done():
			w.Stop()
			return
		}
	}
}

func (w *Worker) Stop() {
	const op = "parser.Stop"
	w.log.Info(fmt.Sprintf("%s: Worker #%v stopped", op, w.id))
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
	const op = "worker.processTask"

	w.log.Info(fmt.Sprintf("%s: worker #%v processing task", op, w.id))

	tickets, err := w.service.GetTicketsFromSource(ticketRequest)
	if err != nil {
		return err
	}

	err = w.service.SaveTickets(tickets)
	if err != nil {
		return err
	}

	err = w.service.SetProcessed(ticketRequest)
	if err != nil {
		return err
	}

	w.log.Info(fmt.Sprintf("%s: worker #%v procesed his task", op, w.id))

	return nil
}
