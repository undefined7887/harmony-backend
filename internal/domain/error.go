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

func IsError(err1 error, err2 *Error) bool {
	domainErr, ok := err1.(*Error)
	return ok && domainErr.Code == err2.Code
}

func ErrUnauthorized() *Error {
	return &Error{
		StatusCode: http.StatusUnauthorized,

		Code: 1,
		Name: "ERR_UNAUTHORIZED",
	}
}

func ErrBadRequest(err error) *Error {
	return &Error{
		StatusCode: http.StatusBadRequest,

		Code: 2,
		Name: "ERR_BAD_REQUEST",

		// Writing reason
		Message: err.Error(),
	}
}

func ErrForbidden() *Error {
	return &Error{
		StatusCode: http.StatusForbidden,

		Code: 3,
		Name: "ERR_FORBIDDEN",
	}
}

func ErrNotImplemented() *Error {
	return &Error{
		StatusCode: http.StatusNotImplemented,

		Code: 4,
		Name: "ERR_NOT_IMPLEMENTED",
	}
}
