package parser

import (
	"fmt"
	"github.com/n-kazachuk/go_parser/internal/domain/models"
)

type OrdersGetterStorage interface {
	GetOrders(fromCity, toCity, date string) ([]*models.Order, error)
}

type OrdersSaverStorage interface {
	SaveOrders(orders []*models.Order) error
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

func (s *Parse) GetOrders(fromCity, toCity, date string) ([]*models.Order, error) {
	const op = "Parse.GetOrders"

	orders, err := s.OrdersGetterStorage.GetOrders(fromCity, toCity, date)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, "Error while getting orders from storage")
	}

	return orders, nil
}

func (s *Parse) SaveOrders(orders []*models.Order) error {
	const op = "Parse.GetOrders"

	err := s.OrdersSaverStorage.SaveOrders(orders)
	if err != nil {
		return fmt.Errorf("%s: %s", op, "Error while saving orders from storage")
	}

	return nil
}
