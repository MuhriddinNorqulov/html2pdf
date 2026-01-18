package middlewares

import (
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	application "slib.uz/src/application/response"
	"slib.uz/src/application/usecase/logusecases"
	"slib.uz/src/core/response"
	"slib.uz/src/core/utils"
	"slib.uz/src/presentation/app/context"
	//"slib.uz/src/core/utils"
)

type ResponseMiddleware struct {
	adminAlert *logusecases.AdminAlertUseCase
}

func NewResponseMiddleware(adminAlert *logusecases.AdminAlertUseCase) *ResponseMiddleware {
	return &ResponseMiddleware{adminAlert: adminAlert}
}

func (this *ResponseMiddleware) Call(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		err := next(c)

		code := 0

		if err != nil {
			println(err.Error())

			var resp *response.Response
			var paymeErr *application.PaymeError
			var paymeResp *application.PaymeSuccessResponse

			if errors.As(err, &resp) {

				this.log(c, resp.Status)
				return c.JSON(resp.Status, resp.ToJsonResponse())

			} else if errors.As(err, &paymeErr) || errors.As(err, &paymeResp) {
				return c.JSON(200, err)

			} else {
				code = 500
			}
		}

		this.log(c, code)

		return err
	}
}

func (this *ResponseMiddleware) log(ctx echo.Context, _code int) {

	c := ctx.(*context.Context)

	code := ctx.Response().Status
	if _code > 0 {
		code = _code
	}

	if code < 400 || code == 401 || code == 403 || code == 404 {
		return
	}

	hub := sentryecho.GetHubFromContext(ctx)
	if hub == nil {
		return
	}

	req := ctx.Request()

	// Kontekstdan olingan ma'lumotlar (MW lar qo'ygan)
	reqQuery, _ := ctx.Get("req.query").(map[string]string)
	reqBody, _ := ctx.Get("req.body").(string)

	if reqBody == "" {
		if rec, ok := ctx.Get("req.rec").(*RespRecorder); ok && rec != nil {
			reqBody = PrettyJSON(rec.Body())
		} else {
			reqBody = "{}" // Agar body bo'lmasa, bo'sh JSON
		}
	}

	respBody := ""
	if rec, ok := ctx.Get("resp.rec").(*RespRecorder); ok && rec != nil {
		respBody = PrettyJSON(rec.Body())
	}

	_request := PrettyJSON(utils.JsonStringify(map[string]interface{}{
		"method":  req.Method,
		"path":    ctx.Path(),
		"query":   reqQuery,
		"body":    reqBody,
		"user":    c.User,
		"headers": SanitizeHeaders(req.Header),
	}))

	_response := PrettyJSON(utils.JsonStringify(map[string]interface{}{
		"status": code,
		"body":   respBody,
	}))

	hub.WithScope(func(s *sentry.Scope) {

		if code >= 500 {
			s.SetLevel(sentry.LevelError)
		} else {
			s.SetLevel(sentry.LevelInfo)
		}

	})

	hub.CaptureException(fmt.Errorf("HTTP %d %s %s", code, req.Method, ctx.Path()))

	this.adminAlert.Execute(fmt.Sprintf("HTTP %d %s %s\n\n\n\nrequest: %s\n\n\n\nresponse:\n %s", code, req.Method, ctx.Path(), _request, _response))

}
