package validation

import (
	"github.com/backend-timedoor/gtimekeeper-framework/container"
	"github.com/backend-timedoor/gtimekeeper-framework/utils/helper/exceptions"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"reflect"
	"strings"
)

const ContainerName string = "validation"

func New() *Validation {
	validator_ := validator.New()
	validator_.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	v := &Validation{
		Error:     &exceptions.GTimeError{},
		Validator: validator_,
	}
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(v.Validator, trans)
	v.Trans = trans

	container.Set(ContainerName, v)

	return v
}
