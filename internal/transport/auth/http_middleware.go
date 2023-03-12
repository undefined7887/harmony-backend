package authtransport

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	authdomain "github.com/undefined7887/harmony-backend/internal/domain/auth"
	jwtservice "github.com/undefined7887/harmony-backend/internal/service/jwt"
	"github.com/undefined7887/harmony-backend/internal/util/http"
)

const (
	claimsKey = "user_claims"
)

type TokenQuery struct {
	Token string `query:"token" binding:"jwt"`
}

func NewHttpAuthMiddleware(jwtService *jwtservice.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, err := ExtractAndValidateToken(ctx, jwtService)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, err)

			return
		}

		ctx.Set(claimsKey, claims)

		// Calling next middleware
		ctx.Next()
	}
}

func ExtractAndValidateToken(ctx *gin.Context, jwtService *jwtservice.Service) (authdomain.Claims, error) {
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
		return authdomain.Claims{}, errors.New("empty token")
	}

	claims := authdomain.Claims{}

	if err := jwtService.Parse(token, &claims); err != nil {
		return authdomain.Claims{}, err
	}

	return claims, nil
}

func GetClaims(ctx *gin.Context) authdomain.Claims {
	return ctx.MustGet(claimsKey).(authdomain.Claims)
}
