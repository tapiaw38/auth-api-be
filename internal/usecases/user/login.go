package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	role_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role"
	user_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	"github.com/tapiaw38/auth-api-be/internal/platform/auth"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/auth-api-be/internal/platform/utils"
)

type (
	LoginUsecase interface {
		Execute(context.Context, LoginInput) (*LoginOutput, apperrors.ApplicationError)
	}

	loginUsecase struct {
		contextFactory appcontext.Factory
	}

	LoginOutput struct {
		Data  UserOutputData `json:"data"`
		Token string         `json:"token"`
	}

	LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		SsoType  string `json:"sso_type"`
		Code     string `json:"code"`
	}
)

func NewLoginUsecase(contextFactory appcontext.Factory) LoginUsecase {
	return &loginUsecase{
		contextFactory: contextFactory,
	}
}

func (u *loginUsecase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	var findUser *string
	if input.SsoType == string(domain.SsoTypeGoogle) {
		userID, appErr := googleLogin(ctx, app, input)
		if appErr != nil {
			return nil, appErr
		}

		if userID == nil {
			return nil, apperrors.NewApplicationError(mappings.UserLoginUserNotFoundError, errors.New("user not found after google login"))
		}

		findUser = userID

	} else {
		userID, appErr := emailAndPasswordLogin(ctx, app, input)
		if appErr != nil {
			return nil, appErr
		}

		if userID == nil {
			return nil, apperrors.NewApplicationError(mappings.UserLoginUserNotFoundError, errors.New("user not found after email login"))
		}

		findUser = userID
	}

	if findUser == nil {
		return nil, apperrors.NewApplicationError(mappings.UserLoginUserNotFoundError, errors.New("user not found"))
	}

	user, appErr := app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{
		ID: *findUser,
	})
	if appErr != nil {
		return nil, appErr
	}

	if user == nil {
		return nil, apperrors.NewApplicationError(mappings.UserLoginUserNotFoundError, errors.New("user not found"))
	}

	token, err := auth.GenerateToken(user, time.Hour*24*7)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserLoginTokenGenerationError, err)
	}

	return &LoginOutput{
		Data:  toUserOutputData(user),
		Token: token,
	}, nil
}

func googleLogin(ctx context.Context, app *appcontext.Context, input LoginInput) (*string, apperrors.ApplicationError) {
	token, err := app.Integrations.SSO.ExchangeCode(ctx, input.Code)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserLoginGoogleExchangeError, err)
	}

	userInfo, err := app.Integrations.SSO.GetUserInfo(ctx, token)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserLoginGoogleUserInfoError, err)
	}

	user, appErr := app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{
		Email: userInfo.Email,
	})
	if appErr != nil {
		return nil, appErr
	}

	if user == nil {
		id, err := uuid.NewUUID()
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
		}

		encodedString, err := utils.GetEncodedString()
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
		}

		userInsert := domain.User{
			ID:                       id.String(),
			FirstName:                userInfo.FirstName,
			LastName:                 userInfo.LastName,
			Username:                 utils.RandomString(30),
			Email:                    userInfo.Email,
			Password:                 "",
			Picture:                  utils.ToPointer(userInfo.Picture),
			IsActive:                 true,
			VerifiedEmail:            userInfo.VerifiedEmail,
			VerifiedEmailToken:       encodedString,
			VerifiedEmailTokenExpiry: time.Now().Add(time.Hour * 24 * 7),
			AuthMethod:               string(domain.AuthMethodGoogle),
			CreatedAt:                time.Now(),
		}

		createdUserID, appErr := app.Repositories.User.Create(ctx, userInsert)
		if appErr != nil {
			return nil, appErr
		}

		defaultRole, appErr := app.Repositories.Role.Get(ctx, role_repo.GetFilterOptions{
			Name: string(domain.RoleUser),
		})
		if appErr != nil {
			return nil, appErr
		}

		if _, err := app.Repositories.UserRole.Create(ctx, domain.UserRole{
			UserID: createdUserID,
			RoleID: defaultRole.ID,
		}); err != nil {
			return nil, apperrors.NewApplicationError(mappings.UserRegisterAssignRoleError, err)
		}

		return &createdUserID, nil
	}

	if !user.IsActive {
		return nil, apperrors.NewApplicationError(mappings.UserLoginUserNotActiveError, errors.New("user is not active"))
	}
	if !user.VerifiedEmail {
		user.VerifiedEmail = userInfo.VerifiedEmail
	}
	if user.Picture == nil || *user.Picture == "" {
		user.Picture = utils.ToPointer(userInfo.Picture)
	}

	updatedUserID, appErr := app.Repositories.User.Patch(ctx, user.ID, user)
	if appErr != nil {
		return nil, appErr
	}

	return &updatedUserID, nil
}

func emailAndPasswordLogin(ctx context.Context, app *appcontext.Context, input LoginInput) (*string, apperrors.ApplicationError) {
	user, appErr := app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{
		Email: input.Email,
	})
	if appErr != nil {
		return nil, appErr
	}

	if user == nil {
		return nil, apperrors.NewApplicationError(mappings.UserLoginUserNotFoundError, errors.New("user not found with email"))
	}

	if !user.IsActive {
		return nil, apperrors.NewApplicationError(mappings.UserLoginUserNotActiveError, errors.New("user is not active"))
	}

	if user.AuthMethod != string(domain.AuthMethodPassword) && user.AuthMethod != string(domain.AuthMethodHybrid) {
		return nil, apperrors.NewApplicationError(mappings.UserLoginSSOAuthMethodError, errors.New("this account uses SSO authentication"))
	}

	if user.Password == "" {
		return nil, apperrors.NewApplicationError(mappings.UserLoginNoPasswordSetError, errors.New("account has no password set"))
	}

	if err := auth.ComparePassword(input.Password, user.Password); err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserLoginInvalidCredentialsError, err)
	}

	return &user.ID, nil
}
