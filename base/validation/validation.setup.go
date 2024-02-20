package validation

import (
	"github.com/backend-timedoor/gtimekeeper-framework/container"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

const ContainerName string = "validation"

func New() *Validation {
	v := &Validation{Validator: validator.New()}
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(v.Validator, trans)
	v.Trans = trans

	container.Set(ContainerName, v)

	return v
}
