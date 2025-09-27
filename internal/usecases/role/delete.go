package role

import (
	"context"
	"errors"

	roleRepo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
)

type (
	DeleteUsecase interface {
		Execute(context.Context, string) error
	}

	deleteUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewDeleteUsecase(contextFactory appcontext.Factory) DeleteUsecase {
	return &deleteUsecase{
		contextFactory: contextFactory,
	}
}

func (u *deleteUsecase) Execute(ctx context.Context, id string) error {
	app := u.contextFactory()

	if id == "" {
		return errors.New("role ID is required")
	}

	role, err := app.Repositories.Role.Get(ctx, roleRepo.GetFilterOptions{ID: id})
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	return app.Repositories.Role.Delete(ctx, id)
}
