package chatrepo

import "go.uber.org/fx"

var Module = fx.Options(
	// Message repository
	fx.Provide(NewMongoMessageRepository),
	fx.Invoke(NewMongoMessageMigrationsRunner),

	// Chat repository
	fx.Provide(NewMongoChatRepository),
)
