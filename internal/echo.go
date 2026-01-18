package internal

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"html2pdf/internal/response"
	myvalidator "html2pdf/internal/validator"
)

func NewEcho() *echo.Echo {
	e := echo.New()
	e.Validator = myvalidator.NewRequestValidator(validator.New())
	return e
}

func GetBody[T any](c echo.Context) (*T, error) {
	data := new(T)
	if err := c.Bind(data); err != nil {
		message := fmt.Sprintf("Failed to bind request body: %v", err)
		return nil, response.NewFailResponse(400, message)
	}
	if err := c.Validate(data); err != nil {
		return nil, err
	}
	return data, nil
}
