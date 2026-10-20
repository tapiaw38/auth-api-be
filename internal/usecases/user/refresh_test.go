package user_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories"
	mock_refresh "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/refresh_token/mocks"
	user_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	mock_user "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user/mocks"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	"github.com/tapiaw38/auth-api-be/internal/platform/auth"
	"github.com/tapiaw38/auth-api-be/internal/platform/config"
	usecase "github.com/tapiaw38/auth-api-be/internal/usecases/user"
	"go.uber.org/mock/gomock"
)

func TestRefreshUsecase(t *testing.T) {
	config.InitConfigService(&config.ConfigurationService{
		ServerConfig: config.ServerConfig{JWTSecret: "test-secret"},
	})

	type fields struct {
		users         *mock_user.MockRepository
		refreshTokens *mock_refresh.MockRepository
	}

	const presented = "presented-token"
	hash := auth.HashRefreshToken(presented)
	activeUser := &domain.User{ID: "user-123", Username: "nymia", IsActive: true}
	revokedAt := time.Now().Add(-time.Minute)

	tests := map[string]struct {
		token          string
		prepare        func(f *fields)
		expectedStatus int
	}{
		"a live token is exchanged and rotated": {
			token: presented,
			prepare: func(f *fields) {
				f.refreshTokens.EXPECT().GetByHash(gomock.Any(), hash).Return(&domain.RefreshToken{
					ID: "rt-1", UserID: "user-123", TokenHash: hash, ExpiresAt: time.Now().Add(time.Hour),
				}, nil)
				f.users.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{ID: "user-123"}).Return(activeUser, nil)
				f.refreshTokens.EXPECT().Revoke(gomock.Any(), "rt-1").Return(nil)
				f.refreshTokens.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		"reusing a spent token drops every session": {
			token: presented,
			prepare: func(f *fields) {
				f.refreshTokens.EXPECT().GetByHash(gomock.Any(), hash).Return(&domain.RefreshToken{
					ID: "rt-1", UserID: "user-123", TokenHash: hash, ExpiresAt: time.Now().Add(time.Hour), RevokedAt: &revokedAt,
				}, nil)
				f.refreshTokens.EXPECT().RevokeAllForUser(gomock.Any(), "user-123").Return(nil)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		"an expired token is refused": {
			token: presented,
			prepare: func(f *fields) {
				f.refreshTokens.EXPECT().GetByHash(gomock.Any(), hash).Return(&domain.RefreshToken{
					ID: "rt-1", UserID: "user-123", TokenHash: hash, ExpiresAt: time.Now().Add(-time.Hour),
				}, nil)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		"an unknown token is refused": {
			token: presented,
			prepare: func(f *fields) {
				f.refreshTokens.EXPECT().GetByHash(gomock.Any(), hash).Return(nil, nil)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		"a deactivated account cannot refresh": {
			token: presented,
			prepare: func(f *fields) {
				f.refreshTokens.EXPECT().GetByHash(gomock.Any(), hash).Return(&domain.RefreshToken{
					ID: "rt-1", UserID: "user-123", TokenHash: hash, ExpiresAt: time.Now().Add(time.Hour),
				}, nil)
				f.users.EXPECT().Get(gomock.Any(), user_repo.GetFilterOptions{ID: "user-123"}).Return(&domain.User{ID: "user-123", IsActive: false}, nil)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		"an empty token is refused before touching the database": {
			token:          "",
			prepare:        func(f *fields) {},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				users:         mock_user.NewMockRepository(ctrl),
				refreshTokens: mock_refresh.NewMockRepository(ctrl),
			}
			tc.prepare(&f)

			contextFactory := func(opts ...appcontext.Option) *appcontext.Context {
				return &appcontext.Context{
					Repositories: &repositories.Repositories{
						User:         f.users,
						RefreshToken: f.refreshTokens,
					},
				}
			}

			output, appErr := usecase.NewRefreshUsecase(contextFactory).Execute(
				context.Background(), usecase.RefreshInput{RefreshToken: tc.token},
			)

			if tc.expectedStatus != 0 {
				assert.NotNil(t, appErr)
				assert.Equal(t, tc.expectedStatus, appErr.StatusCode())
				return
			}

			assert.Nil(t, appErr)
			assert.NotEmpty(t, output.Token)
			assert.NotEmpty(t, output.RefreshToken)
			assert.NotEqual(t, presented, output.RefreshToken)
		})
	}
}
