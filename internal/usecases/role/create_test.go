package role_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories"
	roleRepo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role"
	mock_role "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role/mocks"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	usecase "github.com/tapiaw38/auth-api-be/internal/usecases/role"
	"go.uber.org/mock/gomock"
)

func TestCreateUsecase_Execute(t *testing.T) {
	type fields struct {
		repository *mock_role.MockRepository
	}

	tests := map[string]struct {
		input       usecase.CreateInput
		prepare     func(f *fields)
		expected    *usecase.CreateOutput
		expectedErr error
	}{
		"when creating role successfully": {
			input: usecase.CreateInput{
				Name: "admin",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Create(gomock.Any(), domain.Role{Name: domain.RoleAdmin}).
					Return("role-123", nil)
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-123"}).
					Return(&domain.Role{
						ID:   "role-123",
						Name: domain.RoleAdmin,
					}, nil)
			},
			expected: &usecase.CreateOutput{
				Data: usecase.RoleOutputData{
					ID:   "role-123",
					Name: "admin",
				},
			},
		},
		"when creating user role successfully": {
			input: usecase.CreateInput{
				Name: "user",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Create(gomock.Any(), domain.Role{Name: domain.RoleUser}).
					Return("role-456", nil)
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-456"}).
					Return(&domain.Role{
						ID:   "role-456",
						Name: domain.RoleUser,
					}, nil)
			},
			expected: &usecase.CreateOutput{
				Data: usecase.RoleOutputData{
					ID:   "role-456",
					Name: "user",
				},
			},
		},
		"when creating superadmin role successfully": {
			input: usecase.CreateInput{
				Name: "superadmin",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Create(gomock.Any(), domain.Role{Name: domain.RoleSuperAdmin}).
					Return("role-789", nil)
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-789"}).
					Return(&domain.Role{
						ID:   "role-789",
						Name: domain.RoleSuperAdmin,
					}, nil)
			},
			expected: &usecase.CreateOutput{
				Data: usecase.RoleOutputData{
					ID:   "role-789",
					Name: "superadmin",
				},
			},
		},
		"when creating role with empty name": {
			input: usecase.CreateInput{
				Name: "",
			},
			prepare:     func(f *fields) {},
			expected:    nil,
			expectedErr: errors.New("role name is required"),
		},
		"when creating role with invalid name": {
			input: usecase.CreateInput{
				Name: "invalid-role",
			},
			prepare:     func(f *fields) {},
			expected:    nil,
			expectedErr: errors.New("invalid role name"),
		},
		"when repository create returns error": {
			input: usecase.CreateInput{
				Name: "admin",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Create(gomock.Any(), domain.Role{Name: domain.RoleAdmin}).
					Return("", errors.New("database connection error"))
			},
			expected:    nil,
			expectedErr: errors.New("database connection error"),
		},
		"when repository get returns error after create": {
			input: usecase.CreateInput{
				Name: "admin",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Create(gomock.Any(), domain.Role{Name: domain.RoleAdmin}).
					Return("role-123", nil)
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-123"}).
					Return(nil, errors.New("database connection error"))
			},
			expected:    nil,
			expectedErr: errors.New("database connection error"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repository: mock_role.NewMockRepository(ctrl),
			}

			if tc.prepare != nil {
				tc.prepare(&f)
			}

			contextFactory := func(opts ...appcontext.Option) *appcontext.Context {
				return &appcontext.Context{
					Repositories: &repositories.Repositories{
						Role: f.repository,
					},
				}
			}

			uc := usecase.NewCreateUsecase(contextFactory)
			actual, actualErr := uc.Execute(context.Background(), tc.input)

			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedErr, actualErr)
		})
	}
}
