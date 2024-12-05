package models

import "time"

type Ticket struct {
	FromCity string
	ToCity   string
	Date     time.Time
	FromTime time.Time
	ToTime   time.Time
	Price    float64
	IsFree   bool
}

func NewTicket() *Ticket {
	return &Ticket{}
}
