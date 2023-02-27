package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	httputil "github.com/undefined7887/harmony-backend/internal/util/http"
	"go.uber.org/zap"
	"net/http"
)

func NewHttpLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		logger := logger.With(
			zap.String("method", ctx.Request.Method),
			zap.String("url", ctx.Request.URL.String()),
			zap.String("status", httputil.FullStatus(ctx.Writer.Status())),
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
