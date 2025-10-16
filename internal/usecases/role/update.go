package role

import (
	"context"
	"errors"

	roleRepo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
)

type (
	UpdateUsecase interface {
		Execute(context.Context, string, UpdateInput) (*UpdateOutput, error)
	}

	updateUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateInput struct {
		Name string `json:"name"`
	}

	UpdateOutput struct {
		Data RoleOutputData `json:"data"`
	}
)

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{
		contextFactory: contextFactory,
	}
}

func (u *updateUsecase) Execute(ctx context.Context, id string, input UpdateInput) (*UpdateOutput, error) {
	app := u.contextFactory()

	if id == "" {
		return nil, errors.New("role ID is required")
	}
	if input.Name == "" {
		return nil, errors.New("role name is required")
	}

	roleName := domain.RoleName(input.Name)
	switch roleName {
	case domain.RoleSuperAdmin, domain.RoleAdmin, domain.RoleUser:
	default:
		return nil, errors.New("invalid role name")
	}

	existingRole, err := app.Repositories.Role.Get(ctx, roleRepo.GetFilterOptions{ID: id})
	if err != nil {
		return nil, err
	}
	if existingRole == nil {
		return nil, errors.New("role not found")
	}

	role := domain.Role{
		ID:   id,
		Name: roleName,
	}

	updatedID, err := app.Repositories.Role.Update(ctx, id, &role)
	if err != nil {
		return nil, err
	}

	updatedRole, err := app.Repositories.Role.Get(ctx, roleRepo.GetFilterOptions{ID: updatedID})
	if err != nil {
		return nil, err
	}

	return &UpdateOutput{
		Data: toRoleOutputData(*updatedRole),
	}, nil
}
