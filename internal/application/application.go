package application

import (
	"context"
	"log"
	"net/http"

	"mpj/internal/interfaces"
	"mpj/internal/models"
	openapi "mpj/pkg/openapi/v1"

	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap/zapcore"

	omiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"
	dochandler "gitlab.com/jamietanna/openapi-doc-http-handler/elements"
)

type ApplicationController struct {
	*echo.Echo

	cfg           *models.ConfigApplication
	loggerService interfaces.ILoggerService

	*LiveController
	*UsersController
}

func New(cfg *models.ConfigApplication, loggerService interfaces.ILoggerService,
	liveController *LiveController,
	usersController *UsersController,
) *ApplicationController {
	application := &ApplicationController{
		cfg:           cfg,
		loggerService: loggerService,

		LiveController:  liveController,
		UsersController: usersController,
	}

	application.init()

	return application
}

func (app *ApplicationController) init() {
	spec, err := openapi.GetSwagger()
	if err != nil {
		log.Fatal(err)
	}

	docHandler, err := dochandler.NewHandler(spec, nil)
	if err != nil {
		log.Fatal(err)
	}

	app.Echo = echo.New()

	app.Use(echozap.ZapLogger(app.loggerService.Logger()))
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	app.GET("/doc", echo.WrapHandler(docHandler))

	validator := omiddleware.OapiRequestValidatorWithOptions(spec, &omiddleware.Options{
		ErrorHandler: func(c echo.Context, err *echo.HTTPError) error {
			return c.JSON(http.StatusUnprocessableEntity, openapi.Error{Code: int32(err.Code), Message: err.Error()})
		},
	})

	v1 := app.Group("/api/v1")
	v1.Use(validator)

	ssi := openapi.NewStrictHandler(app, nil)
	openapi.RegisterHandlers(v1, ssi)
}

func (application *ApplicationController) Run(ctx context.Context) {
	application.loggerService.Logf(zapcore.InfoLevel, "listening on :9090")
	if err := application.Start(application.cfg.Bind); err != nil {
		application.loggerService.Logf(zapcore.ErrorLevel, "%+v", err)
	}
}

func (app *ApplicationController) Panic(ctx context.Context, request openapi.PanicRequestObject) (openapi.PanicResponseObject, error) {
	panic("doh")
}

func (app *ApplicationController) Healthcheck(ctx context.Context, request openapi.HealthcheckRequestObject) (openapi.HealthcheckResponseObject, error) {
	return openapi.Healthcheck200JSONResponse(openapi.HealthStates{
		"apiserver": openapi.HealthState{
			Health:  true,
			Message: "The apiserver is healthy",
		},
	}), nil
}
