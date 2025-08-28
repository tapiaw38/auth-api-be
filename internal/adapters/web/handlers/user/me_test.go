package user_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/handlers/user"
	usecase "github.com/tapiaw38/auth-api-be/internal/usecases/user"
	mock_user "github.com/tapiaw38/auth-api-be/internal/usecases/user/mocks"
	"go.uber.org/mock/gomock"
)

func TestMeHandler(t *testing.T) {
	type fields struct {
		usecase *mock_user.MockGetUsecase
	}

	phoneNumber := "+1234567890"
	picture := "https://example.com/avatar.jpg"
	address := "123 Main St, City"

	tests := map[string]struct {
		userID             string
		prepare            func(f *fields)
		expectedStatusCode int
		expectedErr        error
	}{
		"when getting current user successfully": {
			userID: "johndoe",
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), usecase.GetFilterOptions{
					Username: "johndoe",
				}).Return(&usecase.GetOutput{
					Data: usecase.UserOutputData{
						ID:            "user-123",
						FirstName:     "John",
						LastName:      "Doe",
						Email:         "john@example.com",
						PhoneNumber:   &phoneNumber,
						Picture:       &picture,
						Address:       &address,
						IsActive:      true,
						VerifiedEmail: true,
						TokenVersion:  1,
						Roles: []usecase.RoleOutputData{
							{ID: "role-1", Name: "user"},
						},
					},
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		"when getting current user with minimal data": {
			userID: "janesmith",
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), usecase.GetFilterOptions{
					Username: "janesmith",
				}).Return(&usecase.GetOutput{
					Data: usecase.UserOutputData{
						ID:            "user-456",
						FirstName:     "Jane",
						LastName:      "Smith",
						Email:         "jane@example.com",
						PhoneNumber:   nil,
						Picture:       nil,
						Address:       nil,
						IsActive:      true,
						VerifiedEmail: false,
						TokenVersion:  1,
						Roles:         []usecase.RoleOutputData{},
					},
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		"when getting current user with admin role": {
			userID: "adminuser",
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), usecase.GetFilterOptions{
					Username: "adminuser",
				}).Return(&usecase.GetOutput{
					Data: usecase.UserOutputData{
						ID:            "user-789",
						FirstName:     "Admin",
						LastName:      "User",
						Email:         "admin@example.com",
						PhoneNumber:   nil,
						Picture:       nil,
						Address:       nil,
						IsActive:      true,
						VerifiedEmail: true,
						TokenVersion:  2,
						Roles: []usecase.RoleOutputData{
							{ID: "role-1", Name: "admin"},
							{ID: "role-2", Name: "user"},
						},
					},
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		"when usecase returns error": {
			userID: "erroruser",
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), usecase.GetFilterOptions{
					Username: "erroruser",
				}).Return(nil, errors.New("database connection error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        errors.New("database connection error"),
		},
		"when user not found": {
			userID: "nonexistent",
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), usecase.GetFilterOptions{
					Username: "nonexistent",
				}).Return(nil, errors.New("user not found"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        errors.New("user not found"),
		},
		"when user ID is empty": {
			userID: "",
			prepare: func(f *fields) {
				f.usecase.EXPECT().Execute(gomock.Any(), usecase.GetFilterOptions{
					Username: "",
				}).Return(nil, errors.New("user not found"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        errors.New("user not found"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				usecase: mock_user.NewMockGetUsecase(ctrl),
			}

			if tc.prepare != nil {
				tc.prepare(&f)
			}

			gin.SetMode(gin.TestMode)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/me", nil)
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "userID", tc.userID))

			handler := user.NewMeHandler(f.usecase)
			handler(c)

			assert.Equal(t, tc.expectedStatusCode, w.Code)

			if tc.expectedErr != nil {
				// Verificar que el error esté en el body de la respuesta
				assert.Contains(t, w.Body.String(), tc.expectedErr.Error())
			}
		})
	}
}
