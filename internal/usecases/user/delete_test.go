package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories"
	mock_user "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user/mocks"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	usecase "github.com/tapiaw38/auth-api-be/internal/usecases/user"
	"go.uber.org/mock/gomock"
)

func TestDeleteUsecase(t *testing.T) {
	type fields struct {
		repository *mock_user.MockRepository
	}

	tests := map[string]struct {
		userID       string
		prepare      func(f *fields)
		expectedID   string
		expectedErr  error
	}{
		"successful delete": {
			userID: "user-123",
			prepare: func(f *fields) {
				f.repository.EXPECT().Delete(gomock.Any(), "user-123").Return(nil)
			},
			expectedID: "user-123",
		},
		"error - user not found": {
			userID: "non-existent-user",
			prepare: func(f *fields) {
				f.repository.EXPECT().Delete(gomock.Any(), "non-existent-user").Return(errors.New("user not found"))
			},
			expectedErr: errors.New("user not found"),
		},
		"error - database error": {
			userID: "user-456",
			prepare: func(f *fields) {
				f.repository.EXPECT().Delete(gomock.Any(), "user-456").Return(errors.New("database error"))
			},
			expectedErr: errors.New("database error"),
		},
		"error - user has dependencies": {
			userID: "user-789",
			prepare: func(f *fields) {
				f.repository.EXPECT().Delete(gomock.Any(), "user-789").Return(errors.New("cannot delete user with existing dependencies"))
			},
			expectedErr: errors.New("cannot delete user with existing dependencies"),
		},
		"successful delete - empty id should still call repository": {
			userID: "",
			prepare: func(f *fields) {
				f.repository.EXPECT().Delete(gomock.Any(), "").Return(errors.New("invalid user id"))
			},
			expectedErr: errors.New("invalid user id"),
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
			resultID, actualErr := uc.Execute(context.Background(), tc.userID)

			if tc.expectedErr != nil {
				assert.Error(t, actualErr)
				assert.Equal(t, tc.expectedErr.Error(), actualErr.Error())
				assert.Empty(t, resultID)
			} else {
				assert.NoError(t, actualErr)
				assert.Equal(t, tc.expectedID, resultID)
			}
		})
	}
}
