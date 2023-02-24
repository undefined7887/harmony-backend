package authservice

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/undefined7887/harmony-backend/internal/domain/auth"
	"github.com/undefined7887/harmony-backend/internal/domain/user"
	jwtservice "github.com/undefined7887/harmony-backend/internal/service/jwt"
	"github.com/undefined7887/harmony-backend/internal/third_party/google"
	"github.com/undefined7887/harmony-backend/internal/util/crypto"
	"time"
)

type Service struct {
	userRepository userdomain.Repository

	jwtService        *jwtservice.Service
	googleAuthService *google.AuthService
}

func NewService(
	userRepository userdomain.Repository,
	jwtHelper *jwtservice.Service,
	googleAuthService *google.AuthService,
) *Service {
	return &Service{
		userRepository:    userRepository,
		jwtService:        jwtHelper,
		googleAuthService: googleAuthService,
	}
}

func (s *Service) GoogleSignUp(ctx context.Context, idtoken, nickname string) (string, error) {
	claims, err := s.googleAuthService.Auth(ctx, idtoken)
	if err != nil {
		return "", authdomain.ErrWrongGoogleToken()
	}

	if !claims.EmailVerified {
		return "", authdomain.ErrEmailNotVerified()
	}

	user := &userdomain.User{
		ID:        uuid.NewString(),
		Email:     claims.Email,
		Photo:     claims.Picture,
		Nickname:  nickname,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	inserted, err := s.userRepository.Create(ctx, user)
	if err != nil {
		return "", err
	}

	if !inserted {
		return "", authdomain.ErrUserAlreadyExists()
	}

	return s.jwtService.Create(jwt.RegisteredClaims{
		ID: cryptoutil.Token(),

		Issuer:  s.jwtService.Issuer(),
		Subject: user.ID,

		IssuedAt:  s.jwtService.IssuedAt(),
		ExpiresAt: s.jwtService.ExpireAt(),
	}), nil
}

func (s *Service) GoogleSignIn(ctx context.Context, idtoken string) (string, error) {
	claims, err := s.googleAuthService.Auth(ctx, idtoken)
	if err != nil {
		return "", authdomain.ErrWrongGoogleToken()
	}

	user, err := s.userRepository.ReadByEmail(ctx, claims.Email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", authdomain.ErrUserNotFound()
	}

	return s.jwtService.Create(jwt.RegisteredClaims{
		ID: cryptoutil.Token(),

		Issuer:  s.jwtService.Issuer(),
		Subject: user.ID,

		IssuedAt:  s.jwtService.IssuedAt(),
		ExpiresAt: s.jwtService.ExpireAt(),
	}), nil
}
