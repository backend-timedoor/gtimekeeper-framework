package validation

import (
	"fmt"
	"github.com/backend-timedoor/gtimekeeper-framework/utils/helper/exceptions"
	"github.com/backend-timedoor/gtimekeeper-framework/utils/helper/types/protocol"
	"net/http"
	"reflect"
	"strings"

	"github.com/backend-timedoor/gtimekeeper-framework/base/contracts"
	"github.com/backend-timedoor/gtimekeeper-framework/base/validation/custom"
	"github.com/backend-timedoor/gtimekeeper-framework/utils/helper"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Values() any
	Type() protocol.Protocol
	Rules() map[string]string
	Messages() map[string]string
}

type Validation struct {
	Error     *exceptions.GTimeError
	Validator *validator.Validate
	Trans     ut.Translator
}

func (v *Validation) RegisterCustomValidation(validations []contracts.CustomeValidation) {
	validations = append(validations, []contracts.CustomeValidation{
		&custom.UniqueValidator{},
		&custom.ExistsValidator{},
	}...)

	for _, validation := range validations {
		v.Validator.RegisterValidation(validation.Signature(), validation.Handle)
	}
}

// new function

func (v *Validation) Make(c Validator) (any, error) {
	return v.Check(c.Type(), c.Values(), c.Rules(), c.Messages())
}
func (v *Validation) Check(ptc protocol.Protocol, d any, rules map[string]string, messages ...map[string]string) (any, error) {
	errors := map[string]any{}

	message := map[string]string{}
	if len(messages) > 0 {
		message = messages[0]
	}

	// check d type is struct or not if map dont convert to map
	data := d
	if reflect.TypeOf(d).Kind() != reflect.Map {
		if ptc == protocol.GRPC {
			data = helper.ConvertGrpcStructToMap(d)
		} else {
			data = helper.ConvertStructToMap(d)
		}
	}

	v.validation(data, rules, message, "", "", errors)

	if len(errors) > 0 {
		err := v.Error.Make(http.StatusUnprocessableEntity, &exceptions.ErrorMessage{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Unprocessable Entity",
			Errors:     errors,
		})

		if ptc == protocol.HTTP {
			return data, err.HttpError()
		}

		return data, err.GrpcError()
	}

	return data, nil
}

func (v *Validation) validation(d any, rules map[string]string, messages map[string]string, field string, prefix string, errorsBag map[string]any) {
	rd := reflect.ValueOf(d)

	if rd.Kind() == reflect.Ptr {
		rd = rd.Elem()
	}

	switch rd.Kind() {
	case reflect.Map:
		for _, val := range rd.MapKeys() {
			keyName := val.String()
			if prefix != "" {
				keyName = prefix + "." + val.String()
			}

			v.validation(rd.MapIndex(val).Interface(), rules, messages, val.String(), keyName, errorsBag)
		}
	case reflect.Slice:
		for i := 0; i < rd.Len(); i++ {
			keyName := fmt.Sprintf("%s[%d]", prefix, i)

			v.validation(rd.Index(i).Interface(), rules, messages, field, keyName, errorsBag)
		}
	default:
		rule, exists := rules[helper.ReplaceWithPattern(prefix, `\[\d+\]`, ".*")]
		if exists {
			if err := v.Validator.Var(rd.Interface(), rule); err != nil {
				for _, err := range err.(validator.ValidationErrors) {
					lFullName := strings.ToLower(prefix)
					lFieldName := strings.ToLower(field)

					keyMsg := fmt.Sprintf("%s.%s", helper.ReplaceWithPattern(prefix, `\[\d+\]`, ".*"), err.ActualTag())
					msg, ok := messages[keyMsg]
					if ok {
						errorsBag[lFullName] = msg
					} else {
						errorsBag[lFullName] = fmt.Sprintf("The %s field is %s %s", lFieldName, err.ActualTag(), err.Param())
					}
				}
			}
		}
	}
}

func (v *Validation) Validate(i any) error {
	_, err := v.Check(protocol.HTTP, i, map[string]string{})

	return err
}
