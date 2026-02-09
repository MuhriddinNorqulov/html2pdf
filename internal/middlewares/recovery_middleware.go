package middlewares

import (
	"errors"
	"fmt"
	"html2pdf/internal/response"
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v4"
)

type RecoveryMiddleware struct {
}

// @inject
func NewRecoveryMiddleware() *RecoveryMiddleware {
	return &RecoveryMiddleware{}
}

func (this *RecoveryMiddleware) Call(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {

				debug.PrintStack()

				var resp *response.Response

				if err, ok := r.(error); ok && errors.As(err, &resp) {
					println("RecoveryMiddleware: ", resp.Status, resp.Message)
					_ = c.JSON(resp.Status, resp)

					if resp.Status >= 500 || resp.Status == 400 {
						println(fmt.Sprintf("Status=%d RecoveryErrorResponse: %v\n%s", resp.Status, r, debug.Stack()))
					}
				} else {
					println("RecoveryError: " + fmt.Sprintf("%v", r) + "\n" + string(debug.Stack()))
					_ = c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("%v", r)})
				}
			}
		}()
		return next(c)
	}
}
