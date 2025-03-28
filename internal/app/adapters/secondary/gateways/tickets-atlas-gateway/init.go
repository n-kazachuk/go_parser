package tickets_atlas_gateway

import (
	"log/slog"

	"github.com/n-kazachuk/go_parser/internal/app/config"
)

const (
	DOMAIN = "atlasbus.by"
)

type TicketsAtlasGateway struct {
	log *slog.Logger
	cfg *config.GatewayConfig
}

func New(log *slog.Logger, cfg *config.GatewayConfig) *TicketsAtlasGateway {
	return &TicketsAtlasGateway{log, cfg}
}
