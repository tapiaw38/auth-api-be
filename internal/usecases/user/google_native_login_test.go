package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories"
	mockRefreshToken "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/refresh_token/mocks"
	userRepo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	mockUser "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user/mocks"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/integrations"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/integrations/sso"
	mockSSO "github.com/tapiaw38/auth-api-be/internal/adapters/web/integrations/sso/mocks"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	"github.com/tapiaw38/auth-api-be/internal/platform/config"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	usecase "github.com/tapiaw38/auth-api-be/internal/usecases/user"
	"go.uber.org/mock/gomock"
)

func TestGoogleNativeLoginUsecaseExecute(t *testing.T) {
	config.InitConfigService(&config.ConfigurationService{
		ServerConfig: config.ServerConfig{JWTSecret: "test-secret"},
	})

	type fields struct {
		sso           *mockSSO.MockIntegration
		users         *mockUser.MockRepository
		refreshTokens *mockRefreshToken.MockRepository
	}

	picture := "https://example.com/avatar.png"
	googleUser := &domain.User{
		ID:            "user-id",
		Username:      "google-user",
		FirstName:     "Google",
		LastName:      "User",
		Email:         "google@example.com",
		Picture:       &picture,
		IsActive:      true,
		VerifiedEmail: true,
		AuthMethod:    string(domain.AuthMethodGoogle),
		Roles:         []domain.Role{{ID: "role-id", Name: domain.RoleUser}},
	}

	tests := map[string]struct {
		input       usecase.GoogleNativeLoginInput
		prepare     func(*fields)
		expected    *usecase.GoogleNativeLoginOutput
		expectedErr string
	}{
		"when native Google login succeeds": {
			input: usecase.GoogleNativeLoginInput{IDToken: "valid-token"},
			prepare: func(f *fields) {
				f.sso.EXPECT().VerifyIDToken(gomock.Any(), "valid-token").Return(&sso.SocialUser{
					Email:         googleUser.Email,
					FirstName:     googleUser.FirstName,
					LastName:      googleUser.LastName,
					Picture:       picture,
					VerifiedEmail: true,
				}, nil)
				f.users.EXPECT().Get(gomock.Any(), userRepo.GetFilterOptions{Email: googleUser.Email}).Return(googleUser, nil)
				f.users.EXPECT().Patch(gomock.Any(), googleUser.ID, googleUser).Return(googleUser.ID, nil)
				f.users.EXPECT().Get(gomock.Any(), userRepo.GetFilterOptions{ID: googleUser.ID}).Return(googleUser, nil)
				f.refreshTokens.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			expected: &usecase.GoogleNativeLoginOutput{
				Data: usecase.UserOutputData{
					ID:            googleUser.ID,
					Username:      googleUser.Username,
					FirstName:     googleUser.FirstName,
					LastName:      googleUser.LastName,
					Email:         googleUser.Email,
					Picture:       googleUser.Picture,
					IsActive:      true,
					VerifiedEmail: true,
					AuthMethod:    string(domain.AuthMethodGoogle),
					Roles:         []usecase.RoleOutputData{{ID: "role-id", Name: "user"}},
				},
			},
		},
		"when ID token is missing": {
			input:       usecase.GoogleNativeLoginInput{},
			prepare:     func(*fields) {},
			expectedErr: mappings.UserLoginGoogleInvalidIDTokenError.InternalCode,
		},
		"when ID token verification fails": {
			input: usecase.GoogleNativeLoginInput{IDToken: "invalid-token"},
			prepare: func(f *fields) {
				f.sso.EXPECT().VerifyIDToken(gomock.Any(), "invalid-token").Return(nil, errors.New("invalid token"))
			},
			expectedErr: mappings.UserLoginGoogleInvalidIDTokenError.InternalCode,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				sso:           mockSSO.NewMockIntegration(ctrl),
				users:         mockUser.NewMockRepository(ctrl),
				refreshTokens: mockRefreshToken.NewMockRepository(ctrl),
			}
			tc.prepare(&f)

			contextFactory := func(opts ...appcontext.Option) *appcontext.Context {
				return &appcontext.Context{
					Repositories: &repositories.Repositories{
						User:         f.users,
						RefreshToken: f.refreshTokens,
					},
					Integrations: &integrations.Integrations{SSO: f.sso},
				}
			}

			output, appErr := usecase.NewGoogleNativeLoginUsecase(contextFactory).Execute(context.Background(), tc.input)

			if tc.expectedErr != "" {
				assert.Nil(t, output)
				assert.NotNil(t, appErr)
				assert.Equal(t, tc.expectedErr, appErr.InternalCode())
				return
			}

			assert.Nil(t, appErr)
			assert.Equal(t, tc.expected.Data, output.Data)
			assert.NotEmpty(t, output.Token)
			assert.NotEmpty(t, output.RefreshToken)
		})
	}
}
