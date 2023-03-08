package domain

import (
	"fmt"
	"net/http"
)

type Error struct {
	StatusCode int `json:"-"`

	Code int    `json:"code"`
	Name string `json:"name"`

	Message string `json:"message,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("api error code: %d, name: %s", e.Code, e.Name)
}

func ErrBadRequest(err error) *Error {
	return &Error{
		StatusCode: http.StatusBadRequest,

		Code: 1,
		Name: "ERR_BAD_REQUEST",

		// Writing reason
		Message: err.Error(),
	}
}

func ErrForbidden() *Error {
	return &Error{
		StatusCode: http.StatusForbidden,

		Code: 2,
		Name: "ERR_FORBIDDEN",
	}
}
func ErrNotImplemented() *Error {
	return &Error{
		StatusCode: http.StatusNotImplemented,

		Code: 3,
		Name: "ERR_NOT_IMPLEMENTED",
	}
}
