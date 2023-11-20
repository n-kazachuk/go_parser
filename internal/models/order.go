package models

type Order struct {
	timeFrom string
	timeTo   string
	cost     int
	isFree   bool
}

func NewOrder(timeFrom, timeTo string, isFree bool) *Order {
	return &Order{
		timeFrom: timeFrom,
		timeTo:   timeTo,
		isFree:   isFree,
	}
}
