package authdomain

import (
	"github.com/undefined7887/harmony-backend/internal/domain"
	"net/http"
)

const (
	ErrIndex = 100
)

func ErrWrongGoogleToken() *domain.Error {
	return &domain.Error{
		StatusCode: http.StatusBadRequest,

		Code: ErrIndex + 1,
		Name: "ERR_WRONG_GOOGLE_TOKEN",
	}
}

func ErrEmailNotVerified() *domain.Error {
	return &domain.Error{
		StatusCode: http.StatusBadRequest,

		Code: ErrIndex + 2,
		Name: "ERR_EMAIL_NOT_VERIFIED",
	}
}

func ErrUserNotFound() *domain.Error {
	return &domain.Error{
		StatusCode: http.StatusBadRequest,

		Code: ErrIndex + 3,
		Name: "ERR_USER_NOT_FOUND",
	}
}

func ErrUserAlreadyExists() *domain.Error {
	return &domain.Error{
		StatusCode: http.StatusBadRequest,

		Code: ErrIndex + 4,
		Name: "ERR_USER_ALREADY_EXISTS",
	}
}
