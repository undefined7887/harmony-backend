package userdomain

import (
	"github.com/undefined7887/harmony-backend/internal/domain"
	"net/http"
)

const (
	ErrIndex = 100
)

func ErrUserNotFound() *domain.Error {
	return &domain.Error{
		StatusCode: http.StatusNotFound,

		Code: ErrIndex + 1,
		Name: "ERR_USER_NOT_FOUND",
	}
}

func ErrUserAlreadyExists() *domain.Error {
	return &domain.Error{
		StatusCode: http.StatusConflict,

		Code: ErrIndex + 2,
		Name: "ERR_USER_ALREADY_EXISTS",
	}
}
