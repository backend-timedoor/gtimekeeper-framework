package helper

import (
	"reflect"

	"github.com/jinzhu/copier"
)

func Clone(to any, from any) any {
	copier.Copy(to, from)

	return to
}

func ConvertStructToMap(input interface{}) map[string]interface{} {
	output := make(map[string]interface{})
	val := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)

	// If the input is a pointer, dereference it
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Name

		// Recursively handle nested structs
		if field.Kind() == reflect.Struct {
			output[fieldName] = ConvertStructToMap(field.Interface())
		} else if field.Kind() == reflect.Slice {
			sliceLen := field.Len()
			slice := make([]interface{}, sliceLen)
			for j := 0; j < sliceLen; j++ {
				elem := field.Index(j)
				if elem.Kind() == reflect.Ptr {
					slice[j] = ConvertStructToMap(elem.Interface())
				} else if elem.Kind() == reflect.Struct {
					slice[j] = ConvertStructToMap(elem.Interface())
				} else {
					slice[j] = elem.Interface()
				}
			}
			output[fieldName] = slice
		} else {
			output[fieldName] = field.Interface()
		}
	}

	return output
}
