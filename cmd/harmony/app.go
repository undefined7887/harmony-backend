package main

import (
	"os"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/undefined7887/harmony-backend/internal/config"
	mongodatabase "github.com/undefined7887/harmony-backend/internal/infrastructure/database/mongo"
	zaplog "github.com/undefined7887/harmony-backend/internal/infrastructure/log/zap"
	chatrepo "github.com/undefined7887/harmony-backend/internal/repository/chat"
	userrepo "github.com/undefined7887/harmony-backend/internal/repository/user"
	"github.com/undefined7887/harmony-backend/internal/service/auth"
	chatservice "github.com/undefined7887/harmony-backend/internal/service/chat"
	"github.com/undefined7887/harmony-backend/internal/service/jwt"
	userservice "github.com/undefined7887/harmony-backend/internal/service/user"
	"github.com/undefined7887/harmony-backend/internal/third_party/centrifugo"
	"github.com/undefined7887/harmony-backend/internal/third_party/google"
	"github.com/undefined7887/harmony-backend/internal/transport"
	"github.com/undefined7887/harmony-backend/internal/transport/auth"
	chattransport "github.com/undefined7887/harmony-backend/internal/transport/chat"
	usertransport "github.com/undefined7887/harmony-backend/internal/transport/user"
)

func NewApp() *fx.App {
	return fx.New(
		fx.NopLogger,

		config.Module,

		// Infrastructure
		zaplog.Module,
		mongodatabase.Module,

		// Third party
		google.Module,
		centrifugo.Module,

		// Repositories
		userrepo.Module,
		chatrepo.Module,

		// Services
		jwtservice.Module,
		authservice.Module,
		userservice.Module,
		chatservice.Module,

		// Transport
		transport.Module,
		authtransport.Module,
		usertransport.Module,
		chattransport.Module,
	)
}

func FxLogger() fx.Option {
	return fx.WithLogger(func(config *config.App) fxevent.Logger {
		if config.Development {
			return &fxevent.ConsoleLogger{W: os.Stdout}
		}

		return fxevent.NopLogger
	})
}
