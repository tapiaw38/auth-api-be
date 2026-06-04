package user_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories"
	user_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	mock_user "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user/mocks"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	usecase "github.com/tapiaw38/auth-api-be/internal/usecases/user"
	"go.uber.org/mock/gomock"
)

func TestUpdateUsecase(t *testing.T) {
	type fields struct {
		repository *mock_user.MockRepository
	}

	firstName := " John "
	verifiedEmail := true
	isActive := false
	invalidEmail := "not-an-email"
	emptyAddress := (*string)(nil)

	tests := map[string]struct {
		input          usecase.UpdateInput
		prepare        func(f *fields)
		expectedOutput *usecase.UpdateOutput
		expectedErr    apperrors.ApplicationError
	}{
		"partial self update preserves omitted fields": {
			input: usecase.UpdateInput{
				ID:             "user-123",
				AuthUsername:   "john.doe",
				CanManageUsers: false,
				Patch: domain.UserPatch{
					FirstName: domain.Optional[string]{
						Set:   true,
						Value: &firstName,
					},
				},
			},
			prepare: func(f *fields) {
				existingUser := &domain.User{
					ID:            "user-123",
					Username:      "john.doe",
					FirstName:     "Jane",
					LastName:      "Doe",
					Email:         "jane@example.com",
					IsActive:      true,
					VerifiedEmail: true,
				}

				f.repository.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{Username: "john.doe"}).
					Return(existingUser, nil)
				f.repository.EXPECT().
					Patch(gomock.Any(), "user-123", gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, user *domain.User) (string, apperrors.ApplicationError) {
						assert.Equal(t, "John", user.FirstName)
						assert.Equal(t, "Doe", user.LastName)
						assert.Equal(t, "jane@example.com", user.Email)
						assert.True(t, user.IsActive)
						assert.True(t, user.VerifiedEmail)
						return "user-123", nil
					})
				f.repository.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{ID: "user-123"}).
					Return(&domain.User{
						ID:            "user-123",
						Username:      "john.doe",
						FirstName:     "John",
						LastName:      "Doe",
						Email:         "jane@example.com",
						IsActive:      true,
						VerifiedEmail: true,
					}, nil)
			},
			expectedOutput: &usecase.UpdateOutput{
				Data: usecase.UserOutputData{
					ID:            "user-123",
					Username:      "john.doe",
					FirstName:     "John",
					LastName:      "Doe",
					Email:         "jane@example.com",
					IsActive:      true,
					VerifiedEmail: true,
				},
			},
		},
		"self update can clear nullable fields with null": {
			input: usecase.UpdateInput{
				ID:             "user-123",
				AuthUsername:   "john.doe",
				CanManageUsers: false,
				Patch: domain.UserPatch{
					Address: domain.Optional[string]{
						Set:   true,
						Value: emptyAddress,
					},
				},
			},
			prepare: func(f *fields) {
				address := "123 Main St"
				f.repository.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{Username: "john.doe"}).
					Return(&domain.User{
						ID:       "user-123",
						Username: "john.doe",
						Address:  &address,
					}, nil)
				f.repository.EXPECT().
					Patch(gomock.Any(), "user-123", gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, user *domain.User) (string, apperrors.ApplicationError) {
						assert.Nil(t, user.Address)
						return "user-123", nil
					})
				f.repository.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{ID: "user-123"}).
					Return(&domain.User{
						ID:       "user-123",
						Username: "john.doe",
						Address:  nil,
					}, nil)
			},
			expectedOutput: &usecase.UpdateOutput{
				Data: usecase.UserOutputData{
					ID:       "user-123",
					Username: "john.doe",
					Address:  nil,
				},
			},
		},
		"non admin cannot update restricted fields": {
			input: usecase.UpdateInput{
				ID:             "user-123",
				AuthUsername:   "john.doe",
				CanManageUsers: false,
				Patch: domain.UserPatch{
					VerifiedEmail: domain.Optional[bool]{
						Set:   true,
						Value: &verifiedEmail,
					},
				},
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{Username: "john.doe"}).
					Return(&domain.User{
						ID:       "user-123",
						Username: "john.doe",
					}, nil)
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserUpdateRestrictedFieldsError,
				errors.New("user attempted to update restricted fields")),
		},
		"non admin cannot update another user": {
			input: usecase.UpdateInput{
				ID:             "user-456",
				AuthUsername:   "john.doe",
				CanManageUsers: false,
				Patch: domain.UserPatch{
					FirstName: domain.Optional[string]{
						Set:   true,
						Value: &firstName,
					},
				},
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{Username: "john.doe"}).
					Return(&domain.User{
						ID:       "user-123",
						Username: "john.doe",
					}, nil)
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserUpdateUnauthorizedError,
				errors.New("user attempted to update another user without admin role")),
		},
		"admin can update another user including sensitive fields": {
			input: usecase.UpdateInput{
				ID:             "user-456",
				AuthUsername:   "admin.user",
				CanManageUsers: true,
				Patch: domain.UserPatch{
					FirstName: domain.Optional[string]{
						Set:   true,
						Value: &firstName,
					},
					VerifiedEmail: domain.Optional[bool]{
						Set:   true,
						Value: &verifiedEmail,
					},
					IsActive: domain.Optional[bool]{
						Set:   true,
						Value: &isActive,
					},
				},
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{Username: "admin.user"}).
					Return(&domain.User{
						ID:       "admin-1",
						Username: "admin.user",
					}, nil)
				f.repository.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{ID: "user-456"}).
					Return(&domain.User{
						ID:            "user-456",
						Username:      "other.user",
						FirstName:     "Jane",
						LastName:      "Doe",
						Email:         "jane@example.com",
						IsActive:      true,
						VerifiedEmail: false,
					}, nil)
				f.repository.EXPECT().
					Patch(gomock.Any(), "user-456", gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, user *domain.User) (string, apperrors.ApplicationError) {
						assert.Equal(t, "John", user.FirstName)
						assert.False(t, user.IsActive)
						assert.True(t, user.VerifiedEmail)
						assert.Equal(t, "jane@example.com", user.Email)
						return "user-456", nil
					})
				f.repository.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{ID: "user-456"}).
					Return(&domain.User{
						ID:            "user-456",
						Username:      "other.user",
						FirstName:     "John",
						LastName:      "Doe",
						Email:         "jane@example.com",
						IsActive:      false,
						VerifiedEmail: true,
					}, nil)
			},
			expectedOutput: &usecase.UpdateOutput{
				Data: usecase.UserOutputData{
					ID:            "user-456",
					Username:      "other.user",
					FirstName:     "John",
					LastName:      "Doe",
					Email:         "jane@example.com",
					IsActive:      false,
					VerifiedEmail: true,
				},
			},
		},
		"target user not found": {
			input: usecase.UpdateInput{
				ID:             "user-456",
				AuthUsername:   "admin.user",
				CanManageUsers: true,
				Patch: domain.UserPatch{
					FirstName: domain.Optional[string]{
						Set:   true,
						Value: &firstName,
					},
				},
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{Username: "admin.user"}).
					Return(&domain.User{
						ID:       "admin-1",
						Username: "admin.user",
					}, nil)
				f.repository.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{ID: "user-456"}).
					Return(nil, nil)
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserUpdateNotFoundError, nil),
		},
		"empty payload is rejected": {
			input: usecase.UpdateInput{
				ID:             "user-123",
				AuthUsername:   "john.doe",
				CanManageUsers: false,
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{Username: "john.doe"}).
					Return(&domain.User{
						ID:       "user-123",
						Username: "john.doe",
					}, nil)
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserUpdateInvalidInputError,
				errors.New("at least one field is required")),
		},
		"invalid email is rejected": {
			input: usecase.UpdateInput{
				ID:             "user-123",
				AuthUsername:   "john.doe",
				CanManageUsers: false,
				Patch: domain.UserPatch{
					Email: domain.Optional[string]{
						Set:   true,
						Value: &invalidEmail,
					},
				},
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{Username: "john.doe"}).
					Return(&domain.User{
						ID:       "user-123",
						Username: "john.doe",
					}, nil)
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserUpdateInvalidInputError,
				errors.New("invalid email format")),
		},
		"empty string nullable field is rejected": {
			input: usecase.UpdateInput{
				ID:             "user-123",
				AuthUsername:   "john.doe",
				CanManageUsers: false,
				Patch: domain.UserPatch{
					Picture: domain.Optional[string]{
						Set:   true,
						Value: strPtr("   "),
					},
				},
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{Username: "john.doe"}).
					Return(&domain.User{
						ID:       "user-123",
						Username: "john.doe",
					}, nil)
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserUpdateInvalidInputError,
				errors.New("picture cannot be empty; use null to clear it")),
		},
		"null bool is rejected": {
			input: usecase.UpdateInput{
				ID:             "user-123",
				AuthUsername:   "john.doe",
				CanManageUsers: true,
				Patch: domain.UserPatch{
					IsActive: domain.Optional[bool]{
						Set:   true,
						Value: nil,
					},
				},
			},
			prepare: func(f *fields) {
				f.repository.EXPECT().
					Get(gomock.Any(), user_repo.GetFilterOptions{Username: "john.doe"}).
					Return(&domain.User{
						ID:       "user-123",
						Username: "john.doe",
					}, nil)
			},
			expectedErr: apperrors.NewApplicationError(mappings.UserUpdateInvalidInputError,
				errors.New("is_active cannot be null")),
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

			uc := usecase.NewUpdateUsecase(contextFactory)
			actualOutput, actualErr := uc.Execute(context.Background(), tc.input)

			assert.Equal(t, tc.expectedOutput, actualOutput)
			assert.Equal(t, tc.expectedErr, actualErr)
			if actualErr != nil && tc.expectedErr != nil {
				assert.True(t, strings.Contains(actualErr.OriginalMessage(), tc.expectedErr.OriginalMessage()) || actualErr.OriginalMessage() == tc.expectedErr.OriginalMessage())
			}
		})
	}
}

func strPtr(value string) *string {
	return &value
}
