package authtransport

import (
	"github.com/gin-gonic/gin"
	authdomain "github.com/undefined7887/harmony-backend/internal/domain/auth"
	jwtservice "github.com/undefined7887/harmony-backend/internal/service/jwt"
	"github.com/undefined7887/harmony-backend/internal/util/http"
	"net/http"
)

const (
	ClaimsKey = "user_claims"
)

type TokenQuery struct {
	Token string `query:"token" binding:"jwt"`
}

func NewHttpAuthMiddleware(jwtService *jwtservice.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Trying to get token from header
		token := httputil.GetBearerToken(ctx.GetHeader(httputil.HeaderAuthorization))

		if token == "" {
			// Trying to get token from query
			var query TokenQuery

			if err := ctx.ShouldBindQuery(&query); err == nil {
				token = query.Token
			}
		}

		if token == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		claims, err := jwtService.Parse(token, &authdomain.Claims{})
		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, err)

			return
		}

		ctx.Set(ClaimsKey, claims)

		// Calling next middleware
		ctx.Next()
	}
}

func GetClaims(ctx *gin.Context) *authdomain.Claims {
	return ctx.MustGet(ClaimsKey).(*authdomain.Claims)
}
