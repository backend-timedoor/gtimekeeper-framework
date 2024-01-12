package helper

import "github.com/jinzhu/copier"

func Clone(to any, from any) any {
	copier.Copy(to, from)

	return to
}
