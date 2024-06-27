package helper

import (
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
		// If the word ends in "y", replace "y" with "ies"
		return word[:len(word)-1] + "ies"
	} else if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "sh") || strings.HasSuffix(word, "ch") || strings.HasSuffix(word, "x") || strings.HasSuffix(word, "z") {
		// For words ending in s, sh, ch, x, z, add "es"
		return word + "es"
	} else if strings.HasSuffix(word, "f") || strings.HasSuffix(word, "fe") {
		// If the word ends in "f" or "fe", replace "f" or "fe" with "ves"
		if strings.HasSuffix(word, "fe") {
			return word[:len(word)-2] + "ves"
		} else {
			return word[:len(word)-1] + "ves"
		}
	} else {
		// For most cases, just add "s"
		return word + "s"
	}
}
