package user

import (
	"context"
	"errors"

	user_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	"github.com/tapiaw38/auth-api-be/internal/platform/auth"
)

type (
	SetPasswordUsecase interface {
		Execute(context.Context, SetPasswordInput) error
	}

	setPasswordUsecase struct {
		contextFactory appcontext.Factory
	}

	SetPasswordInput struct {
		UserID      string
		NewPassword string
	}
)

func NewSetPasswordUsecase(contextFactory appcontext.Factory) SetPasswordUsecase {
	return &setPasswordUsecase{
		contextFactory: contextFactory,
	}
}

func (u *setPasswordUsecase) Execute(ctx context.Context, input SetPasswordInput) error {
	app := u.contextFactory()

	user, err := app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{
		ID: input.UserID,
	})
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	if user.AuthMethod != string(domain.AuthMethodGoogle) {
		return errors.New("only SSO users can set initial password")
	}

	if err := auth.ValidatePasswordStrength(input.NewPassword); err != nil {
		return err
	}

	hashedPassword, err := auth.HashedPassword(input.NewPassword)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.AuthMethod = string(domain.AuthMethodHybrid)

	_, err = app.Repositories.User.Update(ctx, user.ID, user)
	return err
}
