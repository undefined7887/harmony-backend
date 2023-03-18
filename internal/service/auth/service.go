package authservice

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/undefined7887/harmony-backend/internal/domain"
	"github.com/undefined7887/harmony-backend/internal/domain/auth"
	"github.com/undefined7887/harmony-backend/internal/domain/user"
	"github.com/undefined7887/harmony-backend/internal/repository"
	jwtservice "github.com/undefined7887/harmony-backend/internal/service/jwt"
	"github.com/undefined7887/harmony-backend/internal/third_party/google"
	randutil "github.com/undefined7887/harmony-backend/internal/util/rand"
)

const (
	MinNicknameTag = 1000
	MaxNicknameTag = 9999
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

func (s *Service) GoogleSignUp(ctx context.Context, nonce, idtoken, nickname string) (authdomain.AuthDTO, error) {
	claims, err := s.googleAuthService.Auth(ctx, idtoken)
	if err != nil {
		return authdomain.AuthDTO{}, authdomain.ErrWrongGoogleToken()
	}

	if !claims.EmailVerified {
		return authdomain.AuthDTO{}, authdomain.ErrEmailNotVerified()
	}

	if claims.Nonce != nonce {
		return authdomain.AuthDTO{}, authdomain.ErrWrongGoogleTokenMalformed()
	}

	now := time.Now()

	user := userdomain.User{
		ID:        domain.ID(),
		Status:    userdomain.StatusOffline,
		Email:     claims.Email,
		Photo:     claims.Picture,
		Nickname:  fmt.Sprintf("%s#%d", nickname, randutil.RandomNumber(MinNicknameTag, MaxNicknameTag)),
		CreatedAt: now,
		UpdatedAt: now,
	}

	inserted, err := s.userRepository.Create(ctx, &user)
	if err != nil {
		return authdomain.AuthDTO{}, err
	}

	if !inserted {
		return authdomain.AuthDTO{}, userdomain.ErrUserAlreadyExists()
	}

	return authdomain.AuthDTO{
		User:      userdomain.MapUserDTO(user),
		UserToken: s.createToken(&user),
	}, nil
}

func (s *Service) GoogleSignIn(ctx context.Context, nonce, idtoken string) (authdomain.AuthDTO, error) {
	claims, err := s.googleAuthService.Auth(ctx, idtoken)
	if err != nil {
		return authdomain.AuthDTO{}, authdomain.ErrWrongGoogleToken()
	}

	if !claims.EmailVerified {
		return authdomain.AuthDTO{}, authdomain.ErrEmailNotVerified()
	}

	if claims.Nonce != nonce {
		return authdomain.AuthDTO{}, authdomain.ErrWrongGoogleTokenMalformed()
	}

	user, err := s.userRepository.GetByEmail(ctx, claims.Email)
	if repository.IsNoDocumentsErr(err) {
		return authdomain.AuthDTO{}, userdomain.ErrUserNotFound()
	}

	if err != nil {
		return authdomain.AuthDTO{}, err
	}

	return authdomain.AuthDTO{
		User:      userdomain.MapUserDTO(user),
		UserToken: s.createToken(&user),
	}, nil
}

func (s *Service) createToken(user *userdomain.User) string {
	return s.jwtService.Create(jwt.RegisteredClaims{
		ID: domain.Token(),

		Issuer:  s.jwtService.Issuer(),
		Subject: user.ID,

		IssuedAt:  s.jwtService.IssuedAt(),
		ExpiresAt: s.jwtService.ExpireAt(),
	})
}
