package tickets_dummy_gateway

import (
	"log/slog"

	"github.com/n-kazachuk/go_parser/internal/app/config"
)

type TicketsDummyGateway struct {
	log *slog.Logger
	cfg *config.GatewayConfig
}

func New(log *slog.Logger, cfg *config.GatewayConfig) *TicketsDummyGateway {
	return &TicketsDummyGateway{log, cfg}
}
