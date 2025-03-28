package kafka_handler

import (
	"log/slog"

	"github.com/n-kazachuk/go_parser/internal/app/application/usecases"
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
