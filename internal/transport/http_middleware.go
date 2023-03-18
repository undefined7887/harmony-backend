package transport

import (
	"github.com/gin-contrib/cors"
	"github.com/undefined7887/harmony-backend/internal/config"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	zaplog "github.com/undefined7887/harmony-backend/internal/infrastructure/log/zap"
	"github.com/undefined7887/harmony-backend/internal/util"
	httputil "github.com/undefined7887/harmony-backend/internal/util/http"
)

func NewHttpLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := uuid.NewString()
		requestTimestamp := time.Now()

		logger := logger.With(zap.String("request_id", requestID))

		// Packing logger to the next middlewares
		zaplog.PackLoggerGin(ctx, logger)

		// Setting X-Request-ID header
		ctx.Header(httputil.HeaderXRequestID, requestID)

		ctx.Next()

		logger = logger.With(
			zap.String("request_method", ctx.Request.Method),
			zap.String("request_url", ctx.Request.URL.String()),
			zap.String("request_status", httputil.FullStatus(ctx.Writer.Status())),
			zap.Duration("request_time", time.Since(requestTimestamp)),
		)

		switch ctx.Writer.Status() {
		case http.StatusInternalServerError:
			logger.Error("request processing error", zap.Errors("errors", toErrorsSlice(ctx.Errors)))

		default:
			logger.Info("request processed")
		}
	}
}

func CORSMiddleware(httpConfig *config.Http) gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()

	// Overriding default settings
	corsConfig.AllowOrigins = httpConfig.CorsAllowOrigins
	corsConfig.AllowCredentials = httpConfig.CorsAllowCredentials

	return cors.New(corsConfig)
}

func toErrorsSlice(errors []*gin.Error) []error {
	return util.Map(errors, func(item *gin.Error) error {
		return item
	})
}
