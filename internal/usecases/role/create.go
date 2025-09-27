package role

import (
	"context"
	"errors"

	roleRepo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
)

type (
	CreateUsecase interface {
		Execute(context.Context, CreateInput) (*CreateOutput, error)
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

func (u *createUsecase) Execute(ctx context.Context, input CreateInput) (*CreateOutput, error) {
	app := u.contextFactory()

	if input.Name == "" {
		return nil, errors.New("role name is required")
	}

	roleName := domain.RoleName(input.Name)
	switch roleName {
	case domain.RoleSuperAdmin, domain.RoleAdmin, domain.RoleUser:
	default:
		return nil, errors.New("invalid role name")
	}

	role := domain.Role{
		Name: roleName,
	}

	id, err := app.Repositories.Role.Create(ctx, role)
	if err != nil {
		return nil, err
	}

	createdRole, err := app.Repositories.Role.Get(ctx, roleRepo.GetFilterOptions{ID: id})
	if err != nil {
		return nil, err
	}

	return &CreateOutput{
		Data: toRoleOutputData(*createdRole),
	}, nil
}
