package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	httputil "github.com/undefined7887/harmony-backend/internal/util/http"
	"go.uber.org/zap"
	"net/http"
)

func NewHttpLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		logger := logger.With(
			zap.String("method", c.Request.Method),
			zap.String("url", c.Request.URL.String()),
			zap.String("status", httputil.FullStatus(c.Writer.Status())),
		)

		switch c.Writer.Status() {
		case http.StatusInternalServerError:
			logger.Error("request processing error", zap.Errors("errors", toErrorsSlice(c.Errors)))

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
