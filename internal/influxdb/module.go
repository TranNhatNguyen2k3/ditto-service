package influxdb

import (
	"ditto/config"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(func(cfg *config.Config) *Client {
		return NewClient(
			"http://localhost:8086",
			"my-super-secret-token",
			"myorg",
			"ws_events",
		)
	}),
)
