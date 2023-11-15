package helper

import "github.com/jinzhu/copier"

func Clone(to any, from any) error {
	copier.Copy(to, from)

	return nil
}
