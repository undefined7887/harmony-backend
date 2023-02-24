package validation

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

const (
	NicknameTag = "nickname"
)

var (
	Nickname       = "^[A-z0-9_.]{4,30}$"
	NicknameRegexp = regexp.MustCompile(Nickname)
)

func validateNickname(value validator.FieldLevel) bool {
	val, ok := value.Field().Interface().(string)

	return ok && NicknameRegexp.MatchString(val)
}
