package google

import (
	"context"

	googleidtoken "google.golang.org/api/idtoken"

	"github.com/mitchellh/mapstructure"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Auth(ctx context.Context, idtoken string) (*Claims, error) {
	payload, err := googleidtoken.Validate(ctx, idtoken, "")
	if err != nil {
		return nil, err
	}

	return s.decodeClaims(payload)
}

func (s *AuthService) decodeClaims(payload *googleidtoken.Payload) (*Claims, error) {
	claims := &Claims{}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  claims,
	})
	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(payload.Claims); err != nil {
		return nil, err
	}

	return claims, nil
}
