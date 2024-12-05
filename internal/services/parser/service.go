package parser

import (
	"fmt"
	"github.com/n-kazachuk/go_parser/internal/domain/models"
)

type OrdersGetterStorage interface {
	GetTickets(fromCity, toCity, date string) ([]*models.Ticket, error)
}

type OrdersSaverStorage interface {
	SaveTickets(orders []*models.Ticket) error
}

type Parse struct {
	OrdersGetterStorage
	OrdersSaverStorage
}

func New(
	ordersGetter OrdersGetterStorage,
	ordersSaver OrdersSaverStorage,
) *Parse {
	return &Parse{
		ordersGetter,
		ordersSaver,
	}
}

func (s *Parse) GetTickets(fromCity, toCity, date string) ([]*models.Ticket, error) {
	const op = "Parse.GetOrders"

	orders, err := s.OrdersGetterStorage.GetTickets(fromCity, toCity, date)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, "Error while getting orders from storage")
	}

	return orders, nil
}

func (s *Parse) SaveTickets(tickets []*models.Ticket) error {
	const op = "Parse.GetOrders"

	if len(tickets) == 0 {
		return nil
	}

	err := s.OrdersSaverStorage.SaveTickets(tickets)
	if err != nil {
		return fmt.Errorf("%s: %s", op, "Error while saving orders from storage")
	}

	return nil
}
