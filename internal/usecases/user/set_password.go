package user

import (
	"context"
	"errors"

	user_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	"github.com/tapiaw38/auth-api-be/internal/platform/auth"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

type (
	SetPasswordUsecase interface {
		Execute(context.Context, SetPasswordInput) apperrors.ApplicationError
	}

	setPasswordUsecase struct {
		contextFactory appcontext.Factory
	}

	SetPasswordInput struct {
		Username    string
		NewPassword string
	}
)

func NewSetPasswordUsecase(contextFactory appcontext.Factory) SetPasswordUsecase {
	return &setPasswordUsecase{
		contextFactory: contextFactory,
	}
}

func (u *setPasswordUsecase) Execute(ctx context.Context, input SetPasswordInput) apperrors.ApplicationError {
	app := u.contextFactory()

	user, appErr := app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{
		Username: input.Username,
	})
	if appErr != nil {
		return appErr
	}

	if user == nil {
		return apperrors.NewApplicationError(mappings.UserGetNotFoundError, errors.New("user not found"))
	}

	if user.AuthMethod != string(domain.AuthMethodGoogle) {
		return apperrors.NewApplicationError(mappings.UserSetPasswordNotSSOUserError, errors.New("only SSO users can set initial password"))
	}

	if err := auth.ValidatePasswordStrength(input.NewPassword); err != nil {
		return apperrors.NewApplicationError(mappings.UserSetPasswordWeakPasswordError, err)
	}

	hashedPassword, err := auth.HashedPassword(input.NewPassword)
	if err != nil {
		return apperrors.NewApplicationError(mappings.UserSetPasswordUpdateError, err)
	}

	user.Password = string(hashedPassword)
	user.AuthMethod = string(domain.AuthMethodHybrid)

	if _, appErr = app.Repositories.User.Update(ctx, user.ID, user); appErr != nil {
		return appErr
	}

	return nil
}
