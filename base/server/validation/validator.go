package validation

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/backend-timedoor/gtimekeeper/utils/helper"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type CustomeValidation struct {
	Validator *validator.Validate
	Trans ut.Translator
}

func (v *CustomeValidation) Validate(i interface{}) error {
	messageBag := map[string]any{}

	if err := v.Validator.Struct(i); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			message := fmt.Sprintf("The %s field is %s", e.Field(), e.ActualTag())
			messageBag[strings.ToLower(e.Field())] = message
		}

		return helper.ErrorResponse(http.StatusUnprocessableEntity, map[string]interface{}{
			"message": "Unprocessable Entity",
			"errors":  messageBag,
		})
	}

	return nil
}

func BootCustomValidation() *CustomeValidation {
	v := &CustomeValidation{Validator: validator.New()}
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(v.Validator, trans)
	v.Trans = trans

	return v
}
