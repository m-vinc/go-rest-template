package application

import (
	"context"
	"net/http"

	"mpj/internal/interfaces"
	"mpj/internal/models"
	"mpj/pkg/openapi/v1"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/gorilla/websocket"
	"go.uber.org/zap/zapcore"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type LiveController struct {
	cfg *models.ConfigApplication

	gatewayService interfaces.IGateway
	loggerService  interfaces.ILoggerService
}

func NewLiveController(gatewayService interfaces.IGateway, cfg *models.ConfigApplication, loggerService interfaces.ILoggerService) *LiveController {
	controller := &LiveController{
		cfg:            cfg,
		gatewayService: gatewayService,
		loggerService:  loggerService,
	}
	return controller
}

func (controller *LiveController) GatewayConnect(ctx context.Context, request openapi.GatewayConnectRequestObject) (openapi.GatewayConnectResponseObject, error) {
	ectx := middleware.GetEchoContext(ctx)

	conn, err := wsupgrader.Upgrade(ectx.Response(), ectx.Request(), nil)
	if err != nil {
		controller.loggerService.Logf(zapcore.ErrorLevel, "Failed to set websocket upgrade: %+v\n", err)
		return nil, err
	}

	controller.gatewayService.RegisterClient(conn)
	return nil, nil
}
