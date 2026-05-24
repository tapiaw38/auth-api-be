package user

import (
	"context"
	"errors"
	"time"

	user_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/adapters/queue"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/integrations/notification"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/auth-api-be/internal/platform/utils"
)

type (
	RequestResetPasswordUsecase interface {
		Execute(context.Context, string) (*RequestResetPasswordOutput, apperrors.ApplicationError)
	}

	requestResetPasswordUsecase struct {
		contextFactory appcontext.Factory
	}

	RequestResetPasswordOutput struct {
		Data ResetPasswordOutputData `json:"data"`
	}
)

func NewRequestResetPasswordUsecase(contextFactory appcontext.Factory) RequestResetPasswordUsecase {
	return &requestResetPasswordUsecase{
		contextFactory: contextFactory,
	}
}

func (u *requestResetPasswordUsecase) Execute(ctx context.Context, email string) (*RequestResetPasswordOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	user, appErr := app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{
		Email: email,
	})
	if appErr != nil {
		return nil, apperrors.NewApplicationError(mappings.UserGetQueryError, appErr)
	}

	if user == nil {
		return nil, apperrors.NewApplicationError(mappings.UserRequestResetPasswordUserNotFoundError, errors.New("user email not found"))
	}

	token, err := utils.GetEncodedString()
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserRequestResetPasswordTokenGenerationError, err)
	}

	tokenExpiry := time.Now().Add(time.Hour * 24)

	user.PasswordResetToken = &token
	user.PasswordResetTokenExpiry = &tokenExpiry

	if _, err := app.Repositories.User.Patch(ctx, user.ID, user); err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserRequestResetPasswordUpdateError, err)
	}

	emailResetPassword := notification.SendEmailInput{
		To:           user.Email,
		Subject:      "Restablecer contraseña",
		TemplateName: "reset_password",
		Variables: map[string]string{
			"name": user.FirstName + " " + user.LastName,
			"link": app.ConfigService.ServerConfig.Host + "/auth/reset-password?token=" + token,
		},
	}

	if err := app.Publisher.Publish(queue.TopicSendEmail, emailResetPassword); err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserRequestResetPasswordEmailSendError, err)
	}

	return &RequestResetPasswordOutput{
		Data: ResetPasswordOutputData{
			Email:   user.Email,
			Message: "Password reset token sent",
		},
	}, nil
}
