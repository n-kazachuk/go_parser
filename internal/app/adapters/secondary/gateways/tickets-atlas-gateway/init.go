package tickets_atlas_gateway

import (
	"github.com/n-kazachuk/go_parser/internal/app/config"
	"log/slog"
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
