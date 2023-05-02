package services

import (
	"context"
	"mpj/internal/interfaces"
	"mpj/internal/models"
	"mpj/pkg/ent"
	"mpj/pkg/ent/predicate"
	"mpj/pkg/ent/user"
	"time"

	"github.com/go-playground/validator/v10"
)

func NewUsersService(database *ent.Client, logger interfaces.ILoggerService) *UsersService {
	svc := &UsersService{
		validator: validator.New(),
		database:  database,
		logger:    logger,
	}

	return svc
}

type UsersService struct {
	validator *validator.Validate
	database  *ent.Client

	cryptoService interfaces.ICryptoService

	logger interfaces.ILoggerService
}

func (svc *UsersService) WrapModel(entity *ent.User) *models.User {
	return &models.User{User: entity}
}

func (svc *UsersService) WrapModels(entities []*ent.User) []*models.User {
	users := make([]*models.User, len(entities))
	for i := 0; i < len(entities); i++ {
		users[i] = svc.WrapModel(entities[i])
	}
	return users
}

func (svc *UsersService) User(ctx context.Context, payload *models.UserQuery) (*models.User, error) {
	err := svc.validator.Struct(payload)
	if err != nil {
		return nil, err
	}

	cond := []predicate.User{}

	if payload.Username != nil {
		cond = append(cond, user.UsernameEQ(*payload.Username))
	} else {
		return nil, models.ErrUserNoGetters
	}

	tx := ent.TxFromContext(ctx)
	if tx == nil {
		tx, err = svc.database.Tx(ctx)
		if err != nil {
			return nil, err
		}
	}

	u, err := tx.User.Query().Where(user.Or(cond...)).First(ctx)
	notfound := ent.IsNotFound(err)
	if err != nil && !notfound {
		return nil, ent.Rollback(ctx, tx, err)
	}

	if notfound {
		return nil, ent.Rollback(ctx, tx, models.ErrUserNotFound)
	}

	return svc.WrapModel(u), nil
}

func (svc *UsersService) CreateUser(ctx context.Context, payload *models.UserCreator) (*models.User, error) {
	err := svc.validator.Struct(payload)
	if err != nil {
		return nil, err
	}

	roles := []string{}
	tx := ent.TxFromContext(ctx)
	if tx == nil {
		tx, err = svc.database.Tx(ctx)
		if err != nil {
			return nil, err
		}
	}

	count, err := tx.User.Query().Count(ctx)
	if err != nil {
		return nil, ent.Rollback(ctx, tx, err)
	}

	if count == 0 {
		roles = append(roles, "admin")
	}

	exist, err := tx.User.Query().Where(user.Username(payload.Username)).Exist(ctx)
	if err != nil {
		return nil, ent.Rollback(ctx, tx, err)
	}

	if exist {
		return nil, ent.Rollback(ctx, tx, models.ErrUsernameAlreadyTaken)
	}

	user, err := tx.User.Create().
		SetUsername(payload.Username).
		SetRoles(roles).
		SetCreatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, ent.Rollback(ctx, tx, err)
	}

	err = ent.Commit(ctx, tx)
	if err != nil {
		return nil, ent.Rollback(ctx, tx, err)
	}

	m := svc.WrapModel(user)

	return m, nil
}

func (svc *UsersService) Users(ctx context.Context, pagination *models.Pagination, payload *models.UserQuery) ([]*models.User, *models.PaginationMetadata, error) {
	err := svc.validator.Struct(payload)
	if err != nil {
		return nil, nil, err
	}

	err = pagination.Default().Validate()
	if err != nil {
		return nil, nil, err
	}

	cond := []predicate.User{}

	if payload.Username != nil {
		cond = append(cond, user.UsernameContainsFold(*payload.Username))
	}

	tx := ent.TxFromContext(ctx)
	if tx == nil {
		tx, err = svc.database.Tx(ctx)
		if err != nil {
			return nil, nil, err
		}
	}

	q := tx.User.Query().Where(user.Or(cond...))

	count, err := q.Count(ctx)
	if err != nil {
		return nil, nil, ent.Rollback(ctx, tx, err)
	}

	us, err := q.
		Offset(int(*pagination.Offset)).
		Limit(int(*pagination.Limit)).
		Order(pagination.Order("created_at", ent.Asc, []string{"created_at", "username"})).
		All(ctx)

	notfound := ent.IsNotFound(err)
	if err != nil && !notfound {
		return nil, nil, ent.Rollback(ctx, tx, err)
	}

	if notfound {
		return nil, nil, ent.Rollback(ctx, tx, models.ErrUserNotFound)
	}

	return svc.WrapModels(us), &models.PaginationMetadata{Offset: int(*pagination.Offset), Limit: int(*pagination.Limit), MaxOffset: count}, nil
}
