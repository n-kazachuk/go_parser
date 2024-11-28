package models

type TicketRequest struct {
	FromCity string
	ToCity   string
	Date     string
	FromTime string
	ToTime   string
}

func NewTicketRequest() *TicketRequest {
	return &TicketRequest{}
}
