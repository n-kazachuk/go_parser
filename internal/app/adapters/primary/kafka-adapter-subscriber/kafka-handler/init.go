package kafka_handler

import (
	"github.com/n-kazachuk/go_parser/internal/app/application/usecases"
	"log/slog"
)

type KafkaHandler struct {
	log     *slog.Logger
	service *usecases.UseCases
}

func New(log *slog.Logger, service *usecases.UseCases) *KafkaHandler {
	return &KafkaHandler{
		log,
		service,
	}
}
