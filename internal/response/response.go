package response

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (this *Response) Error() string {
	return this.Message
}

func NewFailResponse(status int, message string) *Response {
	return &Response{Status: status, Message: message}
}
