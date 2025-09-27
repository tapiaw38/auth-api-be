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

func TestDeleteUsecase_Execute(t *testing.T) {
	type fields struct {
		repository *mock_role.MockRepository
	}

	tests := map[string]struct {
		id          string
		prepare     func(f *fields)
		expectedErr error
	}{
		"when deleting role successfully": {
			id: "role-123",
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-123"}).
					Return(&domain.Role{
						ID:   "role-123",
						Name: domain.RoleAdmin,
					}, nil)
				f.repository.EXPECT().
					Delete(gomock.Any(), "role-123").
					Return(nil)
			},
			expectedErr: nil,
		},
		"when deleting user role successfully": {
			id: "role-456",
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-456"}).
					Return(&domain.Role{
						ID:   "role-456",
						Name: domain.RoleUser,
					}, nil)
				f.repository.EXPECT().
					Delete(gomock.Any(), "role-456").
					Return(nil)
			},
			expectedErr: nil,
		},
		"when deleting superadmin role successfully": {
			id: "role-789",
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-789"}).
					Return(&domain.Role{
						ID:   "role-789",
						Name: domain.RoleSuperAdmin,
					}, nil)
				f.repository.EXPECT().
					Delete(gomock.Any(), "role-789").
					Return(nil)
			},
			expectedErr: nil,
		},
		"when deleting role with empty ID": {
			id:      "",
			prepare: func(f *fields) {},
			expectedErr: errors.New("role ID is required"),
		},
		"when role does not exist": {
			id: "non-existent-role",
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "non-existent-role"}).
					Return(nil, nil)
			},
			expectedErr: errors.New("role not found"),
		},
		"when repository get returns error": {
			id: "role-123",
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-123"}).
					Return(nil, errors.New("database connection error"))
			},
			expectedErr: errors.New("database connection error"),
		},
		"when repository delete returns error": {
			id: "role-123",
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-123"}).
					Return(&domain.Role{
						ID:   "role-123",
						Name: domain.RoleAdmin,
					}, nil)
				f.repository.EXPECT().
					Delete(gomock.Any(), "role-123").
					Return(errors.New("database connection error"))
			},
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

			uc := usecase.NewDeleteUsecase(contextFactory)
			actualErr := uc.Execute(context.Background(), tc.id)

			assert.Equal(t, tc.expectedErr, actualErr)
		})
	}
}