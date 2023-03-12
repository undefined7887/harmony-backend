package chatdomain

import (
	"net/http"

	"github.com/undefined7887/harmony-backend/internal/domain"
)

const (
	ErrIndex = 300
)

func ErrMessageNotFound() *domain.Error {
	return &domain.Error{
		StatusCode: http.StatusNotFound,

		Code: ErrIndex + 1,
		Name: "ERR_MESSAGE(S)_NOT_FOUND",
	}
}

func ErrChatsNotFound() *domain.Error {
	return &domain.Error{
		StatusCode: http.StatusNotFound,

		Code: ErrIndex + 2,
		Name: "ERR_CHATS_NOT_FOUND",
	}
}
