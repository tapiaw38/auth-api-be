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
	usecase "github.com/tapiaw38/auth-api-be/internal/usecases/user"
	"go.uber.org/mock/gomock"
)

func TestListUsecase(t *testing.T) {
	type fields struct {
		repository *mock_user.MockRepository
	}

	trueValue := true
	falseValue := false

	tests := map[string]struct {
		filters      usecase.ListFilterOptions
		prepare      func(f *fields)
		expectedLen  int
		expectedErr  error
		validateData func(t *testing.T, result []usecase.UserOutputData)
	}{
		"successful list - all users": {
			filters: usecase.ListFilterOptions{
				Limit:  10,
				Offset: 0,
			},
			prepare: func(f *fields) {
				users := []*domain.User{
					{
						ID:        "user-1",
						FirstName: "John",
						LastName:  "Doe",
						Email:     "john@example.com",
						IsActive:  true,
					},
					{
						ID:        "user-2",
						FirstName: "Jane",
						LastName:  "Smith",
						Email:     "jane@example.com",
						IsActive:  true,
					},
				}
				f.repository.EXPECT().List(gomock.Any(), gomock.Any()).Return(users, nil)
			},
			expectedLen: 2,
			validateData: func(t *testing.T, result []usecase.UserOutputData) {
				assert.Equal(t, "user-1", result[0].ID)
				assert.Equal(t, "John", result[0].FirstName)
				assert.Equal(t, "john@example.com", result[0].Email)
				assert.Equal(t, "user-2", result[1].ID)
				assert.Equal(t, "Jane", result[1].FirstName)
			},
		},
		"successful list - filter by active users": {
			filters: usecase.ListFilterOptions{
				IsActive: &trueValue,
				Limit:    10,
				Offset:   0,
			},
			prepare: func(f *fields) {
				users := []*domain.User{
					{
						ID:        "user-1",
						FirstName: "Active",
						LastName:  "User",
						Email:     "active@example.com",
						IsActive:  true,
					},
				}
				f.repository.EXPECT().List(gomock.Any(), user_repo.ListFilterOptions{
					IsActive: &trueValue,
					Limit:    10,
					Offset:   0,
				}).Return(users, nil)
			},
			expectedLen: 1,
			validateData: func(t *testing.T, result []usecase.UserOutputData) {
				assert.Equal(t, "user-1", result[0].ID)
				assert.True(t, result[0].IsActive)
			},
		},
		"successful list - filter by inactive users": {
			filters: usecase.ListFilterOptions{
				IsActive: &falseValue,
				Limit:    10,
				Offset:   0,
			},
			prepare: func(f *fields) {
				users := []*domain.User{
					{
						ID:        "user-3",
						FirstName: "Inactive",
						LastName:  "User",
						Email:     "inactive@example.com",
						IsActive:  false,
					},
				}
				f.repository.EXPECT().List(gomock.Any(), user_repo.ListFilterOptions{
					IsActive: &falseValue,
					Limit:    10,
					Offset:   0,
				}).Return(users, nil)
			},
			expectedLen: 1,
			validateData: func(t *testing.T, result []usecase.UserOutputData) {
				assert.Equal(t, "user-3", result[0].ID)
				assert.False(t, result[0].IsActive)
			},
		},
		"successful list - filter by verified email": {
			filters: usecase.ListFilterOptions{
				VerifiedEmail: &trueValue,
				Limit:         10,
				Offset:        0,
			},
			prepare: func(f *fields) {
				users := []*domain.User{
					{
						ID:            "user-4",
						FirstName:     "Verified",
						LastName:      "User",
						Email:         "verified@example.com",
						VerifiedEmail: true,
					},
				}
				f.repository.EXPECT().List(gomock.Any(), user_repo.ListFilterOptions{
					VerifiedEmail: &trueValue,
					Limit:         10,
					Offset:        0,
				}).Return(users, nil)
			},
			expectedLen: 1,
			validateData: func(t *testing.T, result []usecase.UserOutputData) {
				assert.Equal(t, "user-4", result[0].ID)
				assert.True(t, result[0].VerifiedEmail)
			},
		},
		"successful list - filter by role": {
			filters: usecase.ListFilterOptions{
				RoleID: "role-123",
				Limit:  10,
				Offset: 0,
			},
			prepare: func(f *fields) {
				users := []*domain.User{
					{
						ID:        "user-5",
						FirstName: "Role",
						LastName:  "User",
						Email:     "role@example.com",
						Roles: []domain.Role{
							{ID: "role-123", Name: "admin"},
						},
					},
				}
				f.repository.EXPECT().List(gomock.Any(), user_repo.ListFilterOptions{
					RoleID: "role-123",
					Limit:  10,
					Offset: 0,
				}).Return(users, nil)
			},
			expectedLen: 1,
			validateData: func(t *testing.T, result []usecase.UserOutputData) {
				assert.Equal(t, "user-5", result[0].ID)
				assert.Len(t, result[0].Roles, 1)
				assert.Equal(t, "role-123", result[0].Roles[0].ID)
			},
		},
		"successful list - empty result": {
			filters: usecase.ListFilterOptions{
				Limit:  10,
				Offset: 0,
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*domain.User{}, nil)
			},
			expectedLen: 0,
		},
		"successful list - with pagination": {
			filters: usecase.ListFilterOptions{
				Limit:  5,
				Offset: 10,
			},
			prepare: func(f *fields) {
				users := []*domain.User{
					{
						ID:        "user-11",
						FirstName: "Page2",
						LastName:  "User1",
						Email:     "page2user1@example.com",
					},
					{
						ID:        "user-12",
						FirstName: "Page2",
						LastName:  "User2",
						Email:     "page2user2@example.com",
					},
				}
				f.repository.EXPECT().List(gomock.Any(), user_repo.ListFilterOptions{
					Limit:  5,
					Offset: 10,
				}).Return(users, nil)
			},
			expectedLen: 2,
		},
		"error - repository returns error": {
			filters: usecase.ListFilterOptions{
				Limit:  10,
				Offset: 0,
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))
			},
			expectedErr: errors.New("database error"),
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

			uc := usecase.NewListUsecase(contextFactory)
			result, actualErr := uc.Execute(context.Background(), tc.filters)

			if tc.expectedErr != nil {
				assert.Error(t, actualErr)
				assert.Equal(t, tc.expectedErr.Error(), actualErr.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, actualErr)
				assert.Len(t, result, tc.expectedLen)
				if tc.validateData != nil {
					tc.validateData(t, result)
				}
			}
		})
	}
}
