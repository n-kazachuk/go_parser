package models

type TicketRequest struct {
	Step     int
	FromCity string
	ToCity   string
	Date     string
	FromTime string
	ToTime   string
}

func NewTicketRequest() *TicketRequest {
	return &TicketRequest{}
}
