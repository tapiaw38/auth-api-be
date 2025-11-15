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
	UpdateUsecase interface {
		Execute(context.Context, string, UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
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

func (u *updateUsecase) Execute(ctx context.Context, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if id == "" {
		return nil, apperrors.NewApplicationError(mappings.RoleUpdateIDRequiredError, nil)
	}
	if input.Name == "" {
		return nil, apperrors.NewApplicationError(mappings.RoleUpdateNameRequiredError, nil)
	}

	roleName := domain.RoleName(input.Name)
	switch roleName {
	case domain.RoleSuperAdmin, domain.RoleAdmin, domain.RoleUser:
	default:
		return nil, apperrors.NewApplicationError(mappings.RoleUpdateInvalidNameError, nil)
	}

	existingRole, appErr := app.Repositories.Role.Get(ctx, roleRepo.GetFilterOptions{ID: id})
	if appErr != nil {
		return nil, appErr
	}
	if existingRole == nil {
		return nil, apperrors.NewApplicationError(mappings.RoleUpdateNotFoundError, nil)
	}

	role := domain.Role{
		ID:   id,
		Name: roleName,
	}

	updatedID, appErr := app.Repositories.Role.Update(ctx, id, &role)
	if appErr != nil {
		return nil, appErr
	}

	updatedRole, appErr := app.Repositories.Role.Get(ctx, roleRepo.GetFilterOptions{ID: updatedID})
	if appErr != nil {
		return nil, appErr
	}

	return &UpdateOutput{
		Data: toRoleOutputData(*updatedRole),
	}, nil
}
