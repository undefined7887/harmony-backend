package domain

import (
	"fmt"
)

type Error struct {
	StatusCode int `json:"-"`

	Code int    `json:"code"`
	Name string `json:"name"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("api error code: %d, name: %s", e.Code, e.Name)
}
