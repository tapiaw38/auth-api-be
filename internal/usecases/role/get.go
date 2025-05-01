package role

import (
	"context"

	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
)

type (
	GetUsecase interface {
		Execute(context.Context, GetFilterOptions) (*GetOutput, error)
	}

	getUsecase struct {
		contextFactory appcontext.Factory
	}

	GetOutput struct {
		Data RoleOutputData `json:"data"`
	}

	GetFilterOptions role.GetFilterOptions
)

func NewGetUsecase(contextFactory appcontext.Factory) GetUsecase {
	return &getUsecase{
		contextFactory: contextFactory,
	}
}

func (u *getUsecase) Execute(ctx context.Context, filters GetFilterOptions) (*GetOutput, error) {
	app := u.contextFactory()

	role, err := app.Repositories.Role.Get(ctx, role.GetFilterOptions(filters))
	if err != nil {
		return nil, err
	}

	return &GetOutput{
		Data: toRoleOutputData(*role),
	}, nil
}
