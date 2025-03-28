package tickets_dummy_gateway

import (
	"math/rand"
	"time"

	"github.com/n-kazachuk/go_parser/internal/app/domain/ticket"

	ticketsRequest "github.com/n-kazachuk/go_parser/internal/app/domain/tickets-request"
)

func (s *TicketsDummyGateway) GetTickets(ticketRequest *ticketsRequest.TicketRequest) ([]*ticket.Ticket, error) {
	tickets := make([]*ticket.Ticket, rand.Intn(10)+1)
	for i := range tickets {
		tickets[i] = &ticket.Ticket{
			FromCity: ticketRequest.FromCity,
			ToCity:   ticketRequest.ToCity,
			Date:     ticketRequest.Date,
			FromTime: normalizeTime(randomTime(ticketRequest.FromTime, ticketRequest.ToTime)),
			ToTime:   normalizeTime(randomTime(ticketRequest.FromTime, ticketRequest.ToTime)),
			Price:    randomPrice(10, 100),
			IsFree:   rand.Intn(2) == 0,
		}
	}

	return tickets, nil
}

func randomTime(from, to time.Time) time.Time {
	if from.After(to) {
		from, to = to, from
	}
	delta := to.Sub(from)
	randomOffset := time.Duration(rand.Int63n(int64(delta)))
	return from.Add(randomOffset)
}

func normalizeTime(t time.Time) time.Time {
	minutes := t.Minute()
	roundedMinutes := (minutes / 15) * 15
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), roundedMinutes, 0, 0, t.Location())
}

func randomPrice(min, max int) float64 {
	return float64(min + rand.Intn(max-min+1))
}
