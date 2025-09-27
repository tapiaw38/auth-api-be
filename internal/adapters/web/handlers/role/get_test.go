package role_test

import (
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

func TestGetHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := map[string]struct {
		url          string
		setupUsecase func(*mock_role.MockGetUsecase)
		expectedCode int
		expectedBody string
	}{
		"when get usecase executes successfully by ID": {
			url: "/role/role-123",
			setupUsecase: func(mockUsecase *mock_role.MockGetUsecase) {
				mockUsecase.EXPECT().Execute(gomock.Any(), roleUsecase.GetFilterOptions{ID: "role-123"}).Return(&roleUsecase.GetOutput{
					Data: roleUsecase.RoleOutputData{
						ID:   "role-123",
						Name: "admin",
					},
				}, nil)
			},
			expectedCode: 200,
			expectedBody: `{"data":{"id":"role-123","name":"admin"}}`,
		},
		"when get usecase executes successfully by name query": {
			url: "/role?name=user",
			setupUsecase: func(mockUsecase *mock_role.MockGetUsecase) {
				mockUsecase.EXPECT().Execute(gomock.Any(), roleUsecase.GetFilterOptions{Name: "user"}).Return(&roleUsecase.GetOutput{
					Data: roleUsecase.RoleOutputData{
						ID:   "role-456",
						Name: "user",
					},
				}, nil)
			},
			expectedCode: 200,
			expectedBody: `{"data":{"id":"role-456","name":"user"}}`,
		},
		"when get usecase returns role not found": {
			url: "/role/role-123",
			setupUsecase: func(mockUsecase *mock_role.MockGetUsecase) {
				mockUsecase.EXPECT().Execute(gomock.Any(), roleUsecase.GetFilterOptions{ID: "role-123"}).Return(nil, nil)
			},
			expectedCode: 404,
			expectedBody: `{"message":"role not found"}`,
		},
		"when get usecase returns error": {
			url: "/role/role-123",
			setupUsecase: func(mockUsecase *mock_role.MockGetUsecase) {
				mockUsecase.EXPECT().Execute(gomock.Any(), roleUsecase.GetFilterOptions{ID: "role-123"}).Return(nil, assert.AnError)
			},
			expectedCode: 500,
			expectedBody: `{"message":"` + assert.AnError.Error() + `"}`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mock_role.NewMockGetUsecase(ctrl)
			if tc.setupUsecase != nil {
				tc.setupUsecase(mockUsecase)
			}

			handler := role_handler.NewGetHandler(mockUsecase)

			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Set the URL parameter if it exists in the test URL
			if tc.url == "/role/role-123" {
				c.Params = []gin.Param{
					{Key: "id", Value: "role-123"},
				}
			}

			handler(c)

			assert.Equal(t, tc.expectedCode, w.Code)
			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, w.Body.String())
			}
		})
	}
}