package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories"
	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	mock_user "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user/mocks"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	usecase "github.com/tapiaw38/auth-api-be/internal/usecases/user"
	"go.uber.org/mock/gomock"
)

func TestGetUsecase_Execute(t *testing.T) {
	type fields struct {
		repository *mock_user.MockRepository
	}

	validDateTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	phoneNumber := "+1234567890"
	picture := "https://example.com/avatar.jpg"
	address := "123 Main St, City"

	tests := map[string]struct {
		input       user.GetFilterOptions
		prepare     func(f *fields)
		expected    *usecase.GetOutput
		expectedErr error
	}{
		"when getting user by ID successfully": {
			input: user.GetFilterOptions{
				ID: "user-123",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), user.GetFilterOptions{ID: "user-123"}).
					Return(&domain.User{
						ID:                       "user-123",
						FirstName:                "John",
						LastName:                 "Doe",
						Username:                 "johndoe",
						Email:                    "john@example.com",
						Password:                 "hashedpassword",
						PhoneNumber:              &phoneNumber,
						Picture:                  &picture,
						Address:                  &address,
						IsActive:                 true,
						VerifiedEmail:            true,
						VerifiedEmailToken:       "token123",
						VerifiedEmailTokenExpiry: validDateTime.Add(24 * time.Hour),
						PasswordResetToken:       nil,
						PasswordResetTokenExpiry: nil,
						TokenVersion:             1,
						Roles: []domain.Role{
							{ID: "role-1", Name: domain.RoleUser},
						},
						CreatedAt: validDateTime,
						UpdatedAt: validDateTime,
					}, nil)
			},
			expected: &usecase.GetOutput{
				Data: usecase.UserOutputData{
					ID:            "user-123",
					FirstName:     "John",
					LastName:      "Doe",
					Email:         "john@example.com",
					PhoneNumber:   &phoneNumber,
					Picture:       &picture,
					Address:       &address,
					IsActive:      true,
					VerifiedEmail: true,
					TokenVersion:  1,
					Roles: []usecase.RoleOutputData{
						{ID: "role-1", Name: "user"},
					},
				},
			},
		},
		"when getting user by username successfully": {
			input: user.GetFilterOptions{
				Username: "johndoe",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), user.GetFilterOptions{Username: "johndoe"}).
					Return(&domain.User{
						ID:                       "user-456",
						FirstName:                "Jane",
						LastName:                 "Smith",
						Username:                 "johndoe",
						Email:                    "jane@example.com",
						Password:                 "hashedpassword",
						PhoneNumber:              nil,
						Picture:                  nil,
						Address:                  nil,
						IsActive:                 true,
						VerifiedEmail:            false,
						VerifiedEmailToken:       "token456",
						VerifiedEmailTokenExpiry: validDateTime.Add(24 * time.Hour),
						PasswordResetToken:       nil,
						PasswordResetTokenExpiry: nil,
						TokenVersion:             1,
						Roles: []domain.Role{
							{ID: "role-1", Name: domain.RoleAdmin},
						},
						CreatedAt: validDateTime,
						UpdatedAt: validDateTime,
					}, nil)
			},
			expected: &usecase.GetOutput{
				Data: usecase.UserOutputData{
					ID:            "user-456",
					FirstName:     "Jane",
					LastName:      "Smith",
					Email:         "jane@example.com",
					PhoneNumber:   nil,
					Picture:       nil,
					Address:       nil,
					IsActive:      true,
					VerifiedEmail: false,
					TokenVersion:  1,
					Roles: []usecase.RoleOutputData{
						{ID: "role-1", Name: "admin"},
					},
				},
			},
		},
		"when repository returns error": {
			input: user.GetFilterOptions{
				ID: "user-123",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), user.GetFilterOptions{ID: "user-123"}).
					Return(nil, errors.New("database connection error"))
			},
			expectedErr: errors.New("database connection error"),
		},
		"when user not found": {
			input: user.GetFilterOptions{
				ID: "non-existent-user",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), user.GetFilterOptions{ID: "non-existent-user"}).
					Return(nil, nil)
			},
			expectedErr: errors.New("user not found"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repository: mock_user.NewMockRepository(ctrl),
			}

			if tc.prepare != nil {
				tc.prepare(&f)
			}

			contextFactory := func(opts ...appcontext.Option) *appcontext.Context {
				return &appcontext.Context{
					Repositories: &repositories.Repositories{
						User: f.repository,
					},
				}
			}

			uc := usecase.NewGetUsecase(contextFactory)
			actual, actualErr := uc.Execute(context.Background(), usecase.GetFilterOptions(tc.input))

			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedErr, actualErr)
		})
	}
}
