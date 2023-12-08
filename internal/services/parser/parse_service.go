package parser

import "github.com/n-kazachuk/go_parser/internal/models"

type WebParser interface {
	ParseOrders(fromCity, toCity, date string) ([]*models.Order, error)
}

type ParseService struct {
	WebParser
}

func NewParseService() *ParseService {
	return &ParseService{
		WebParser: NewAtlasParser(),
	}
}

func (s *ParseService) GetOrders(fromCity, toCity, date string) ([]*models.Order, error) {
	orders, err := s.WebParser.ParseOrders(fromCity, toCity, date)

	if err != nil {
		return nil, err
	}

	return orders, nil
}
