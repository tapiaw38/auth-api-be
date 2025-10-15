package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/handlers/user"
	mock_usecase "github.com/tapiaw38/auth-api-be/internal/usecases/user/mocks"
	"go.uber.org/mock/gomock"
)

func TestSetPasswordHandler(t *testing.T) {
	type fields struct {
		usecase *mock_usecase.MockSetPasswordUsecase
	}

	tests := map[string]struct {
		body               any
		prepare            func(f *fields)
		expectedStatusCode int
		expectedErr        error
	}{
		"when password set successfully": {
			body: user.SetPasswordRequest{
				NewPassword: "NewPassword123!",
			},
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedStatusCode: http.StatusNoContent,
		},
		"when usecase returns an error": {
			body: user.SetPasswordRequest{
				NewPassword: "NewPassword123!",
			},
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(errors.New("some error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        errors.New("some error"),
		},
		"when user is not SSO user": {
			body: user.SetPasswordRequest{
				NewPassword: "NewPassword123!",
			},
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(errors.New("only SSO users can set initial password"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        errors.New("only SSO users can set initial password"),
		},
		"when password is too weak": {
			body: user.SetPasswordRequest{
				NewPassword: "weak",
			},
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(errors.New("password must be at least 8 characters long"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        errors.New("password must be at least 8 characters long"),
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
				usecase: mock_usecase.NewMockSetPasswordUsecase(ctrl),
			}

			if tc.prepare != nil {
				tc.prepare(&f)
			}

			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			bodyBytes, _ := json.Marshal(tc.body)
			c.Request = httptest.NewRequest(http.MethodPost, "/user/me/set-password", bytes.NewReader(bodyBytes))
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "userID", "user-123"))

			handler := user.NewSetPasswordHandler(f.usecase)
			handler(c)

			assert.Equal(t, tc.expectedStatusCode, w.Code)

			if tc.expectedErr != nil {
				assert.Contains(t, w.Body.String(), tc.expectedErr.Error())
			}
		})
	}
}
