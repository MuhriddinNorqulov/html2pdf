package middlewares

import (
	"errors"
	"html2pdf/internal/response"

	"github.com/labstack/echo/v4"
)

type ResponseMiddleware struct {
}

// @inject
func NewResponseMiddleware() *ResponseMiddleware {
	return &ResponseMiddleware{}
}

func (this *ResponseMiddleware) Call(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		err := next(c)

		if err != nil {
			println(err.Error())

			var resp *response.Response

			if errors.As(err, &resp) {

				return c.JSON(resp.Status, resp)

			}
		}
		return err
	}
}
