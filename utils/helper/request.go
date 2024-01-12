package helper

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func MapAndValidate(c echo.Context, model any, request any) error {
	if err := Validate(c, request); err != nil {
		return err
	}

	Clone(model, request)

	return nil
}

func Validate(c echo.Context, request any) error {
	if err := c.Bind(&request); err != nil {
		return ErrorResponse(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(request); err != nil {
		return err
	}

	return nil
}
