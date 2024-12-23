package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type TicketRequest struct {
	FromCity string    `json:"from_city"`
	ToCity   string    `json:"to_city"`
	Date     time.Time `json:"date"`
	FromTime time.Time `json:"from_time"`
	ToTime   time.Time `json:"to_time"`
}

func NewTicketRequest() *TicketRequest {
	return &TicketRequest{}
}

func (t *TicketRequest) UnmarshalJSON(data []byte) error {
	type Alias TicketRequest
	aux := &struct {
		Date     string `json:"date"`
		FromTime string `json:"from_time"`
		ToTime   string `json:"to_time"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var err error

	dateFormat := "2006-01-02"
	timeFormat := "15:04:05"

	t.Date, err = time.Parse(dateFormat, aux.Date)
	if err != nil {
		return fmt.Errorf("failed to parse date: %v", err)
	}

	t.FromTime, err = time.Parse(timeFormat, aux.FromTime)
	if err != nil {
		return fmt.Errorf("failed to parse from_time: %v", err)
	}

	t.ToTime, err = time.Parse(timeFormat, aux.ToTime)
	if err != nil {
		return fmt.Errorf("failed to parse to_time: %v", err)
	}

	return nil
}
