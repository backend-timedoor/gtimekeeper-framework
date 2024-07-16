package helper

import (
	"reflect"
	"strings"

	"github.com/jinzhu/copier"
)

func Clone(to any, from any) any {
	copier.Copy(to, from)

	return to
}

func ConvertGrpcStructToMap(input any) map[string]any {
	output := make(map[string]any)
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
		tags := fieldType.Tag.Get("json")
		split := strings.Split(tags, ",")
		fieldName := split[0]

		if fieldType.Tag.Get("json") != "" {
			if field.Kind() == reflect.Ptr {
				if !field.IsNil() {
					output[fieldName] = ConvertGrpcStructToMap(field.Interface())
				}
			} else if field.Kind() == reflect.Struct {
				output[fieldName] = ConvertGrpcStructToMap(field.Interface())
			} else if field.Kind() == reflect.Slice {
				field := reflect.ValueOf(field.Interface())
				sliceLen := field.Len()
				slice := make([]any, sliceLen)
				for j := 0; j < sliceLen; j++ {
					elem := field.Index(j)
					if elem.Kind() == reflect.Ptr || elem.Kind() == reflect.Struct {
						slice[j] = ConvertGrpcStructToMap(elem.Interface())
					} else {
						slice[j] = elem.Interface()
					}
				}
				output[fieldName] = slice
			} else {
				output[fieldName] = field.Interface()
			}
		}
	}

	return output
}

func ConvertStructToMap(input any) map[string]any {
	output := make(map[string]any)
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
			slice := make([]any, sliceLen)
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

func MergeMapRequest(map1, map2 map[string]string) map[string]string {
	mergedMap := make(map[string]string)

	for key, value := range map1 {
		mergedMap[key] = value
	}

	for key, value := range map2 {
		mergedMap[key] = value
	}

	return mergedMap
}
