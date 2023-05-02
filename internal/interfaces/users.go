package interfaces

import (
	"context"
	"mpj/internal/models"
)

type IUsersService interface {
	Users(ctx context.Context, pagination *models.Pagination, payload *models.UserQuery) ([]*models.User, *models.PaginationMetadata, error)
	User(ctx context.Context, payload *models.UserQuery) (*models.User, error)

	CreateUser(ctx context.Context, payload *models.UserCreator) (*models.User, error)
}
