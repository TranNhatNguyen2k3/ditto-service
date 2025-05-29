package ditto

import (
	"ditto/config"
	"ditto/internal/influxdb"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(func(cfg *config.Config) *Client {
		return NewClient(
			cfg.Ditto.URL,
			cfg.Ditto.Username,
			cfg.Ditto.Password,
		)
	}),
	fx.Provide(NewService),
	influxdb.Module,
)
