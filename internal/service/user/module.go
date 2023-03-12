package userservice

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewService),
	fx.Invoke(NewServiceRunner),
)
