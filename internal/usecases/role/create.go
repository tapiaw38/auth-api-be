package role

import (
	"context"
	"errors"

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
		return nil, apperrors.NewApplicationError(mappings.RoleCreateNameRequiredError, errors.New("role name is required"))
	}

	roleName := domain.RoleName(input.Name)
	switch roleName {
	case domain.RoleSuperAdmin, domain.RoleAdmin, domain.RoleUser:
	default:
		return nil, apperrors.NewApplicationError(mappings.RoleCreateInvalidNameError, errors.New("invalid role name"))
	}

	role := domain.Role{
		Name: roleName,
	}

	id, err := app.Repositories.Role.Create(ctx, role)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.RoleCreateQueryError, err)
	}

	createdRole, err := app.Repositories.Role.Get(ctx, roleRepo.GetFilterOptions{ID: id})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.RoleGetQueryError, err)
	}

	return &CreateOutput{
		Data: toRoleOutputData(*createdRole),
	}, nil
}
