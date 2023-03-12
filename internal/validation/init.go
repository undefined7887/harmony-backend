package validation

import (
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	engine := binding.Validator.Engine().(*validator.Validate)

	if err := engine.RegisterValidation(NicknameTag, validateNickname); err != nil {
		panic(fmt.Sprintf("validation: unexpected \"%s\" registration error: %v", NicknameTag, err))
	}

	if err := engine.RegisterValidation(NicknameExtendedTag, validateNicknameExtended); err != nil {
		panic(fmt.Sprintf("validation: unexpected \"%s\" registration error: %v", NicknameExtendedTag, err))
	}

	if err := engine.RegisterValidation(IdTag, validateID); err != nil {
		panic(fmt.Sprintf("validation: unexpected \"%s\" registration error: %v", IdTag, err))
	}
}
