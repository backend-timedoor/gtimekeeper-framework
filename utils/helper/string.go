package helper

import "github.com/iancoleman/strcase"

func ToSnakeCase(v string) string {
	return strcase.ToKebab(v)
}

func ToCamelCase(v string) string {
	return strcase.ToLowerCamel(v)
}

func ToPascalCase(v string) string {
	return strcase.ToCamel(v)
}
