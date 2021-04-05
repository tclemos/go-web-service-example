package api

import (
	"github.com/labstack/echo/v4"
)

func NewError(ctx echo.Context, code int, m string) error {
	err := ctx.JSON(code, Error{
		Message: m,
	})

	return err
}
