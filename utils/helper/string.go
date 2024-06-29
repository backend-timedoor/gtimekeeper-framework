package helper

import (
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

func ToSnakeCase(v string) string {
	return strcase.ToKebab(v)
}

func ToCamelCase(v string) string {
	return strcase.ToLowerCamel(v)
}

func ToPascalCase(v string) string {
	return strcase.ToCamel(v)
}

func Pluralize(word string) string {
	if strings.HasSuffix(word, "y") {
		return word[:len(word)-1] + "ies"
	} else if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "sh") || strings.HasSuffix(word, "ch") || strings.HasSuffix(word, "x") || strings.HasSuffix(word, "z") {
		return word + "es"
	} else if strings.HasSuffix(word, "f") || strings.HasSuffix(word, "fe") {
		if strings.HasSuffix(word, "fe") {
			return word[:len(word)-2] + "ves"
		} else {
			return word[:len(word)-1] + "ves"
		}
	} else {
		return word + "s"
	}
}

func ReplaceWithPattern(s string, pattern string, replacement string) string {
	re := regexp.MustCompile(pattern)

	return re.ReplaceAllString(s, replacement)
}
