package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

const (
	NicknameTag         = "nickname"
	NicknameExtendedTag = "nickname-extended"
)

var (
	Nickname       = "^[A-z0-9-_.]{4,30}$"
	NicknameRegexp = regexp.MustCompile(Nickname)

	NicknameExtended       = "^[A-z0-9-_.]{4,30}#[0-9]{4}$"
	NicknameExtendedRegexp = regexp.MustCompile(NicknameExtended)
)

func validateNickname(value validator.FieldLevel) bool {
	val, ok := value.Field().Interface().(string)

	return ok && NicknameRegexp.MatchString(val)
}

func validateNicknameExtended(value validator.FieldLevel) bool {
	val, ok := value.Field().Interface().(string)

	return ok && NicknameExtendedRegexp.MatchString(val)
}
