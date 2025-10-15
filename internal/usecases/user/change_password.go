
package user

import (
	"context"
	"errors"

	user_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	"github.com/tapiaw38/auth-api-be/internal/platform/auth"
)

type (
	ChangePasswordUsecase interface {
		Execute(context.Context, ChangePasswordInput) error
	}

	changePasswordUsecase struct {
		contextFactory appcontext.Factory
	}

	ChangePasswordInput struct {
		ID          string
		OldPassword string
		NewPassword string
	}
)

func NewChangePasswordUsecase(contextFactory appcontext.Factory) ChangePasswordUsecase {
	return &changePasswordUsecase{
		contextFactory: contextFactory,
	}
}

func (u *changePasswordUsecase) Execute(ctx context.Context, input ChangePasswordInput) error {
	app := u.contextFactory()

	user, err := app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{
		ID: input.ID,
	})
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	if err := auth.ComparePassword(input.OldPassword, user.Password);
	err != nil {
		return err
	}

	if input.OldPassword == input.NewPassword {
		return errors.New("new password must be different from current password")
	}

	if err := auth.ValidatePasswordStrength(input.NewPassword); err != nil {
		return err
	}

	hashedPassword, err := auth.HashedPassword(input.NewPassword)
	if err != nil {
		return err
	}

	return app.Repositories.User.ChangePassword(ctx, input.ID, string(hashedPassword))
}
