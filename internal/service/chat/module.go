package chatservice

import "go.uber.org/fx"

var Module = fx.Provide(
	NewService,
)
