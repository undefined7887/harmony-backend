package calldomain

import (
	"github.com/undefined7887/harmony-backend/internal/domain"
	"net/http"
)

const (
	ErrIndex = 400
)

func ErrCallNotFound() *domain.Error {
	return &domain.Error{
		StatusCode: http.StatusNotFound,

		Code: ErrIndex + 1,
		Name: "ERR_CALL(S)_NOT_FOUND",
	}
}

func ErrCallAlreadyExists() *domain.Error {
	return &domain.Error{
		StatusCode: http.StatusNotFound,

		Code: ErrIndex + 2,
		Name: "ERR_CALL(S)_ALREADY_EXISTS",
	}
}
