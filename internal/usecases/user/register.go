package user

import (
	"context"
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
	RegisterUsecase interface {
		Execute(context.Context, *domain.User) (*RegisterOutput, error)
	}

	registerUsecase struct {
		contextFactory appcontext.Factory
	}

	RegisterOutput struct {
		Data  UserOutputData `json:"data"`
		Token string         `json:"token"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) RegisterUsecase {
	return &registerUsecase{
		contextFactory: contextFactory,
	}
}

func (u *registerUsecase) Execute(ctx context.Context, user *domain.User) (*RegisterOutput, error) {
	app := u.contextFactory()

	err := AddVerifiedEmailToken(user)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := auth.HashedPassword(user.Password)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	user.ID = id.String()
	user.Password = string(hashedPassword)
	user.IsActive = true
	user.VerifiedEmail = false

	userID, err := app.Repositories.User.Create(ctx, *user)
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

	createdUser, err := app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{
		ID: userID,
	})
	if err != nil {
		return nil, err
	}

	return &RegisterOutput{
		Data: toUserOutputData(createdUser),
	}, nil
}

func AddVerifiedEmailToken(user *domain.User) error {
	encodedString, err := utils.GetEncodedString()
	if err != nil {
		return err
	}

	user.VerifiedEmailToken = encodedString
	user.VerifiedEmailTokenExpiry = time.Now().Add(time.Hour * 24 * 7)

	return nil
}
