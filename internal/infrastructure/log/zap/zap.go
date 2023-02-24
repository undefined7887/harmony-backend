package zaplog

import (
	"context"
	"github.com/undefined7887/harmony-backend/internal/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(config *config.Logger, appConfig *config.App) (*zap.Logger, error) {
	level, err := zap.ParseAtomicLevel(config.Level)
	if err != nil {
		return nil, err
	}

	return zap.Config{
		Level:       level,
		Development: appConfig.Development,

		DisableCaller:     true,
		DisableStacktrace: true,

		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "timestamp",
			NameKey:        "logger",
			CallerKey:      zapcore.OmitKey,
			FunctionKey:    zapcore.OmitKey,
			StacktraceKey:  zapcore.OmitKey,
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.MillisDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},

		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}.Build()
}

func NewLifecycleLogger(lifecycle fx.Lifecycle, logger *zap.Logger, appConfig *config.App) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("starting harmony", zap.Bool("development", appConfig.Development))

			return nil
		},

		OnStop: func(ctx context.Context) error {
			logger.Info("stopping harmony")

			return nil
		},
	})
}

type loggerKey struct {
}

func PackLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func UnpackLogger(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(loggerKey{}).(*zap.Logger)
	if !ok {
		return zap.NewNop()
	}

	return logger
}
