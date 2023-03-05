package transport

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/undefined7887/harmony-backend/internal/config"
	"github.com/undefined7887/harmony-backend/internal/domain"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
)

type HttpEndpoint interface {
	Register(c *gin.RouterGroup)
}

func HttpBind(ctx *gin.Context, params, request, query any) bool {
	if params != nil && !HttpBindURI(ctx, params) {
		return false
	}

	if request != nil && !HttpBindJSON(ctx, request) {
		return false
	}

	if query != nil && !HttpBindQuery(ctx, query) {
		return false
	}

	return true
}

func HttpBindURI(ctx *gin.Context, request any) bool {
	if err := ctx.ShouldBindUri(request); err != nil {
		HttpHandleError(ctx, domain.ErrBadRequest(err))

		return false
	}

	return true
}

func HttpBindJSON(ctx *gin.Context, request any) bool {
	if err := ctx.ShouldBindJSON(request); err != nil {
		HttpHandleError(ctx, domain.ErrBadRequest(err))

		return false
	}

	return true
}

func HttpBindQuery(ctx *gin.Context, request any) bool {
	if err := ctx.ShouldBindQuery(request); err != nil {
		HttpHandleError(ctx, domain.ErrBadRequest(err))

		return false
	}

	return true
}

func HttpHandleError(ctx *gin.Context, err error) {
	domainErr, ok := err.(*domain.Error)
	if !ok {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	ctx.AbortWithStatusJSON(domainErr.StatusCode, domainErr)
}

type HttpServer struct {
	logger *zap.Logger
	inner  *http.Server
}

func NewHttpServer(
	config *config.Http,
	logger *zap.Logger,
	endpoints []HttpEndpoint,
) *HttpServer {
	server := &HttpServer{
		logger: logger,
	}

	server.inner = &http.Server{
		Addr:    config.Address,
		Handler: server.handler(endpoints),

		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	return server
}

func NewHttpServerRunner(lifecycle fx.Lifecycle, logger *zap.Logger, server *HttpServer) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info("running http server", zap.String("address", server.Address()))

				if err := server.Start(); err != nil {
					logger.Info("failed to start http server", zap.String("address", server.Address()))
				}
			}()

			return nil
		},

		OnStop: func(ctx context.Context) error {
			if err := server.Stop(ctx); err != nil {
				logger.Info("failed to stop http server", zap.Error(err))
			} else {
				logger.Info("stopped http server")
			}

			return nil
		},
	})
}

func (s *HttpServer) Start() error {
	if err := s.inner.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil // Not returning error on successful shutdown
		}

		return err
	}

	return nil
}

func (s *HttpServer) Stop(ctx context.Context) error {
	return s.inner.Shutdown(ctx)
}

func (s *HttpServer) Address() string {
	return s.inner.Addr
}

func (s *HttpServer) handler(endpoints []HttpEndpoint) http.Handler {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	// Global middlewares
	engine.Use(NewHttpLoggerMiddleware(s.logger))

	v1 := engine.Group("/api/v1")
	{
		for _, endpoint := range endpoints {
			endpoint.Register(v1)
		}
	}

	return engine.Handler()
}
