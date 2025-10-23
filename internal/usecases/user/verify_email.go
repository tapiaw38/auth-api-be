package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

type (
	VerifyEmailUsecase interface {
		Execute(context.Context, string) (string, apperrors.ApplicationError)
	}

	verifyEmailUsecase struct {
		contextFactory appcontext.Factory
	}
)

func NewVerifyEmailUsecase(contextFactory appcontext.Factory) VerifyEmailUsecase {
	return &verifyEmailUsecase{
		contextFactory: contextFactory,
	}
}

func (u *verifyEmailUsecase) Execute(ctx context.Context, token string) (string, apperrors.ApplicationError) {
	app := u.contextFactory()

	user, err := app.Repositories.User.Get(
		ctx,
		user.GetFilterOptions{
			VerifiedEmailToken: token,
		},
	)
	if err != nil {
		return "", apperrors.NewApplicationError(mappings.UserVerifyEmailInvalidTokenError, err)
	}

	if time.Now().After(user.VerifiedEmailTokenExpiry) {
		return "", apperrors.NewApplicationError(mappings.UserVerifyEmailInvalidTokenError, errors.New("token expired"))
	}

	user.VerifiedEmail = true

	if _, err = app.Repositories.User.Update(
		ctx,
		user.ID,
		user,
	); err != nil {
		return "", apperrors.NewApplicationError(mappings.UserVerifyEmailUpdateError, err)
	}

	redirectURL := fmt.Sprintf(
		"%s/",
		app.ConfigService.GCPConfig.OAuth2Config.FrontendURL,
	)

	return redirectURL, nil
}
