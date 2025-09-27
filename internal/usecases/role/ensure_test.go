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

func TestEnsureUsecase_Execute(t *testing.T) {
	type fields struct {
		repository *mock_role.MockRepository
	}

	tests := map[string]struct {
		prepare     func(f *fields)
		expectedErr error
	}{
		"when all roles already exist": {
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{Name: "superadmin"}).
					Return(&domain.Role{
						ID:   "role-123",
						Name: domain.RoleSuperAdmin,
					}, nil)
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{Name: "admin"}).
					Return(&domain.Role{
						ID:   "role-456",
						Name: domain.RoleAdmin,
					}, nil)
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{Name: "user"}).
					Return(&domain.Role{
						ID:   "role-789",
						Name: domain.RoleUser,
					}, nil)
			},
			expectedErr: nil,
		},
		"when some roles need to be created": {
			prepare: func(f *fields) {
				// superadmin exists
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{Name: "superadmin"}).
					Return(&domain.Role{
						ID:   "role-123",
						Name: domain.RoleSuperAdmin,
					}, nil)
				// admin does not exist
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{Name: "admin"}).
					Return(nil, nil)
				f.repository.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(domain.Role{})).
					DoAndReturn(func(ctx context.Context, r domain.Role) (string, error) {
						// Verify that the role name is correct
						if r.Name != domain.RoleAdmin {
							t.Errorf("Expected admin role, got %s", r.Name)
						}
						return "role-456", nil
					})
				// user does not exist
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{Name: "user"}).
					Return(nil, nil)
				f.repository.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(domain.Role{})).
					DoAndReturn(func(ctx context.Context, r domain.Role) (string, error) {
						// Verify that the role name is correct
						if r.Name != domain.RoleUser {
							t.Errorf("Expected user role, got %s", r.Name)
						}
						return "role-789", nil
					})
			},
			expectedErr: nil,
		},
		"when no roles exist and all need to be created": {
			prepare: func(f *fields) {
				// superadmin does not exist
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{Name: "superadmin"}).
					Return(nil, nil)
				f.repository.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(domain.Role{})).
					DoAndReturn(func(ctx context.Context, r domain.Role) (string, error) {
						// Verify that the role name is correct
						if r.Name != domain.RoleSuperAdmin {
							t.Errorf("Expected superadmin role, got %s", r.Name)
						}
						return "role-123", nil
					})
				// admin does not exist
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{Name: "admin"}).
					Return(nil, nil)
				f.repository.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(domain.Role{})).
					DoAndReturn(func(ctx context.Context, r domain.Role) (string, error) {
						// Verify that the role name is correct
						if r.Name != domain.RoleAdmin {
							t.Errorf("Expected admin role, got %s", r.Name)
						}
						return "role-456", nil
					})
				// user does not exist
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{Name: "user"}).
					Return(nil, nil)
				f.repository.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(domain.Role{})).
					DoAndReturn(func(ctx context.Context, r domain.Role) (string, error) {
						// Verify that the role name is correct
						if r.Name != domain.RoleUser {
							t.Errorf("Expected user role, got %s", r.Name)
						}
						return "role-789", nil
					})
			},
			expectedErr: nil,
		},
		"when repository get returns error": {
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{Name: "superadmin"}).
					Return(nil, errors.New("database connection error"))
			},
			expectedErr: errors.New("database connection error"),
		},
		"when repository create returns error": {
			prepare: func(f *fields) {
				// superadmin exists
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{Name: "superadmin"}).
					Return(&domain.Role{
						ID:   "role-123",
						Name: domain.RoleSuperAdmin,
					}, nil)
				// admin does not exist
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{Name: "admin"}).
					Return(nil, nil)
				f.repository.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(domain.Role{})).
					DoAndReturn(func(ctx context.Context, r domain.Role) (string, error) {
						// Verify that the role name is correct
						if r.Name != domain.RoleAdmin {
							t.Errorf("Expected admin role, got %s", r.Name)
						}
						return "", errors.New("database connection error")
					})
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

			uc := usecase.NewEnsureUseCase(contextFactory)
			actualErr := uc.Execute(context.Background())

			assert.Equal(t, tc.expectedErr, actualErr)
		})
	}
}