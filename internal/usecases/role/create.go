package role

import (
	"context"

	roleRepo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

type (
	CreateUsecase interface {
		Execute(context.Context, CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		Name string `json:"name"`
	}

	CreateOutput struct {
		Data RoleOutputData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{
		contextFactory: contextFactory,
	}
}

func (u *createUsecase) Execute(ctx context.Context, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if input.Name == "" {
		return nil, apperrors.NewApplicationError(mappings.RoleCreateNameRequiredError, nil)
	}

	roleName := domain.RoleName(input.Name)
	switch roleName {
	case domain.RoleSuperAdmin, domain.RoleAdmin, domain.RoleUser:
	default:
		return nil, apperrors.NewApplicationError(mappings.RoleCreateInvalidNameError, nil)
	}

	role := domain.Role{
		Name: roleName,
	}

	id, appErr := app.Repositories.Role.Create(ctx, role)
	if appErr != nil {
		return nil, appErr
	}

	createdRole, appErr := app.Repositories.Role.Get(ctx, roleRepo.GetFilterOptions{ID: id})
	if appErr != nil {
		return nil, appErr
	}

	return &CreateOutput{
		Data: toRoleOutputData(*createdRole),
	}, nil
}
