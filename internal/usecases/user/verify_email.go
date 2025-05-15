package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
)

type (
	VerifyEmailUsecase interface {
		Execute(context.Context, string) (string, error)
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

func (u *verifyEmailUsecase) Execute(ctx context.Context, token string) (string, error) {
	app := u.contextFactory()

	user, err := app.Repositories.User.Get(
		ctx,
		user.GetFilterOptions{
			VerifiedEmailToken: token,
		},
	)
	if err != nil {
		return "", err
	}

	if time.Now().After(user.VerifiedEmailTokenExpiry) {
		return "", errors.New("token expired")
	}

	user.VerifiedEmail = true

	if _, err = app.Repositories.User.Update(
		ctx,
		user.ID,
		user,
	); err != nil {
		return "", err
	}

	redirectURL := fmt.Sprintf(
		"%s/",
		app.ConfigService.GCPConfig.OAuth2Config.FrontendURL,
	)

	return redirectURL, nil
}
