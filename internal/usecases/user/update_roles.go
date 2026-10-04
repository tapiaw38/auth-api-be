package user

import (
	"context"
	"errors"
	"slices"

	role_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role"
	user_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

var assignableRoles = []domain.RoleName{domain.RoleUser, domain.RoleAdmin}

type (
	UpdateRolesUsecase interface {
		Execute(context.Context, UpdateRolesInput) (*UpdateRolesOutput, apperrors.ApplicationError)
	}

	updateRolesUsecase struct {
		contextFactory appcontext.Factory
	}

	UpdateRolesInput struct {
		TargetID      string
		ActorUsername string
		Roles         []string
	}

	UpdateRolesOutput struct {
		Data UserOutputData `json:"data"`
	}
)

func NewUpdateRolesUsecase(contextFactory appcontext.Factory) UpdateRolesUsecase {
	return &updateRolesUsecase{contextFactory: contextFactory}
}

func (u *updateRolesUsecase) Execute(ctx context.Context, input UpdateRolesInput) (*UpdateRolesOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	actor, appErr := app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{Username: input.ActorUsername})
	if appErr != nil {
		return nil, appErr
	}
	if actor == nil {
		return nil, apperrors.NewApplicationError(mappings.AuthUnauthorizedError,
			errors.New("authenticated user not found"))
	}
	if actor.ID == input.TargetID {
		return nil, apperrors.NewApplicationError(mappings.UserRolesSelfUpdateError,
			errors.New("user attempted to change their own roles"))
	}

	wanted := make([]domain.RoleName, 0, len(input.Roles))
	for _, name := range input.Roles {
		roleName := domain.RoleName(name)
		if !slices.Contains(assignableRoles, roleName) {
			return nil, apperrors.NewApplicationError(mappings.UserRolesNotAssignableError,
				errors.New("role cannot be assigned through the API: "+name))
		}
		if !slices.Contains(wanted, roleName) {
			wanted = append(wanted, roleName)
		}
	}

	target, appErr := app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{ID: input.TargetID})
	if appErr != nil {
		return nil, appErr
	}

	if target == nil {
		return nil, apperrors.NewApplicationError(mappings.UserUpdateNotFoundError, nil)
	}

	current := make([]domain.RoleName, 0, len(target.Roles))
	for _, r := range target.Roles {
		if slices.Contains(assignableRoles, r.Name) {
			current = append(current, r.Name)
		}
	}

	for _, name := range wanted {
		if slices.Contains(current, name) {
			continue
		}
		if appErr := u.grant(ctx, app, input.TargetID, name); appErr != nil {
			return nil, appErr
		}
	}

	for _, name := range current {
		if slices.Contains(wanted, name) {
			continue
		}
		if appErr := u.revoke(ctx, app, input.TargetID, name); appErr != nil {
			return nil, appErr
		}
	}

	if appErr := app.Repositories.User.IncrementTokenVersion(ctx, input.TargetID); appErr != nil {
		return nil, appErr
	}

	updated, appErr := app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{ID: input.TargetID})
	if appErr != nil {
		return nil, appErr
	}

	return &UpdateRolesOutput{Data: toUserOutputData(updated)}, nil
}

func (u *updateRolesUsecase) grant(ctx context.Context, app *appcontext.Context, userID string, name domain.RoleName) apperrors.ApplicationError {
	role, appErr := u.findRole(ctx, app, name)
	if appErr != nil {
		return appErr
	}

	if _, err := app.Repositories.UserRole.Create(ctx, domain.UserRole{UserID: userID, RoleID: role.ID}); err != nil {
		return apperrors.NewApplicationError(mappings.UserRolesUpdateError, err)
	}

	return nil
}

func (u *updateRolesUsecase) revoke(ctx context.Context, app *appcontext.Context, userID string, name domain.RoleName) apperrors.ApplicationError {
	role, appErr := u.findRole(ctx, app, name)
	if appErr != nil {
		return appErr
	}

	if _, err := app.Repositories.UserRole.Delete(ctx, userID, role.ID); err != nil {
		return apperrors.NewApplicationError(mappings.UserRolesUpdateError, err)
	}

	return nil
}

func (u *updateRolesUsecase) findRole(ctx context.Context, app *appcontext.Context, name domain.RoleName) (*domain.Role, apperrors.ApplicationError) {
	role, appErr := app.Repositories.Role.Get(ctx, role_repo.GetFilterOptions{Name: string(name)})
	if appErr != nil {
		return nil, appErr
	}
	if role == nil {
		return nil, apperrors.NewApplicationError(mappings.UserRolesUpdateError,
			errors.New("role not found: "+string(name)))
	}

	return role, nil
}
