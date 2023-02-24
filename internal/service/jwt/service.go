package jwtservice

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/undefined7887/harmony-backend/internal/config"
	"github.com/undefined7887/harmony-backend/internal/util/crypto"
	"time"
)

type Service struct {
	config *config.Jwt

	privateKey crypto.Signer
}

func NewService(config *config.Jwt) (*Service, error) {
	privateKey, err := cryptoutil.ReadPrivateKey(config.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	helper := &Service{
		config:     config,
		privateKey: privateKey,
	}

	if helper.signingMethod() == nil {
		return nil, fmt.Errorf("key %T not supported", privateKey)
	}

	return helper, nil
}

func (h *Service) Issuer() string {
	return h.config.Issuer
}

func (h *Service) IssuedAt() *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now())
}

func (h *Service) ExpireAt() *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(h.config.Lifetime))
}

func (h *Service) Create(claims jwt.Claims) string {
	token, err := jwt.
		NewWithClaims(h.signingMethod(), claims).
		SignedString(h.privateKey)

	if err != nil {
		panic(fmt.Sprintf("unexpected jwt signing error: %v", err))
	}

	return token
}

func (h *Service) Parse(token string, claims jwt.Claims) (jwt.Claims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != h.signingMethod() {
			return nil, errors.New("wrong signing method")
		}

		return h.privateKey.Public(), nil
	})

	if !parsedToken.Valid {
		return nil, err
	}

	return parsedToken.Claims, nil
}

func (h *Service) signingMethod() jwt.SigningMethod {
	switch key := h.privateKey.(type) {
	case ed25519.PrivateKey:
		return jwt.SigningMethodEdDSA

	case *ecdsa.PrivateKey:
		switch key.Params().BitSize {
		case 256:
			return jwt.SigningMethodES256
		case 384:
			return jwt.SigningMethodES384
		case 512:
			return jwt.SigningMethodES512
		}
	}

	return nil
}
