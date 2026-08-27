package internal

import (
	"fmt"
	"html2pdf/internal/converter"
	"html2pdf/internal/middlewares"
	"html2pdf/internal/validator"

	"github.com/labstack/echo/v4"
)

type App struct {
	echo *echo.Echo

	// middlewares
	responseMiddleware *middlewares.ResponseMiddleware
	recoveryMiddleware *middlewares.RecoveryMiddleware
	pdfConverter       *converter.PdfConverter
}

// @inject
func NewApp(echo *echo.Echo, responseMiddleware *middlewares.ResponseMiddleware, recoveryMiddleware *middlewares.RecoveryMiddleware, pdfConverter *converter.PdfConverter) *App {
	return &App{echo: echo, responseMiddleware: responseMiddleware, recoveryMiddleware: recoveryMiddleware, pdfConverter: pdfConverter}
}

func (this *App) Handle(c echo.Context) error {
	req, err := GetBody[validator.Request](c)
	if err != nil {
		return err
	}
	params := WkhtmlParamsFromQuery(c.QueryParams())
	r, err := this.pdfConverter.ConvertStream(c.Request().Context(), req.Content, params)
	if err != nil {
		return echo.NewHTTPError(502, fmt.Sprintf("convert error: %v", err))
	}
	defer r.Close()

	c.Response().Header().Set(echo.HeaderContentType, "application/pdf")
	c.Response().Header().Set(echo.HeaderContentDisposition, `inline; filename="doc.pdf"`)

	return c.Stream(200, "application/pdf", r)
}

func (this *App) HealthCheck(c echo.Context) error {
	return c.NoContent(200)
}

func (this *App) Init() {

	this.echo.Use(this.responseMiddleware.Call)
	this.echo.Use(this.recoveryMiddleware.Call)

	this.echo.POST("/to-pdf", this.Handle)
	this.echo.GET("/health", this.HealthCheck)
}

func (this *App) Start() {
	this.echo.Logger.Fatal(this.echo.Start(":8700"))
}
