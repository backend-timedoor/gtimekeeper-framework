package helper

import (
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

type Resp map[string]any

func ErrorReponse(code int, err any) error {
	return echo.NewHTTPError(code, err)
}

func UnauthenticatedErrorReponse(args ...any) error {
	message := "Unauthenticated"

	if len(args) >= 1 {
		if args[0] != "" {

			switch t := reflect.TypeOf(args[0]); t.Kind() {
			// case reflect.Int:
			// 	fmt.Printf("Value: %v, Type: int\n", v)
			default:
				message = args[0].(string)
			}
		}
	}

	return ErrorReponse(http.StatusUnauthorized, map[string]string{
		"message": message,
	})
}
