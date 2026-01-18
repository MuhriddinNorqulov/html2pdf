package internal

import (
	"github.com/labstack/echo/v4"
	"html2pdf/internal/middlewares"
)

type App struct {
	echo *echo.Echo

	// middlewares
	responseMiddleware *middlewares.ResponseMiddleware
	recoveryMiddleware *middlewares.RecoveryMiddleware
}

func NewApp(echo *echo.Echo, responseMiddleware *middlewares.ResponseMiddleware, recoveryMiddleware *middlewares.RecoveryMiddleware) *App {
	return &App{echo: echo, responseMiddleware: responseMiddleware, recoveryMiddleware: recoveryMiddleware}
}

//func (this *App) Handle(c echo.Context) error {
//	req, err := GetBody[schema.Request](c)
//	if err != nil {
//		return err
//	}
//
//
//}
