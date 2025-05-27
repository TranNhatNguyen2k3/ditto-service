package ditto

import (
	"ditto/config"
	"ditto/internal/influxdb"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(func(cfg *config.Config) *Client {
		return NewClient(
			"localhost:8080", // Ditto WebSocket port
			"ditto",          // username
			"ditto",          // password
		)
	}),
	fx.Provide(NewService),
	influxdb.Module,
)
