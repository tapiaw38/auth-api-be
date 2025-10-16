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
	handler "github.com/tapiaw38/auth-api-be/internal/adapters/web/handlers/user"
	usecase "github.com/tapiaw38/auth-api-be/internal/usecases/user"
	mock_usecase "github.com/tapiaw38/auth-api-be/internal/usecases/user/mocks"
	"go.uber.org/mock/gomock"
)

func TestRegisterHandler(t *testing.T) {
	type fields struct {
		usecase *mock_usecase.MockRegisterUsecase
	}

	tests := map[string]struct {
		body               any
		prepare            func(f *fields)
		expectedStatusCode int
		expectedResponse   *usecase.RegisterOutput
		expectedErr        error
	}{
		"when registration is successful": {
			body: map[string]string{
				"first_name": "John",
				"last_name":  "Doe",
				"username":   "johndoe",
				"email":      "john.doe@example.com",
				"password":   "SecurePassword123!",
			},
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(&usecase.RegisterOutput{
					Data: usecase.UserOutputData{
						ID:        "user-123",
						FirstName: "John",
						LastName:  "Doe",
						Email:     "john.doe@example.com",
					},
					Token: "",
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: &usecase.RegisterOutput{
				Data: usecase.UserOutputData{
					ID:        "user-123",
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@example.com",
				},
			},
		},
		"when usecase returns an error - email already exists": {
			body: map[string]string{
				"first_name": "Jane",
				"last_name":  "Smith",
				"username":   "janesmith",
				"email":      "existing@example.com",
				"password":   "SecurePassword123!",
			},
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, errors.New("email already exists"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        errors.New("email already exists"),
		},
		"when usecase returns an error - username already exists": {
			body: map[string]string{
				"first_name": "Bob",
				"last_name":  "Wilson",
				"username":   "existinguser",
				"email":      "bob@example.com",
				"password":   "SecurePassword123!",
			},
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, errors.New("username already exists"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        errors.New("username already exists"),
		},
		"when usecase returns an error - weak password": {
			body: map[string]string{
				"first_name": "Alice",
				"last_name":  "Brown",
				"username":   "alicebrown",
				"email":      "alice@example.com",
				"password":   "weak",
			},
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, errors.New("password must be at least 8 characters long"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        errors.New("password must be at least 8 characters long"),
		},
		"when usecase returns an error - missing uppercase": {
			body: map[string]string{
				"first_name": "Charlie",
				"last_name":  "Davis",
				"username":   "charliedavis",
				"email":      "charlie@example.com",
				"password":   "password123!",
			},
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, errors.New("password must contain at least one uppercase letter"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        errors.New("password must contain at least one uppercase letter"),
		},
		"when usecase returns an error - missing special character": {
			body: map[string]string{
				"first_name": "David",
				"last_name":  "Evans",
				"username":   "davidevans",
				"email":      "david@example.com",
				"password":   "Password123",
			},
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, errors.New("password must contain at least one special character (!@#$%^&*()_+-=[]{}|;:,.<>?)"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        errors.New("password must contain at least one special character"),
		},
		"when request body is invalid": {
			body:               "invalid body",
			prepare:            nil,
			expectedStatusCode: http.StatusBadRequest,
		},
		"when request has missing required fields": {
			body: map[string]string{
				"first_name": "Emma",
				// Missing other required fields
			},
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, errors.New("validation error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        errors.New("validation error"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				usecase: mock_usecase.NewMockRegisterUsecase(ctrl),
			}

			if tc.prepare != nil {
				tc.prepare(&f)
			}

			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			bodyBytes, _ := json.Marshal(tc.body)
			c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(bodyBytes))

			handlerFunc := handler.NewRegisterHandler(f.usecase)
			handlerFunc(c)

			assert.Equal(t, tc.expectedStatusCode, w.Code)

			if tc.expectedResponse != nil {
				var response usecase.RegisterOutput
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse.Data.ID, response.Data.ID)
				assert.Equal(t, tc.expectedResponse.Data.Email, response.Data.Email)
				assert.Equal(t, tc.expectedResponse.Data.FirstName, response.Data.FirstName)
				assert.Equal(t, tc.expectedResponse.Data.LastName, response.Data.LastName)
			}

			if tc.expectedErr != nil {
				assert.Contains(t, w.Body.String(), tc.expectedErr.Error())
			}
		})
	}
}
