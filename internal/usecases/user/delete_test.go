package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories"
	mock_user "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user/mocks"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	usecase "github.com/tapiaw38/auth-api-be/internal/usecases/user"
	"go.uber.org/mock/gomock"
)

func TestDeleteUsecase(t *testing.T) {
	type fields struct {
		repository *mock_user.MockRepository
	}

	tests := map[string]struct {
		userID      string
		prepare     func(f *fields)
		expectedErr apperrors.ApplicationError
	}{
		"successful delete": {
			userID: "user-123",
			prepare: func(f *fields) {
				f.repository.EXPECT().Delete(gomock.Any(), "user-123").Return(nil)
			},
		},
		"error - user not found": {
			userID: "non-existent-user",
			prepare: func(f *fields) {
				f.repository.EXPECT().Delete(gomock.Any(), "non-existent-user").Return(apperrors.NewApplicationError(mappings.UserDeleteNotFoundError, nil))
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserDeleteNotFoundError, nil),
		},
		"error - database error": {
			userID: "user-456",
			prepare: func(f *fields) {
				f.repository.EXPECT().Delete(gomock.Any(), "user-456").Return(apperrors.NewApplicationError(mappings.UserDeleteQueryError, nil))
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserDeleteQueryError, nil),
		},
		"error - user has dependencies": {
			userID: "user-789",
			prepare: func(f *fields) {
				f.repository.EXPECT().Delete(gomock.Any(), "user-789").Return(apperrors.NewApplicationError(mappings.UserDeleteQueryError, nil))
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserDeleteQueryError, nil),
		},
		"error - empty id validation": {
			userID: "",
			prepare: func(f *fields) {
				// No repository call expected - validation happens before
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserDeleteNotFoundError, nil),
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

			uc := usecase.NewDeleteUsecase(contextFactory)
			actualErr := uc.Execute(context.Background(), tc.userID)

			assert.Equal(t, tc.expectedErr, actualErr)
		})
	}
}
