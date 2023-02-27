package authtransport

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	jwtservice "github.com/undefined7887/harmony-backend/internal/service/jwt"
	"github.com/undefined7887/harmony-backend/internal/util/http"
	"net/http"
)

func NewHttpAuthMiddleware(jwtService *jwtservice.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := httputil.GetBearerToken(ctx.GetHeader(httputil.HeaderAuthorization))

		if token == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		claims, err := jwtService.Parse(token, &jwt.RegisteredClaims{})
		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, err)

			return
		}

		httputil.SetClaims(ctx, claims.(*jwt.RegisteredClaims))

		// Calling next middleware
		ctx.Next()
	}
}
