package helper

import (
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

type Resp map[string]any

func SuccessResponse(c echo.Context, code int, data any) error {
	return c.JSON(code, Resp{
		"data": data,
	})
}

func ErrorResponse(code int, err any) error {
	return echo.NewHTTPError(code, Resp{
		"message": err,
	})
}

func UnauthenticatedErrorReponse(args ...any) error {
	errMessage := "Unauthenticated"

	if len(args) >= 1 {
		if args[0] != "" {

			switch t := reflect.TypeOf(args[0]); t.Kind() {
			// case reflect.Int:
			// 	fmt.Printf("Value: %v, Type: int\n", v)
			default:
				errMessage = args[0].(string)
			}
		}
	}

	return ErrorResponse(http.StatusUnauthorized, errMessage)
}
