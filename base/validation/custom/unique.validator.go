package custom

import (
	"github.com/go-playground/validator/v10"
)

type UniqueValidator struct{}

func (u *UniqueValidator) Signature() string {
	return "example"
}

func (u *UniqueValidator) Handle(fl validator.FieldLevel) bool {
	// tagParts := strings.Split(fl.Param(), ":")

	// if len(tagParts) != 2 {
	// 	return false
	// }

	// tableName := tagParts[0]
	// fieldName := tagParts[1]
	// fieldValue := fl.Field().String()

	// var count int64
	// app.DB.Table(tableName).Where(fieldName+" = ?", fieldValue).Count(&count)

	return true
}
