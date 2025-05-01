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
	"github.com/tapiaw38/auth-api-be/internal/platform/utils"
)

type (
	LoginUsecase interface {
		Execute(context.Context, LoginInput) (*LoginOutput, error)
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

func (u *loginUsecase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	app := u.contextFactory()

	var findUser *string
	if input.SsoType == string(domain.SsoTypeGoogle) {
		userID, err := googleLogin(ctx, app, input.Code)
		if err != nil {
			return nil, err
		}

		if userID == nil {
			return nil, errors.New("user not found")
		}

		findUser = userID

	} else {
		userID, err := emailAndPasswordLogin(ctx, app, input)
		if err != nil {
			return nil, err
		}

		if userID == nil {
			return nil, errors.New("user not found")
		}

		findUser = userID
	}

	if findUser == nil {
		return nil, errors.New("user not found")
	}

	user, err := app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{
		ID: *findUser,
	})
	if err != nil {
		return nil, err
	}

	token, err := auth.GenerateToken(user, time.Hour*24*7)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		Data:  toUserOutputData(user),
		Token: token,
	}, nil
}

func googleLogin(ctx context.Context, app *appcontext.Context, code string) (*string, error) {
	token, err := app.Integrations.SSO.ExchangeCode(ctx, code)
	if err != nil {
		return nil, err
	}

	userInfo, err := app.Integrations.SSO.GetUserInfo(ctx, token)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	userInsert := domain.User{
		ID:            id.String(),
		FirstName:     userInfo.FirstName,
		LastName:      userInfo.LastName,
		Username:      utils.RandomString(30),
		Email:         userInfo.Email,
		Password:      "",
		Picture:       utils.ToPointer(userInfo.Picture),
		IsActive:      true,
		VerifiedEmail: userInfo.VerifiedEmail,
		CreatedAt:     time.Now(),
	}

	userID, err := app.Repositories.User.Create(ctx, userInsert)
	if err != nil {
		return nil, err
	}

	defaultRole, err := app.Repositories.Role.Get(ctx, role_repo.GetFilterOptions{
		Name: string(domain.RoleUser),
	})

	if _, err = app.Repositories.UserRole.Create(ctx, domain.UserRole{
		UserID: userID,
		RoleID: defaultRole.ID,
	}); err != nil {
		return nil, err
	}

	return &userID, nil
}

func emailAndPasswordLogin(ctx context.Context, app *appcontext.Context, input LoginInput) (*string, error) {
	user, err := app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{
		Email: input.Email,
	})
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("user is not active")
	}

	err = auth.ComparePassword(input.Password, user.Password)
	if err != nil {
		return nil, err
	}

	return &user.ID, nil
}
