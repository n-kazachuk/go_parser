package parser

import (
	"github.com/n-kazachuk/go_parser/internal/domain/models"
)

type Storage interface {
	GetOrders(fromCity, toCity, date string) ([]*models.Order, error)
}

type Parser struct {
	Storage
}

func NewParser(storage Storage) *Parser {
	return &Parser{
		Storage: storage,
	}
}

func (s *Parser) GetOrders(fromCity, toCity, date string) ([]*models.Order, error) {
	orders, err := s.Storage.GetOrders(fromCity, toCity, date)

	if err != nil {
		return nil, err
	}

	return orders, nil
}
