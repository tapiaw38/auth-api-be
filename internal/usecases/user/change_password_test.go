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
		prepare     func(f *fields)
		expectedErr error
	}{
		"successful password change": {
			input: usecase.ChangePasswordInput{
				ID:          "user-123",
				OldPassword: "oldpassword",
				NewPassword: "NewPassword123!",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{ID: "user-123"}).Return(&domain.User{
					ID:       "user-123",
					Password: string(hashedPassword),
				}, nil)
				f.repository.EXPECT().ChangePassword(gomock.Any(), "user-123", gomock.Any()).Return(nil)
			},
			expectedErr: nil,
		},
		"user not found - error from repository": {
			input: usecase.ChangePasswordInput{
				ID: "non-existent-user",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{ID: "non-existent-user"}).Return(nil, errors.New("not found"))
			},
			expectedErr: errors.New("not found"),
		},
		"user not found - nil user": {
			input: usecase.ChangePasswordInput{
				ID: "non-existent-user",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{ID: "non-existent-user"}).Return(nil, nil)
			},
			expectedErr: errors.New("user not found"),
		},
		"incorrect old password": {
			input: usecase.ChangePasswordInput{
				ID:          "user-123",
				OldPassword: "wrongpassword",
				NewPassword: "NewPassword123!",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{ID: "user-123"}).Return(&domain.User{
					ID:       "user-123",
					Password: string(hashedPassword),
				}, nil)
			},
			expectedErr: errors.New("invalid credentials"),
		},
		"new password same as old password": {
			input: usecase.ChangePasswordInput{
				ID:          "user-123",
				OldPassword: "oldpassword",
				NewPassword: "oldpassword",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{ID: "user-123"}).Return(&domain.User{
					ID:       "user-123",
					Password: string(hashedPassword),
				}, nil)
			},
			expectedErr: errors.New("new password must be different from current password"),
		},
		"new password too short": {
			input: usecase.ChangePasswordInput{
				ID:          "user-123",
				OldPassword: "oldpassword",
				NewPassword: "Pass1!",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{ID: "user-123"}).Return(&domain.User{
					ID:       "user-123",
					Password: string(hashedPassword),
				}, nil)
			},
			expectedErr: errors.New("password must be at least 8 characters long"),
		},
		"new password missing uppercase": {
			input: usecase.ChangePasswordInput{
				ID:          "user-123",
				OldPassword: "oldpassword",
				NewPassword: "password123!",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{ID: "user-123"}).Return(&domain.User{
					ID:       "user-123",
					Password: string(hashedPassword),
				}, nil)
			},
			expectedErr: errors.New("password must contain at least one uppercase letter"),
		},
		"new password missing lowercase": {
			input: usecase.ChangePasswordInput{
				ID:          "user-123",
				OldPassword: "oldpassword",
				NewPassword: "PASSWORD123!",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{ID: "user-123"}).Return(&domain.User{
					ID:       "user-123",
					Password: string(hashedPassword),
				}, nil)
			},
			expectedErr: errors.New("password must contain at least one lowercase letter"),
		},
		"new password missing number": {
			input: usecase.ChangePasswordInput{
				ID:          "user-123",
				OldPassword: "oldpassword",
				NewPassword: "Password!",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{ID: "user-123"}).Return(&domain.User{
					ID:       "user-123",
					Password: string(hashedPassword),
				}, nil)
			},
			expectedErr: errors.New("password must contain at least one number"),
		},
		"new password missing special character": {
			input: usecase.ChangePasswordInput{
				ID:          "user-123",
				OldPassword: "oldpassword",
				NewPassword: "Password123",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{ID: "user-123"}).Return(&domain.User{
					ID:       "user-123",
					Password: string(hashedPassword),
				}, nil)
			},
			expectedErr: errors.New("password must contain at least one special character (!@#$%^&*()_+-=[]{}|;:,.<>?)"),
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
			actualErr := uc.Execute(context.Background(), tc.input)

			assert.Equal(t, tc.expectedErr, actualErr)
		})
	}
}