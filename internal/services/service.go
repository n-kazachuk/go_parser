package services

import (
	"github.com/n-kazachuk/go_parser/internal/services/parser"
)

type Service struct {
	Parser *parser.ParseService
}

func NewService() *Service {
	return &Service{
		Parser: parser.NewParseService(),
	}
}
