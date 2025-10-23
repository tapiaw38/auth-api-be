package role

import (
	"context"

	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

type (
	GetUsecase interface {
		Execute(context.Context, GetFilterOptions) (*GetOutput, apperrors.ApplicationError)
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

func (u *getUsecase) Execute(ctx context.Context, filters GetFilterOptions) (*GetOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	role, appErr := app.Repositories.Role.Get(ctx, role.GetFilterOptions(filters))
	if appErr != nil {
		return nil, appErr
	}

	if role == nil {
		return nil, apperrors.NewApplicationError(mappings.RoleGetNotFoundError, nil)
	}

	return &GetOutput{
		Data: toRoleOutputData(*role),
	}, nil
}
