package validation

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/backend-timedoor/gtimekeeper-framework/base/contracts"
	"github.com/backend-timedoor/gtimekeeper-framework/base/validation/custom"
	"github.com/backend-timedoor/gtimekeeper-framework/utils/helper"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type ValidationType string

const (
	ValidationGrpc ValidationType = "grpc"
	ValidationHttp ValidationType = "http"
)

type Validator interface {
	Type() ValidationType
	Rules() map[string]string
}

type Validation struct {
	Validator *validator.Validate
	Trans     ut.Translator
}

type ValidationError struct {
	Code     int   `json:"-"`
	Message  any   `json:"message"`
	Internal error `json:"-"` // Stores the error returned by an external dependency
}

func (e ValidationError) Error() string {
	return fmt.Sprint(e.Message)
}

func (v *Validation) RegisterCustomeValidation(validations []contracts.CustomeValidation) {
	validations = append(validations, []contracts.CustomeValidation{
		&custom.UniqueValidator{},
	}...)

	for _, validation := range validations {
		v.Validator.RegisterValidation(validation.Signature(), validation.Handle)
	}
}

// new function

func (v *Validation) Make(i any, c Validator) error {
	return v.Check(c.Type(), i, c.Rules())
}
func (v *Validation) Check(vType ValidationType, d any, rules map[string]string) error {
	resp := &ValidationError{Code: http.StatusUnprocessableEntity}
	errors := map[string]any{}

	// check d type is struct or not if map dont convert to map
	data := d
	if reflect.TypeOf(d).Kind() != reflect.Map {
		data = helper.ConvertStructToMap(d)
	}

	v.validation(data, rules, "", "", errors)

	if len(errors) > 0 {
		respMessage := helper.Resp{
			"message": "Unprocessable Entity",
			"errors":  errors,
		}

		if vType == ValidationHttp {
			return helper.ErrorResponse(http.StatusUnprocessableEntity, respMessage)
		}

		respMessage["status"] = http.StatusUnprocessableEntity
		resp.Message = respMessage

		return resp
	}

	return nil
}

func (v *Validation) validation(d any, rules map[string]string, field string, prefix string, errorsBag map[string]any) {
	rd := reflect.ValueOf(d)

	switch rd.Kind() {
	case reflect.Map:
		for _, val := range rd.MapKeys() {
			keyName := val.String()
			if prefix != "" {
				keyName = prefix + "." + val.String()
			}
			v.validation(rd.MapIndex(val).Interface(), rules, val.String(), keyName, errorsBag)
		}
	case reflect.Slice:
		for i := 0; i < rd.Len(); i++ {
			keyName := fmt.Sprintf("%s[%d]", prefix, i)
			v.validation(rd.Index(i).Interface(), rules, field, keyName, errorsBag)
		}
	default:
		if rule, exists := rules[helper.ReplaceWithPattern(prefix, `\[\d+\]`, "")]; exists {
			if err := v.Validator.Var(rd.Interface(), rule); err != nil {
				for _, err := range err.(validator.ValidationErrors) {
					lFullName := strings.ToLower(prefix)
					lFieldName := strings.ToLower(field)
					errorsBag[lFullName] = fmt.Sprintf("The %s field is %s %s", lFieldName, err.ActualTag(), err.Param())
				}
			}
		}
	}
}

func (v *Validation) Validate(i any) error {
	return v.Check(ValidationHttp, i, map[string]string{})
}
