package validator

type Request struct {
	Content string `json:"content" validate:"required"`
}
