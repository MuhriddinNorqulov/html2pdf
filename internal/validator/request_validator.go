package validator

import (
	"github.com/go-playground/validator/v10"
	"html2pdf/internal/response"
)

type RequestValidator struct {
	validator *validator.Validate
}

func NewRequestValidator(validator *validator.Validate) *RequestValidator {
	return &RequestValidator{validator: validator}
}

func (this *RequestValidator) Validate(i interface{}) error {
	if err := this.validator.Struct(i); err != nil {
		return response.NewFailResponse(400, err.Error())
	}
	if v, ok := i.(Validatable); ok {
		if ok, err := v.Validate(); !ok && err != nil {
			return response.NewFailResponse(400, err.Error())
		}
	}
	return nil
}
