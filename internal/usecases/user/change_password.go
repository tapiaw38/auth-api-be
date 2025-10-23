package user

import (
	"context"

	user_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	"github.com/tapiaw38/auth-api-be/internal/platform/auth"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

type (
	ChangePasswordUsecase interface {
		Execute(context.Context, ChangePasswordInput, string) apperrors.ApplicationError
	}

	changePasswordUsecase struct {
		contextFactory appcontext.Factory
	}

	ChangePasswordInput struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
)

func NewChangePasswordUsecase(contextFactory appcontext.Factory) ChangePasswordUsecase {
	return &changePasswordUsecase{
		contextFactory: contextFactory,
	}
}

func (u *changePasswordUsecase) Execute(ctx context.Context, input ChangePasswordInput, username string) apperrors.ApplicationError {
	app := u.contextFactory()

	user, appErr := app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{
		Username: username,
	})
	if appErr != nil {
		return appErr
	}

	if user == nil {
		return apperrors.NewApplicationError(mappings.UserGetNotFoundError, nil)
	}

	if err := auth.ComparePassword(input.OldPassword, user.Password); err != nil {
		return apperrors.NewApplicationError(mappings.UserChangePasswordInvalidOldPasswordError, err)
	}

	if input.OldPassword == input.NewPassword {
		return apperrors.NewApplicationError(mappings.UserChangePasswordSamePasswordError, nil)
	}

	if err := auth.ValidatePasswordStrength(input.NewPassword); err != nil {
		return apperrors.NewApplicationError(mappings.UserChangePasswordWeakPasswordError, err)
	}

	hashedPassword, err := auth.HashedPassword(input.NewPassword)
	if err != nil {
		return apperrors.NewApplicationError(mappings.UserChangePasswordUpdateError, err)
	}

	if appErr := app.Repositories.User.ChangePassword(ctx, user.ID, string(hashedPassword)); appErr != nil {
		return appErr
	}

	return nil
}
