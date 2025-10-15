package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	role_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role"
	user_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/adapters/queue"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/integrations/notification"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	"github.com/tapiaw38/auth-api-be/internal/platform/auth"
	"github.com/tapiaw38/auth-api-be/internal/platform/utils"
)

type (
	RegisterUsecase interface {
		Execute(context.Context, RegisterInput) (*RegisterOutput, error)
	}

	registerUsecase struct {
		contextFactory appcontext.Factory
	}

	RegisterOutput struct {
		Data  UserOutputData `json:"data"`
		Token string         `json:"token"`
	}

	RegisterInput struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) RegisterUsecase {
	return &registerUsecase{
		contextFactory: contextFactory,
	}
}

func (u *registerUsecase) Execute(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	app := u.contextFactory()

	// Validate first name and last name
	if input.FirstName == "" {
		return nil, errors.New("first name is required")
	}
	if input.LastName == "" {
		return nil, errors.New("last name is required")
	}

	// Validate email format
	if err := auth.ValidateEmail(input.Email); err != nil {
		return nil, err
	}

	// Generate username automatically from first name, last name and random string
	generatedUsername := auth.GenerateUsername(input.FirstName, input.LastName)

	user := domain.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Username:  generatedUsername,
		Email:     input.Email,
		Password:  input.Password,
	}

	// Check if email already exists
	existingUser, err := app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{
		Email: user.Email,
	})
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("email already in use")
	}

	err = AddVerifiedEmailToken(&user)
	if err != nil {
		return nil, err
	}

	if err := auth.ValidatePasswordStrength(user.Password); err != nil {
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
	user.AuthMethod = string(domain.AuthMethodPassword)

	userID, err := app.Repositories.User.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	defaultRole, err := app.Repositories.Role.Get(ctx, role_repo.GetFilterOptions{
		Name: string(domain.RoleUser),
	})
	if err != nil {
		return nil, err
	}

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

	emailConfirmation := notification.SendEmailInput{
		To:           createdUser.Email,
		Subject:      "Confirmación de registro",
		TemplateName: "email_verification",
		Variables: map[string]string{
			"name": user.FirstName + " " + user.LastName,
			"link": app.ConfigService.ServerConfig.Host + "/auth/verify-email?token=" + user.VerifiedEmailToken,
		},
	}

	if err = app.Publisher.Publish(queue.TopicSendEmail, emailConfirmation); err != nil {
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
