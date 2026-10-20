package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/integrations"
	mockSSO "github.com/tapiaw38/auth-api-be/internal/adapters/web/integrations/sso/mocks"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	usecase "github.com/tapiaw38/auth-api-be/internal/usecases/user"
	"go.uber.org/mock/gomock"
)

func TestGoogleNativeLoginUsecaseExecute(t *testing.T) {
	tests := map[string]struct {
		input       usecase.GoogleNativeLoginInput
		prepare     func(*mockSSO.MockIntegration)
		expectedErr string
	}{
		"when ID token is missing": {
			input:       usecase.GoogleNativeLoginInput{},
			prepare:     func(*mockSSO.MockIntegration) {},
			expectedErr: mappings.UserLoginGoogleInvalidIDTokenError.InternalCode,
		},
		"when ID token verification fails": {
			input: usecase.GoogleNativeLoginInput{IDToken: "invalid-token"},
			prepare: func(sso *mockSSO.MockIntegration) {
				sso.EXPECT().VerifyIDToken(gomock.Any(), "invalid-token").Return(nil, errors.New("invalid token"))
			},
			expectedErr: mappings.UserLoginGoogleInvalidIDTokenError.InternalCode,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			sso := mockSSO.NewMockIntegration(ctrl)
			tc.prepare(sso)

			contextFactory := func(opts ...appcontext.Option) *appcontext.Context {
				return &appcontext.Context{
					Integrations: &integrations.Integrations{SSO: sso},
				}
			}

			output, appErr := usecase.NewGoogleNativeLoginUsecase(contextFactory).Execute(context.Background(), tc.input)

			assert.Nil(t, output)
			assert.NotNil(t, appErr)
			assert.Equal(t, tc.expectedErr, appErr.InternalCode())
		})
	}
}
