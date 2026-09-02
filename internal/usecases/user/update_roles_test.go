package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories"
	role_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role"
	mock_role "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role/mocks"
	user_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	mock_user "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user/mocks"
	mock_user_role "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user/role/mocks"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	usecase "github.com/tapiaw38/auth-api-be/internal/usecases/user"
	"go.uber.org/mock/gomock"
)

func TestUpdateRolesUsecase(t *testing.T) {
	type fields struct {
		users     *mock_user.MockRepository
		roles     *mock_role.MockRepository
		userRoles *mock_user_role.MockRepository
	}

	actor := &domain.User{ID: "actor-id", Username: "actor-username"}
	adminRole := &domain.Role{ID: "role-admin", Name: domain.RoleAdmin}
	userRole := &domain.Role{ID: "role-user", Name: domain.RoleUser}

	tests := map[string]struct {
		input       usecase.UpdateRolesInput
		prepare     func(f *fields)
		expectedErr string
	}{
		"promotes a student to teacher": {
			input: usecase.UpdateRolesInput{
				TargetID:      "target-id",
				ActorUsername: "actor-username",
				Roles:         []string{"admin"},
			},
			prepare: func(f *fields) {
				f.users.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{Username: "actor-username"}).
					Return(actor, nil)
				f.users.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{ID: "target-id"}).
					Return(&domain.User{
						ID:    "target-id",
						Roles: []domain.Role{{ID: "role-user", Name: domain.RoleUser}},
					}, nil)
				f.roles.EXPECT().
					Get(gomock.Any(), role_repo.GetFilterOptions{Name: "admin"}).
					Return(adminRole, nil)
				f.userRoles.EXPECT().
					Create(gomock.Any(), domain.UserRole{UserID: "target-id", RoleID: "role-admin"}).
					Return(&domain.UserRole{}, nil)
				f.roles.EXPECT().
					Get(gomock.Any(), role_repo.GetFilterOptions{Name: "user"}).
					Return(userRole, nil)
				f.userRoles.EXPECT().
					Delete(gomock.Any(), "target-id", "role-user").
					Return(&domain.UserRole{}, nil)
				f.users.EXPECT().
					IncrementTokenVersion(gomock.Any(), "target-id").
					Return(nil)
				f.users.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{ID: "target-id"}).
					Return(&domain.User{ID: "target-id"}, nil)
			},
		},
		"leaves untouched what is already correct": {
			input: usecase.UpdateRolesInput{
				TargetID:      "target-id",
				ActorUsername: "actor-username",
				Roles:         []string{"admin"},
			},
			prepare: func(f *fields) {
				f.users.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{Username: "actor-username"}).
					Return(actor, nil)
				f.users.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{ID: "target-id"}).
					Return(&domain.User{
						ID:    "target-id",
						Roles: []domain.Role{{ID: "role-admin", Name: domain.RoleAdmin}},
					}, nil)
				f.users.EXPECT().
					IncrementTokenVersion(gomock.Any(), "target-id").
					Return(nil)
				f.users.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{ID: "target-id"}).
					Return(&domain.User{ID: "target-id"}, nil)
			},
		},
		"keeps a superadmin role it cannot manage": {
			input: usecase.UpdateRolesInput{
				TargetID:      "target-id",
				ActorUsername: "actor-username",
				Roles:         []string{"admin"},
			},
			prepare: func(f *fields) {
				f.users.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{Username: "actor-username"}).
					Return(actor, nil)
				f.users.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{ID: "target-id"}).
					Return(&domain.User{
						ID: "target-id",
						Roles: []domain.Role{
							{ID: "role-super", Name: domain.RoleSuperAdmin},
							{ID: "role-admin", Name: domain.RoleAdmin},
						},
					}, nil)
				f.users.EXPECT().
					IncrementTokenVersion(gomock.Any(), "target-id").
					Return(nil)
				f.users.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{ID: "target-id"}).
					Return(&domain.User{ID: "target-id"}, nil)
			},
		},
		"refuses to change your own roles": {
			input: usecase.UpdateRolesInput{
				TargetID:      "actor-id",
				ActorUsername: "actor-username",
				Roles:         []string{"admin"},
			},
			prepare: func(f *fields) {
				f.users.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{Username: "actor-username"}).
					Return(actor, nil)
			},
			expectedErr: mappings.UserRolesSelfUpdateError.InternalCode,
		},
		"refuses to grant superadmin": {
			input: usecase.UpdateRolesInput{
				TargetID:      "target-id",
				ActorUsername: "actor-username",
				Roles:         []string{"superadmin"},
			},
			prepare: func(f *fields) {
				f.users.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{Username: "actor-username"}).
					Return(actor, nil)
			},
			expectedErr: mappings.UserRolesNotAssignableError.InternalCode,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				users:     mock_user.NewMockRepository(ctrl),
				roles:     mock_role.NewMockRepository(ctrl),
				userRoles: mock_user_role.NewMockRepository(ctrl),
			}
			tc.prepare(&f)

			contextFactory := func(opts ...appcontext.Option) *appcontext.Context {
				return &appcontext.Context{
					Repositories: &repositories.Repositories{
						User:     f.users,
						Role:     f.roles,
						UserRole: f.userRoles,
					},
				}
			}

			uc := usecase.NewUpdateRolesUsecase(contextFactory)
			output, appErr := uc.Execute(context.Background(), tc.input)

			if tc.expectedErr != "" {
				assert.Nil(t, output)
				assert.NotNil(t, appErr)
				assert.Equal(t, tc.expectedErr, appErr.InternalCode())
				return
			}

			assert.Nil(t, appErr)
			assert.NotNil(t, output)
		})
	}
}
