package google

import "go.uber.org/fx"

var Module = fx.Provide(
	NewAuthService,
)
