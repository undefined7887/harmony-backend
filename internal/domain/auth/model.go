package authdomain

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	Channels []string `json:"channels,omitempty"`

	jwt.RegisteredClaims
}
