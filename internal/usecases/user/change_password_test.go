package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories"
	user_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	mock_user "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user/mocks"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	"github.com/tapiaw38/auth-api-be/internal/platform/auth"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	usecase "github.com/tapiaw38/auth-api-be/internal/usecases/user"
	"go.uber.org/mock/gomock"
)

func TestChangePasswordUsecase(t *testing.T) {
	type fields struct {
		repository *mock_user.MockRepository
	}

	hashedPassword, _ := auth.HashedPassword("oldpassword")

	tests := map[string]struct {
		input       usecase.ChangePasswordInput
		username    string
		prepare     func(f *fields)
		expectedErr apperrors.ApplicationError
	}{
		"successful password change": {
			input: usecase.ChangePasswordInput{
				OldPassword: "oldpassword",
				NewPassword: "NewPassword123!",
			},
			username: "testuser",
			prepare: func(f *fields) {
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{Username: "testuser"}).Return(&domain.User{
					ID:       "user-123",
					Username: "testuser",
					Password: string(hashedPassword),
				}, nil)
				f.repository.EXPECT().ChangePassword(gomock.Any(), "user-123", gomock.Any()).Return(nil)
			},
			expectedErr: nil,
		},
		"user not found - error from repository": {
			input:    usecase.ChangePasswordInput{},
			username: "non-existent-user",
			prepare: func(f *fields) {
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{Username: "non-existent-user"}).Return(nil, apperrors.NewApplicationError(mappings.UserGetQueryError, nil))
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserGetQueryError, nil),
		},
		"user not found - nil user": {
			input:    usecase.ChangePasswordInput{},
			username: "non-existent-user",
			prepare: func(f *fields) {
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{Username: "non-existent-user"}).Return(nil, nil)
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserGetNotFoundError, nil),
		},
		"incorrect old password": {
			input: usecase.ChangePasswordInput{
				OldPassword: "wrongpassword",
				NewPassword: "NewPassword123!",
			},
			username: "testuser",
			prepare: func(f *fields) {
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{Username: "testuser"}).Return(&domain.User{
					ID:       "user-123",
					Username: "testuser",
					Password: string(hashedPassword),
				}, nil)
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserChangePasswordInvalidOldPasswordError, errors.New("invalid credentials")),
		},
		"new password same as old password": {
			input: usecase.ChangePasswordInput{
				OldPassword: "oldpassword",
				NewPassword: "oldpassword",
			},
			username: "testuser",
			prepare: func(f *fields) {
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{Username: "testuser"}).Return(&domain.User{
					ID:       "user-123",
					Username: "testuser",
					Password: string(hashedPassword),
				}, nil)
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserChangePasswordSamePasswordError, nil),
		},
		"new password too short": {
			input: usecase.ChangePasswordInput{
				OldPassword: "oldpassword",
				NewPassword: "Pass1!",
			},
			username: "testuser",
			prepare: func(f *fields) {
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{Username: "testuser"}).Return(&domain.User{
					ID:       "user-123",
					Username: "testuser",
					Password: string(hashedPassword),
				}, nil)
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserChangePasswordWeakPasswordError, errors.New("password must be at least 8 characters long")),
		},
		"new password missing uppercase": {
			input: usecase.ChangePasswordInput{
				OldPassword: "oldpassword",
				NewPassword: "password123!",
			},
			username: "testuser",
			prepare: func(f *fields) {
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{Username: "testuser"}).Return(&domain.User{
					ID:       "user-123",
					Username: "testuser",
					Password: string(hashedPassword),
				}, nil)
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserChangePasswordWeakPasswordError, errors.New("password must contain at least one uppercase letter")),
		},
		"new password missing lowercase": {
			input: usecase.ChangePasswordInput{
				OldPassword: "oldpassword",
				NewPassword: "PASSWORD123!",
			},
			username: "testuser",
			prepare: func(f *fields) {
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{Username: "testuser"}).Return(&domain.User{
					ID:       "user-123",
					Username: "testuser",
					Password: string(hashedPassword),
				}, nil)
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserChangePasswordWeakPasswordError, errors.New("password must contain at least one lowercase letter")),
		},
		"new password missing number": {
			input: usecase.ChangePasswordInput{
				OldPassword: "oldpassword",
				NewPassword: "Password!",
			},
			username: "testuser",
			prepare: func(f *fields) {
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{Username: "testuser"}).Return(&domain.User{
					ID:       "user-123",
					Username: "testuser",
					Password: string(hashedPassword),
				}, nil)
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserChangePasswordWeakPasswordError, errors.New("password must contain at least one number")),
		},
		"new password missing special character": {
			input: usecase.ChangePasswordInput{
				OldPassword: "oldpassword",
				NewPassword: "Password123",
			},
			username: "testuser",
			prepare: func(f *fields) {
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{Username: "testuser"}).Return(&domain.User{
					ID:       "user-123",
					Username: "testuser",
					Password: string(hashedPassword),
				}, nil)
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserChangePasswordWeakPasswordError, errors.New("password must contain at least one special character (!@#$%^&*()_+-=[]{}|;:,.<>?)")),
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

			uc := usecase.NewChangePasswordUsecase(contextFactory)
			actualErr := uc.Execute(context.Background(), tc.input, tc.username)

			assert.Equal(t, tc.expectedErr, actualErr)
		})
	}
}