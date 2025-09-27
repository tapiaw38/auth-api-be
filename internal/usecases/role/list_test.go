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

func TestListUsecase_Execute(t *testing.T) {
	type fields struct {
		repository *mock_role.MockRepository
	}

	tests := map[string]struct {
		input       usecase.ListFilterOptions
		prepare     func(f *fields)
		expected    []usecase.RoleOutputData
		expectedErr error
	}{
		"when listing roles successfully": {
			input: usecase.ListFilterOptions{},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					List(gomock.Any(), roleRepo.ListFilterOptions{}).
					Return([]domain.Role{
						{
							ID:   "role-123",
							Name: domain.RoleAdmin,
						},
						{
							ID:   "role-456",
							Name: domain.RoleUser,
						},
					}, nil)
			},
			expected: []usecase.RoleOutputData{
				{
					ID:   "role-123",
					Name: "admin",
				},
				{
					ID:   "role-456",
					Name: "user",
				},
			},
		},
		"when listing roles with name filter successfully": {
			input: usecase.ListFilterOptions{
				Name: "admin",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					List(gomock.Any(), roleRepo.ListFilterOptions{Name: "admin"}).
					Return([]domain.Role{
						{
							ID:   "role-123",
							Name: domain.RoleAdmin,
						},
					}, nil)
			},
			expected: []usecase.RoleOutputData{
				{
					ID:   "role-123",
					Name: "admin",
				},
			},
		},
		"when listing empty results": {
			input: usecase.ListFilterOptions{},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					List(gomock.Any(), roleRepo.ListFilterOptions{}).
					Return([]domain.Role{}, nil)
			},
			expected: []usecase.RoleOutputData{},
		},
		"when listing with nil results": {
			input: usecase.ListFilterOptions{},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					List(gomock.Any(), roleRepo.ListFilterOptions{}).
					Return(nil, nil)
			},
			expected: nil,
		},
		"when repository returns error": {
			input: usecase.ListFilterOptions{},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					List(gomock.Any(), roleRepo.ListFilterOptions{}).
					Return(nil, errors.New("database connection error"))
			},
			expected:    nil,
			expectedErr: errors.New("database connection error"),
		},
		"when listing superadmin roles": {
			input: usecase.ListFilterOptions{
				Name: "superadmin",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					List(gomock.Any(), roleRepo.ListFilterOptions{Name: "superadmin"}).
					Return([]domain.Role{
						{
							ID:   "role-999",
							Name: domain.RoleSuperAdmin,
						},
					}, nil)
			},
			expected: []usecase.RoleOutputData{
				{
					ID:   "role-999",
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

			uc := usecase.NewListUsecase(contextFactory)
			actual, actualErr := uc.Execute(context.Background(), tc.input)

			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedErr, actualErr)
		})
	}
}