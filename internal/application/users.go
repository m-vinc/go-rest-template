package application

import (
	"context"
	"mpj/internal/application/serializers"
	"mpj/internal/interfaces"
	"mpj/internal/models"
	"mpj/pkg/openapi/v1"
)

type UsersController struct {
	cfg *models.ConfigApplication

	gatewayService interfaces.IGateway
	usersService   interfaces.IUsersService

	loggerService interfaces.ILoggerService
}

func NewUsersController(usersService interfaces.IUsersService, gatewayService interfaces.IGateway, cfg *models.ConfigApplication, loggerService interfaces.ILoggerService) *UsersController {
	controller := &UsersController{
		cfg: cfg,

		gatewayService: gatewayService,
		usersService:   usersService,
		loggerService:  loggerService,
	}
	return controller
}

func (controller *UsersController) Users(ctx context.Context, request openapi.UsersRequestObject) (openapi.UsersResponseObject, error) {
	users, metadata, err := controller.usersService.Users(ctx, &models.Pagination{
		Limit:    request.Params.Limit,
		Offset:   request.Params.Offset,
		OrderBy:  request.Params.OrderBy,
		OrderDir: (*string)(request.Params.OrderDir),
	}, &models.UserQuery{
		Username: request.Params.Q,
	})
	if err != nil {
		return nil, err
	}

	return openapi.Users200JSONResponse{
		Metadata: &openapi.PageMetadata{Total: &metadata.MaxOffset, Offset: &metadata.Offset, Limit: &metadata.Limit},
		Users:    serializers.Users(users),
	}, nil
}

func (controller *UsersController) CreateUser(ctx context.Context, request openapi.CreateUserRequestObject) (openapi.CreateUserResponseObject, error) {
	user, err := controller.usersService.CreateUser(ctx, &models.UserCreator{
		Username:    request.Body.Username,
		Firstname:   &request.Body.FirstName,
		Lastname:    &request.Body.LastName,
		DateOfBirth: request.Body.DateOfBirth,
		Description: request.Body.Description,
	})

	if err != nil {
		return nil, err
	}

	return openapi.CreateUser201JSONResponse(*serializers.User(user)), nil
}
