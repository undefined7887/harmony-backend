package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/lo"
	httputil "github.com/undefined7887/harmony-backend/internal/util/http"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func NewHttpLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := uuid.NewString()
		requestTimestamp := time.Now()

		logger := logger.With(zap.String("request_id", requestID))

		// Packing logger to the next middlewares
		//ctx = zaplog.PackLogger(ctx, logger)

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

func toErrorsSlice(errors []*gin.Error) []error {
	return lo.Map(errors, func(item *gin.Error, index int) error {
		return item
	})
}
