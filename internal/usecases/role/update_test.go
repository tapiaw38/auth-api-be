package role_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories"
	roleRepo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role"
	mock_role "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/role/mocks"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	usecase "github.com/tapiaw38/auth-api-be/internal/usecases/role"
	"go.uber.org/mock/gomock"
)

func TestUpdateUsecase_Execute(t *testing.T) {
	type fields struct {
		repository *mock_role.MockRepository
	}

	tests := map[string]struct {
		id          string
		input       usecase.UpdateInput
		prepare     func(f *fields)
		expected    *usecase.UpdateOutput
		expectedErr error
	}{
		"when updating role successfully": {
			id: "role-123",
			input: usecase.UpdateInput{
				Name: "user",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-123"}).
					Return(&domain.Role{
						ID:   "role-123",
						Name: domain.RoleAdmin,
					}, nil)
				f.repository.EXPECT().
					Update(gomock.Any(), "role-123", &domain.Role{
						ID:   "role-123",
						Name: domain.RoleUser,
					}).
					Return("role-123", nil)
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-123"}).
					Return(&domain.Role{
						ID:   "role-123",
						Name: domain.RoleUser,
					}, nil)
			},
			expected: &usecase.UpdateOutput{
				Data: usecase.RoleOutputData{
					ID:   "role-123",
					Name: "user",
				},
			},
		},
		"when updating role to admin": {
			id: "role-456",
			input: usecase.UpdateInput{
				Name: "admin",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-456"}).
					Return(&domain.Role{
						ID:   "role-456",
						Name: domain.RoleUser,
					}, nil)
				f.repository.EXPECT().
					Update(gomock.Any(), "role-456", &domain.Role{
						ID:   "role-456",
						Name: domain.RoleAdmin,
					}).
					Return("role-456", nil)
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-456"}).
					Return(&domain.Role{
						ID:   "role-456",
						Name: domain.RoleAdmin,
					}, nil)
			},
			expected: &usecase.UpdateOutput{
				Data: usecase.RoleOutputData{
					ID:   "role-456",
					Name: "admin",
				},
			},
		},
		"when updating role to superadmin": {
			id: "role-789",
			input: usecase.UpdateInput{
				Name: "superadmin",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-789"}).
					Return(&domain.Role{
						ID:   "role-789",
						Name: domain.RoleUser,
					}, nil)
				f.repository.EXPECT().
					Update(gomock.Any(), "role-789", &domain.Role{
						ID:   "role-789",
						Name: domain.RoleSuperAdmin,
					}).
					Return("role-789", nil)
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-789"}).
					Return(&domain.Role{
						ID:   "role-789",
						Name: domain.RoleSuperAdmin,
					}, nil)
			},
			expected: &usecase.UpdateOutput{
				Data: usecase.RoleOutputData{
					ID:   "role-789",
					Name: "superadmin",
				},
			},
		},
		"when updating role with empty ID": {
			id: "",
			input: usecase.UpdateInput{
				Name: "user",
			},
			prepare:     func(f *fields) {},
			expected:    nil,
			expectedErr: apperrors.NewApplicationError(mappings.RoleUpdateIDRequiredError, nil),
		},
		"when updating role with empty name": {
			id: "role-123",
			input: usecase.UpdateInput{
				Name: "",
			},
			prepare:     func(f *fields) {},
			expected:    nil,
			expectedErr: apperrors.NewApplicationError(mappings.RoleUpdateNameRequiredError, nil),
		},
		"when updating role with invalid name": {
			id: "role-123",
			input: usecase.UpdateInput{
				Name: "invalid-role",
			},
			prepare:     func(f *fields) {},
			expected:    nil,
			expectedErr: apperrors.NewApplicationError(mappings.RoleUpdateInvalidNameError, nil),
		},
		"when role does not exist": {
			id: "non-existent-role",
			input: usecase.UpdateInput{
				Name: "user",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "non-existent-role"}).
					Return(nil, nil)
			},
			expected:    nil,
			expectedErr: apperrors.NewApplicationError(mappings.RoleUpdateNotFoundError, nil),
		},
		"when repository get returns error": {
			id: "role-123",
			input: usecase.UpdateInput{
				Name: "user",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-123"}).
					Return(nil, apperrors.NewApplicationError(mappings.RoleGetQueryError, nil))
			},
			expected:    nil,
			expectedErr: apperrors.NewApplicationError(mappings.RoleGetQueryError, nil),
		},
		"when repository update returns error": {
			id: "role-123",
			input: usecase.UpdateInput{
				Name: "user",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-123"}).
					Return(&domain.Role{
						ID:   "role-123",
						Name: domain.RoleAdmin,
					}, nil)
				f.repository.EXPECT().
					Update(gomock.Any(), "role-123", &domain.Role{
						ID:   "role-123",
						Name: domain.RoleUser,
					}).
					Return("", apperrors.NewApplicationError(mappings.RoleUpdateQueryError, nil))
			},
			expected:    nil,
			expectedErr: apperrors.NewApplicationError(mappings.RoleUpdateQueryError, nil),
		},
		"when repository get returns error after update": {
			id: "role-123",
			input: usecase.UpdateInput{
				Name: "user",
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-123"}).
					Return(&domain.Role{
						ID:   "role-123",
						Name: domain.RoleAdmin,
					}, nil)
				f.repository.EXPECT().
					Update(gomock.Any(), "role-123", &domain.Role{
						ID:   "role-123",
						Name: domain.RoleUser,
					}).
					Return("role-123", nil)
				f.repository.EXPECT().
					Get(gomock.Any(), roleRepo.GetFilterOptions{ID: "role-123"}).
					Return(nil, apperrors.NewApplicationError(mappings.RoleGetQueryError, nil))
			},
			expected:    nil,
			expectedErr: apperrors.NewApplicationError(mappings.RoleGetQueryError, nil),
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

			uc := usecase.NewUpdateUsecase(contextFactory)
			actual, actualErr := uc.Execute(context.Background(), tc.id, tc.input)

			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedErr, actualErr)
		})
	}
}