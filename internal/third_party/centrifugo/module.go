package centrifugo

import "go.uber.org/fx"

var Module = fx.Provide(
	NewClient,
)
