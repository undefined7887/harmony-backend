package httputil

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	KeyClaims = "claims"
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
