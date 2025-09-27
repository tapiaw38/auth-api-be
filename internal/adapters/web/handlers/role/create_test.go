package role_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	mock_role "github.com/tapiaw38/auth-api-be/internal/usecases/role/mocks"
	roleUsecase "github.com/tapiaw38/auth-api-be/internal/usecases/role"
	role_handler "github.com/tapiaw38/auth-api-be/internal/adapters/web/handlers/role"
	"go.uber.org/mock/gomock"
)

func TestCreateHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := map[string]struct {
		body         string
		setupUsecase func(*mock_role.MockCreateUsecase)
		expectedCode int
		expectedBody string
	}{
		"when create usecase executes successfully": {
			body: `{"name":"admin"}`,
			setupUsecase: func(mockUsecase *mock_role.MockCreateUsecase) {
				mockUsecase.EXPECT().Execute(gomock.Any(), roleUsecase.CreateInput{Name: "admin"}).Return(&roleUsecase.CreateOutput{
					Data: roleUsecase.RoleOutputData{
						ID:   "role-123",
						Name: "admin",
					},
				}, nil)
			},
			expectedCode: 201,
			expectedBody: `{"data":{"id":"role-123","name":"admin"}}`,
		},
		"when create usecase with user role": {
			body: `{"name":"user"}`,
			setupUsecase: func(mockUsecase *mock_role.MockCreateUsecase) {
				mockUsecase.EXPECT().Execute(gomock.Any(), roleUsecase.CreateInput{Name: "user"}).Return(&roleUsecase.CreateOutput{
					Data: roleUsecase.RoleOutputData{
						ID:   "role-456",
						Name: "user",
					},
				}, nil)
			},
			expectedCode: 201,
			expectedBody: `{"data":{"id":"role-456","name":"user"}}`,
		},
		"when create usecase returns error": {
			body: `{"name":"admin"}`,
			setupUsecase: func(mockUsecase *mock_role.MockCreateUsecase) {
				mockUsecase.EXPECT().Execute(gomock.Any(), roleUsecase.CreateInput{Name: "admin"}).Return(nil, assert.AnError)
			},
			expectedCode: 500,
			expectedBody: `{"message":"` + assert.AnError.Error() + `"}`,
		},
		"when request body is invalid": {
			body:         `{"invalid":}`,
			setupUsecase: func(mockUsecase *mock_role.MockCreateUsecase) {},
			expectedCode: 400,
			expectedBody: `{"message":"Invalid request body"}`,
		},
		"when name is missing": {
			body:         `{}`,
			setupUsecase: func(mockUsecase *mock_role.MockCreateUsecase) {},
			expectedCode: 400,
			expectedBody: `{"message":"Invalid request body"}`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mock_role.NewMockCreateUsecase(ctrl)
			if tc.setupUsecase != nil {
				tc.setupUsecase(mockUsecase)
			}

			handler := role_handler.NewCreateHandler(mockUsecase)

			req := httptest.NewRequest(http.MethodPost, "/role", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			handler(c)

			assert.Equal(t, tc.expectedCode, w.Code)
			if tc.expectedBody != "" {
				if tc.expectedCode == 400 {
					// For 400 errors, we want to check that the response contains the error message
					assert.Contains(t, w.Body.String(), "Invalid request body")
				} else {
					assert.Contains(t, w.Body.String(), tc.expectedBody)
				}
			}
		})
	}
}