package httputil

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

const (
	claimsKey = "claims"
)

const (
	HeaderAuthorization = "Authorization"
)

func FullStatus(statusCode int) string {
	return fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode))
}

func GetBearerToken(header string) string {
	parts := strings.Split(header, " ")

	if len(parts) != 2 {
		return ""
	}

	if parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func SetClaims(ctx *gin.Context, claims *jwt.RegisteredClaims) {
	ctx.Set(claimsKey, claims)
}

func GetClaims(ctx *gin.Context) *jwt.RegisteredClaims {
	return ctx.MustGet(claimsKey).(*jwt.RegisteredClaims)
}
