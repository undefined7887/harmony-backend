package authtransport

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	jwtservice "github.com/undefined7887/harmony-backend/internal/service/jwt"
	"github.com/undefined7887/harmony-backend/internal/util/http"
	"net/http"
)

func NewHttpAuthMiddleware(jwtService *jwtservice.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := httputil.GetBearerToken(c.GetHeader(httputil.HeaderAuthorization))

		if token == "" {
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		claims, err := jwtService.Parse(token, &jwt.RegisteredClaims{})
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, err)

			return
		}

		c.Set(httputil.KeyClaims, claims)

		// Calling next middleware
		c.Next()
	}
}
