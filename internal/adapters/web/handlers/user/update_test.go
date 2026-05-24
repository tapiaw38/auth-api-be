package user_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/handlers/user"
	"github.com/tapiaw38/auth-api-be/internal/platform/auth"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	usecase "github.com/tapiaw38/auth-api-be/internal/usecases/user"
	mock_user "github.com/tapiaw38/auth-api-be/internal/usecases/user/mocks"
	"go.uber.org/mock/gomock"
)

func TestUpdateByIDHandler(t *testing.T) {
	type fields struct {
		usecase *mock_user.MockUpdateUsecase
	}

	tests := map[string]struct {
		paramID            string
		body               string
		ctxUserID          string
		ctxRoles           []auth.RoleClaim
		prepare            func(f *fields)
		expectedStatusCode int
		expectedBody       string
	}{
		"when update succeeds for own user": {
			paramID:   "user-123",
			body:      `{"first_name":"John","last_name":"Doe","email":"john@example.com"}`,
			ctxUserID: "john.doe",
			ctxRoles:  []auth.RoleClaim{{ID: "r1", Name: "user"}},
			prepare: func(f *fields) {
				f.usecase.EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(&usecase.UpdateOutput{
						Data: usecase.UserOutputData{
							ID:        "user-123",
							Username:  "john.doe",
							FirstName: "John",
							LastName:  "Doe",
							Email:     "john@example.com",
						},
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `"id":"user-123"`,
		},
		"when user is not authenticated": {
			paramID:   "user-123",
			body:      `{"first_name":"John"}`,
			ctxUserID: "",
			ctxRoles:  nil,
			prepare:   func(f *fields) {},
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `"code":"auth:unauthorized"`,
		},
		"when request body is invalid JSON": {
			paramID:   "user-123",
			body:      `invalid json`,
			ctxUserID: "john.doe",
			ctxRoles:  []auth.RoleClaim{{ID: "r1", Name: "user"}},
			prepare:   func(f *fields) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `"code":"common:request-body-parsing-error"`,
		},
		"when usecase returns unauthorized error": {
			paramID:   "other-user-456",
			body:      `{"first_name":"Hacker"}`,
			ctxUserID: "john.doe",
			ctxRoles:  []auth.RoleClaim{{ID: "r1", Name: "user"}},
			prepare: func(f *fields) {
				f.usecase.EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(nil, apperrors.NewApplicationError(mappings.UserUpdateUnauthorizedError, nil))
			},
			expectedStatusCode: http.StatusForbidden,
			expectedBody:       `"code":"user:update:unauthorized"`,
		},
		"when usecase returns internal error": {
			paramID:   "user-123",
			body:      `{"first_name":"John"}`,
			ctxUserID: "john.doe",
			ctxRoles:  []auth.RoleClaim{{ID: "r1", Name: "user"}},
			prepare: func(f *fields) {
				f.usecase.EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(nil, apperrors.NewApplicationError(mappings.UserUpdateQueryError, nil))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `"code":"user:update:query-error"`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{usecase: mock_user.NewMockUpdateUsecase(ctrl)}
			if tc.prepare != nil {
				tc.prepare(&f)
			}

			gin.SetMode(gin.TestMode)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest(http.MethodPut, "/user/"+tc.paramID, bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")

			ctx := req.Context()
			if tc.ctxUserID != "" {
				ctx = context.WithValue(ctx, "userID", tc.ctxUserID)
			}
			if tc.ctxRoles != nil {
				ctx = context.WithValue(ctx, "userRoles", tc.ctxRoles)
			}
			req = req.WithContext(ctx)

			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: tc.paramID}}

			handler := user.NewUpdateByIDHandler(f.usecase)
			handler(c)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			if tc.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tc.expectedBody)
			}
		})
	}
}
