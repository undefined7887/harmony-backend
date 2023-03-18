package authtransport

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/undefined7887/harmony-backend/internal/domain"
	authdomain "github.com/undefined7887/harmony-backend/internal/domain/auth"
	zaplog "github.com/undefined7887/harmony-backend/internal/infrastructure/log/zap"
	jwtservice "github.com/undefined7887/harmony-backend/internal/service/jwt"
	"github.com/undefined7887/harmony-backend/internal/transport"
	"github.com/undefined7887/harmony-backend/internal/util/http"
	"go.uber.org/zap"
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
			transport.HttpHandleError(ctx, domain.ErrUnauthorized())

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
		var err error

		// Trying to get token from 'token' cookie
		token, err = ctx.Cookie("token")
		if err != nil {
			zaplog.
				UnpackLogger(ctx).
				Info("failed to get token from cookie", zap.Error(err))
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
