package helper

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func MapAndValidate(c echo.Context, model any, request any) error {
	if err := Validation(c, request); err != nil {
		return err
	}

	Clone(model, request)

	return nil
}

func Validation(c echo.Context, request any) error {
	if err := c.Bind(&request); err != nil {
		return ErrorReponse(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(request); err != nil {
		return err
	}

	return nil
}
