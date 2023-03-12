package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

const (
	IdTag = "id"
)

var (
	ID       = "^[A-z0-9-_]{38}$"
	IdRegexp = regexp.MustCompile(ID)
)

func validateID(value validator.FieldLevel) bool {
	val, ok := value.Field().Interface().(string)

	return ok && IdRegexp.MatchString(val)
}
