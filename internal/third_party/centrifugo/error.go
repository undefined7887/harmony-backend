package centrifugo

import (
	"fmt"
	"net/http"
)

type HttpError struct {
	StatusCode int
}

func (s *HttpError) Error() string {
	return fmt.Sprintf("%s", http.StatusText(s.StatusCode))
}

type ApiError struct {
	Code    int
	Message string
}

func (a *ApiError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", a.Code, a.Message)
}
