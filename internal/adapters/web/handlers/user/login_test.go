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
	usecase "github.com/tapiaw38/auth-api-be/internal/usecases/user"
	mock_usecase "github.com/tapiaw38/auth-api-be/internal/usecases/user/mocks"
	"go.uber.org/mock/gomock"
)

func TestLoginHandler(t *testing.T) {
	type fields struct {
		usecase *mock_usecase.MockLoginUsecase
	}

	tests := map[string]struct {
		body               any
		prepare            func(f *fields)
		expectedStatusCode int
		expectedResponse   *usecase.LoginOutput
		expectedErr        error
	}{
		"when login is successful with email and password": {
			body: usecase.LoginInput{
				Email:    "test@example.com",
				Password: "Password123!",
			},
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(&usecase.LoginOutput{
					Data: usecase.UserOutputData{
						ID:        "user-123",
						FirstName: "Test",
						LastName:  "User",
						Email:     "test@example.com",
					},
					Token: "fake-jwt-token",
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: &usecase.LoginOutput{
				Data: usecase.UserOutputData{
					ID:        "user-123",
					FirstName: "Test",
					LastName:  "User",
					Email:     "test@example.com",
				},
				Token: "fake-jwt-token",
			},
		},
		"when login is successful with Google SSO": {
			body: usecase.LoginInput{
				SsoType: "google",
				Code:    "google-auth-code",
			},
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(&usecase.LoginOutput{
					Data: usecase.UserOutputData{
						ID:        "user-456",
						FirstName: "Google",
						LastName:  "User",
						Email:     "google@example.com",
					},
					Token: "fake-jwt-token-sso",
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: &usecase.LoginOutput{
				Data: usecase.UserOutputData{
					ID:        "user-456",
					FirstName: "Google",
					LastName:  "User",
					Email:     "google@example.com",
				},
				Token: "fake-jwt-token-sso",
			},
		},
		"when usecase returns an error - invalid credentials": {
			body: usecase.LoginInput{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, errors.New("invalid credentials"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        errors.New("invalid credentials"),
		},
		"when usecase returns an error - user not found": {
			body: usecase.LoginInput{
				Email:    "nonexistent@example.com",
				Password: "Password123!",
			},
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, errors.New("user not found"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        errors.New("user not found"),
		},
		"when usecase returns an error - user is not active": {
			body: usecase.LoginInput{
				Email:    "inactive@example.com",
				Password: "Password123!",
			},
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, errors.New("user is not active"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        errors.New("user is not active"),
		},
		"when usecase returns an error - SSO authentication required": {
			body: usecase.LoginInput{
				Email:    "ssouser@example.com",
				Password: "Password123!",
			},
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, errors.New("this account uses SSO authentication. Please use Google login"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        errors.New("this account uses SSO authentication. Please use Google login"),
		},
		"when request body is invalid": {
			body:               "invalid body",
			prepare:            nil,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				usecase: mock_usecase.NewMockLoginUsecase(ctrl),
			}

			if tc.prepare != nil {
				tc.prepare(&f)
			}

			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			bodyBytes, _ := json.Marshal(tc.body)
			c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(bodyBytes))

			handler := user.NewLoginHandler(f.usecase)
			handler(c)

			assert.Equal(t, tc.expectedStatusCode, w.Code)

			if tc.expectedResponse != nil {
				var response usecase.LoginOutput
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse.Data.ID, response.Data.ID)
				assert.Equal(t, tc.expectedResponse.Data.Email, response.Data.Email)
				assert.Equal(t, tc.expectedResponse.Token, response.Token)
			}

			if tc.expectedErr != nil {
				assert.Contains(t, w.Body.String(), tc.expectedErr.Error())
			}
		})
	}
}
