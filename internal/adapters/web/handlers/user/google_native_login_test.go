package user_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/handlers/user"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	usecase "github.com/tapiaw38/auth-api-be/internal/usecases/user"
	mockUsecase "github.com/tapiaw38/auth-api-be/internal/usecases/user/mocks"
	"go.uber.org/mock/gomock"
)

func TestGoogleNativeLoginHandler(t *testing.T) {
	tests := map[string]struct {
		body               string
		prepare            func(*mockUsecase.MockGoogleNativeLoginUsecase)
		expectedStatusCode int
		expectedCode       string
	}{
		"when native Google login succeeds": {
			body: `{"id_token":"valid-token"}`,
			prepare: func(mock *mockUsecase.MockGoogleNativeLoginUsecase) {
				mock.EXPECT().Execute(gomock.Any(), usecase.GoogleNativeLoginInput{IDToken: "valid-token"}).Return(&usecase.GoogleNativeLoginOutput{
					Token:        "access-token",
					RefreshToken: "refresh-token",
					Data:         usecase.UserOutputData{ID: "user-id"},
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		"when native Google login rejects token": {
			body: `{"id_token":"invalid-token"}`,
			prepare: func(mock *mockUsecase.MockGoogleNativeLoginUsecase) {
				mock.EXPECT().Execute(gomock.Any(), usecase.GoogleNativeLoginInput{IDToken: "invalid-token"}).Return(nil, apperrors.NewApplicationError(mappings.UserLoginGoogleInvalidIDTokenError, errors.New("invalid token")))
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedCode:       mappings.UserLoginGoogleInvalidIDTokenError.InternalCode,
		},
		"when ID token is absent": {
			body:               `{}`,
			prepare:            func(*mockUsecase.MockGoogleNativeLoginUsecase) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedCode:       mappings.RequestBodyParsingError.InternalCode,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mockUsecase.NewMockGoogleNativeLoginUsecase(ctrl)
			tc.prepare(mock)

			gin.SetMode(gin.TestMode)
			request := httptest.NewRequest(http.MethodPost, "/auth/google/native", bytes.NewBufferString(tc.body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(response)
			context.Request = request

			user.NewGoogleNativeLoginHandler(mock)(context)

			assert.Equal(t, tc.expectedStatusCode, response.Code)
			if tc.expectedCode != "" {
				var body struct {
					Code string `json:"code"`
				}
				assert.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
				assert.Equal(t, tc.expectedCode, body.Code)
			}
		})
	}
}
