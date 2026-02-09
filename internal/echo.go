package internal

import (
	"fmt"
	"html2pdf/internal/response"
	myvalidator "html2pdf/internal/validator"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// @inject
func NewEcho() *echo.Echo {
	e := echo.New()
	e.Validator = myvalidator.NewRequestValidator(validator.New())

	e.Server.ReadTimeout = 30 * time.Second
	e.Server.WriteTimeout = 180 * time.Second // PDF stream uchun katta bo‘lsin
	e.Server.IdleTimeout = 120 * time.Second

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
