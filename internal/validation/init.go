package validation

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	engine := binding.Validator.Engine().(*validator.Validate)

	if err := engine.RegisterValidation(NicknameTag, validateNickname); err != nil {
		panic(fmt.Sprintf("validation: unexpected \"nickname\" registration error: %v", err))
	}

	if err := engine.RegisterValidation(NicknameExtendedTag, validateNicknameExtended); err != nil {
		panic(fmt.Sprintf("validation: unexpected \"nickname-extended\" registration error: %v", err))
	}
}
