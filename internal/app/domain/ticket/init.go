package ticket

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

func New() *Ticket {
	return &Ticket{}
}
