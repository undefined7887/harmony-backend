package chatrepo

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewMongoRepository),
	fx.Invoke(NewMongoMigrationsRunner),
)
