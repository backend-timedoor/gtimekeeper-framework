package validation

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/backend-timedoor/gtimekeeper-framework/base/contracts"
	"github.com/backend-timedoor/gtimekeeper-framework/base/validation/custom"
	"github.com/backend-timedoor/gtimekeeper-framework/utils/helper"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type Validation struct {
	Validator *validator.Validate
	Trans     ut.Translator
}

func (v *Validation) RegisterCustomeValidation(validations []contracts.CustomeValidation) {
	validations = append(validations, []contracts.CustomeValidation{
		&custom.UniqueValidator{},
	}...)

	for _, validation := range validations {
		v.Validator.RegisterValidation(validation.Signature(), validation.Handle)
	}
}

func (v *Validation) Validate(i interface{}) error {
	if err := v.GValidate(i); err != nil {
		return helper.ErrorResponse(http.StatusUnprocessableEntity, map[string]any{
			"message": "Unprocessable Entity",
			"errors":  err,
		})
	}

	return nil
}

func (v *Validation) GValidate(i interface{}) map[string]any {
	messageBag := map[string]any{}

	if err := v.Validator.Struct(i); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			message := fmt.Sprintf("The %s field is %s", e.Field(), e.ActualTag())
			messageBag[strings.ToLower(e.Field())] = message
		}

		return map[string]any{
			"status":  http.StatusUnprocessableEntity,
			"message": "Unprocessable Entity",
			"errors":  messageBag,
		}
	}

	return nil
}
