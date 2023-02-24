package main

import (
	"github.com/undefined7887/harmony-backend/internal/config"
	mongodatabase "github.com/undefined7887/harmony-backend/internal/infrastructure/database/mongo"
	zaplog "github.com/undefined7887/harmony-backend/internal/infrastructure/log/zap"
	userrepo "github.com/undefined7887/harmony-backend/internal/repository/user"
	"github.com/undefined7887/harmony-backend/internal/service/auth"
	"github.com/undefined7887/harmony-backend/internal/service/jwt"
	"github.com/undefined7887/harmony-backend/internal/third_party/centrifugo"
	"github.com/undefined7887/harmony-backend/internal/third_party/google"
	"github.com/undefined7887/harmony-backend/internal/transport"
	"github.com/undefined7887/harmony-backend/internal/transport/auth"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"os"
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

		// Services
		jwtservice.Module,
		authservice.Module,

		// Transport
		transport.Module,
		authtransport.Module,
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
