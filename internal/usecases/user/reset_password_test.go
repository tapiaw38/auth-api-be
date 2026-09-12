package user_test

import (
	"context"
	"testing"
	"time"

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

func TestResetPasswordUsecase(t *testing.T) {
	type fields struct {
		repository *mock_user.MockRepository
	}

	now := time.Now()
	expiredTime := now.Add(-time.Hour)
	validTime := now.Add(time.Hour)
	hashedPassword, _ := auth.HashedPassword("oldpassword")

	tests := map[string]struct {
		token       string
		password    string
		prepare     func(f *fields)
		expectedErr apperrors.ApplicationError
	}{
		"successful password reset": {
			token:    "valid-token",
			password: "NewPassword123!",
			prepare: func(f *fields) {
				user := &domain.User{
					ID:                       "user-123",
					Email:                    "test@example.com",
					Password:                 string(hashedPassword),
					PasswordResetToken:       &(&struct{ s string }{"valid-token"}).s,
					PasswordResetTokenExpiry: &validTime,
				}
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{PasswordResetToken: "valid-token"}).Return(user, nil)
				f.repository.EXPECT().Patch(gomock.Any(), "user-123", gomock.Any()).Return("user-123", nil)
				f.repository.EXPECT().InvalidatePasswordResetToken(gomock.Any(), "user-123").Return(nil)
				f.repository.EXPECT().IncrementTokenVersion(gomock.Any(), "user-123").Return(nil)
			},
			expectedErr: nil,
		},
		"invalid token - user not found": {
			token:    "invalid-token",
			password: "NewPassword123!",
			prepare: func(f *fields) {
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{PasswordResetToken: "invalid-token"}).Return(nil, apperrors.NewApplicationError(mappings.UserGetQueryError, nil))
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserResetPasswordInvalidTokenError, nil),
		},
		"invalid token - nil user": {
			token:    "invalid-token2",
			password: "NewPassword123!",
			prepare: func(f *fields) {
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{PasswordResetToken: "invalid-token2"}).Return(nil, nil)
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserResetPasswordInvalidTokenError, nil),
		},
		"expired token": {
			token:    "expired-token",
			password: "NewPassword123!",
			prepare: func(f *fields) {
				user := &domain.User{
					ID:                       "user-123",
					PasswordResetToken:       &(&struct{ s string }{"expired-token"}).s,
					PasswordResetTokenExpiry: &expiredTime,
				}
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{PasswordResetToken: "expired-token"}).Return(user, nil)
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserResetPasswordInvalidTokenError, nil),
		},
		"password update failure": {
			token:    "valid-token-update-fail",
			password: "NewPassword123!",
			prepare: func(f *fields) {
				user := &domain.User{
					ID:                       "user-123",
					Email:                    "test@example.com",
					Password:                 string(hashedPassword),
					PasswordResetToken:       &(&struct{ s string }{"valid-token-update-fail"}).s,
					PasswordResetTokenExpiry: &validTime,
				}
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{PasswordResetToken: "valid-token-update-fail"}).Return(user, nil)
				f.repository.EXPECT().Patch(gomock.Any(), "user-123", gomock.Any()).Return("", apperrors.NewApplicationError(mappings.UserUpdateQueryError, nil))
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserUpdateQueryError, nil),
		},
		"token invalidation failure": {
			token:    "valid-token-invalidate-fail",
			password: "NewPassword123!",
			prepare: func(f *fields) {
				user := &domain.User{
					ID:                       "user-123",
					Email:                    "test@example.com",
					Password:                 string(hashedPassword),
					PasswordResetToken:       &(&struct{ s string }{"valid-token-invalidate-fail"}).s,
					PasswordResetTokenExpiry: &validTime,
				}
				f.repository.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{PasswordResetToken: "valid-token-invalidate-fail"}).Return(user, nil)
				f.repository.EXPECT().Patch(gomock.Any(), "user-123", gomock.Any()).Return("user-123", nil)
				f.repository.EXPECT().InvalidatePasswordResetToken(gomock.Any(), "user-123").Return(apperrors.NewApplicationError(mappings.UserResetPasswordUpdateError, nil))
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserResetPasswordUpdateError, nil),
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

			uc := usecase.NewResetPasswordUsecase(contextFactory)
			_, actualErr := uc.Execute(context.Background(), usecase.ResetPasswordInput{
				Token:    tc.token,
				Password: tc.password,
			})

			assert.Equal(t, tc.expectedErr, actualErr)
		})
	}
}
