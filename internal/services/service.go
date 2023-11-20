package services

import (
	"github.com/n-kazachuk/go_parser/internal/models"
	"github.com/n-kazachuk/go_parser/internal/parsers"
)

type Parser interface {
	GetOrders(fromCity, toCity, date string) ([]*models.Order, error)
}

type Service struct {
	Parser
}

func NewService() *Service {
	return &Service{
		Parser: parsers.NewAtlasParser(),
	}
}

func (s *Service) ParseOrders(fromCity, toCity, date string) ([]*models.Order, error) {
	orders, err := s.Parser.GetOrders(fromCity, toCity, date)

	if err != nil {
		return nil, err
	}

	return orders, nil
}
