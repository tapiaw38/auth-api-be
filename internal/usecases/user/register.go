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
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/auth-api-be/internal/platform/auth"
	"github.com/tapiaw38/auth-api-be/internal/platform/utils"
)

type (
	RegisterUsecase interface {
		Execute(context.Context, RegisterInput) (*RegisterOutput, apperrors.ApplicationError)
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

func (u *registerUsecase) Execute(ctx context.Context, input RegisterInput) (*RegisterOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if input.FirstName == "" {
		return nil, apperrors.NewApplicationError(mappings.UserRegisterFirstNameRequiredError, errors.New("first name is required"))
	}
	if input.LastName == "" {
		return nil, apperrors.NewApplicationError(mappings.UserRegisterLastNameRequiredError, errors.New("last name is required"))
	}

	if err := auth.ValidateEmail(input.Email); err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserRegisterInvalidEmailError, err)
	}

	generatedUsername := auth.GenerateUsername(input.FirstName, input.LastName)

	user := domain.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Username:  generatedUsername,
		Email:     input.Email,
		Password:  input.Password,
	}

	existingUser, appErr := app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{
		Email: user.Email,
	})
	if appErr != nil {
		return nil, apperrors.NewApplicationError(mappings.UserGetQueryError, appErr)
	}

	if existingUser != nil {
		return nil, apperrors.NewApplicationError(mappings.UserRegisterEmailInUseError, errors.New("email already in use"))
	}

	if err := AddVerifiedEmailToken(&user); err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserRegisterTokenGenerationError, err)
	}

	if err := auth.ValidatePasswordStrength(user.Password); err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserRegisterWeakPasswordError, err)
	}

	hashedPassword, err := auth.HashedPassword(user.Password)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserRegisterPasswordHashError, err)
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserRegisterUUIDGenerationError, err)
	}

	user.ID = id.String()
	user.Password = string(hashedPassword)
	user.IsActive = true
	user.VerifiedEmail = false
	user.AuthMethod = string(domain.AuthMethodPassword)

	userID, err := app.Repositories.User.Create(ctx, user)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserRegisterCreateUserError, err)
	}

	defaultRole, appErr := app.Repositories.Role.Get(ctx, role_repo.GetFilterOptions{
		Name: string(domain.RoleUser),
	})
	if appErr != nil {
		return nil, apperrors.NewApplicationError(mappings.UserRegisterGetRoleError, appErr)
	}

	if _, err := app.Repositories.UserRole.Create(ctx, domain.UserRole{
		UserID: userID,
		RoleID: defaultRole.ID,
	}); err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserRegisterAssignRoleError, err)
	}

	createdUser, appErr := app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{
		ID: userID,
	})
	if appErr != nil {
		return nil, apperrors.NewApplicationError(mappings.UserRegisterGetUserError, appErr)
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

	if err := app.Publisher.Publish(queue.TopicSendEmail, emailConfirmation); err != nil {
		return nil, apperrors.NewApplicationError(mappings.UserRegisterEmailSendError, err)
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
