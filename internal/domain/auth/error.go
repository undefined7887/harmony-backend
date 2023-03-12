package authdomain

import (
	"net/http"

	"github.com/undefined7887/harmony-backend/internal/domain"
)

const (
	ErrIndex = 200
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
