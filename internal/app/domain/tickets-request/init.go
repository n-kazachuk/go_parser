package tickets_request

import (
	"time"
)

type TicketRequest struct {
	FromCity string    `json:"from_city"`
	ToCity   string    `json:"to_city"`
	Date     time.Time `json:"date"`
	FromTime time.Time `json:"from_time"`
	ToTime   time.Time `json:"to_time"`
}

func New() *TicketRequest {
	return &TicketRequest{}
}
