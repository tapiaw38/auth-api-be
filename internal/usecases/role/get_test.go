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

func TestGetUsecase_Execute(t *testing.T) {
	type fields struct {
		repository *mock_role.MockRepository
	}

	tests := map[string]struct {
		input       usecase.GetFilterOptions
		prepare     func(f *fields)
		expected    *usecase.GetOutput
		expectedErr error
	}{
		"when getting role by ID successfully": {
			input: usecase.GetFilterOptions{
				ID: "role-123",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-123"}).
					Return(&domain.Role{
						ID:   "role-123",
						Name: domain.RoleAdmin,
					}, nil)
			},
			expected: &usecase.GetOutput{
				Data: usecase.RoleOutputData{
					ID:   "role-123",
					Name: "admin",
				},
			},
		},
		"when getting role by name successfully": {
			input: usecase.GetFilterOptions{
				Name: "user",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{Name: "user"}).
					Return(&domain.Role{
						ID:   "role-456",
						Name: domain.RoleUser,
					}, nil)
			},
			expected: &usecase.GetOutput{
				Data: usecase.RoleOutputData{
					ID:   "role-456",
					Name: "user",
				},
			},
		},
		"when repository returns error": {
			input: usecase.GetFilterOptions{
				ID: "role-123",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-123"}).
					Return(nil, errors.New("database connection error"))
			},
			expectedErr: errors.New("database connection error"),
		},
		"when role not found": {
			input: usecase.GetFilterOptions{
				ID: "non-existent-role",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "non-existent-role"}).
					Return(nil, nil)
			},
			expectedErr: errors.New("role not found"),
		},
		"when getting superadmin role": {
			input: usecase.GetFilterOptions{
				ID: "role-789",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-789"}).
					Return(&domain.Role{
						ID:   "role-789",
						Name: domain.RoleSuperAdmin,
					}, nil)
			},
			expected: &usecase.GetOutput{
				Data: usecase.RoleOutputData{
					ID:   "role-789",
					Name: "superadmin",
				},
			},
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

			uc := usecase.NewGetUsecase(contextFactory)
			actual, actualErr := uc.Execute(context.Background(), tc.input)

			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedErr, actualErr)
		})
	}
}