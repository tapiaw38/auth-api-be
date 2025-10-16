package role_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	mock_role "github.com/tapiaw38/auth-api-be/internal/usecases/role/mocks"
	role_handler "github.com/tapiaw38/auth-api-be/internal/adapters/web/handlers/role"
	"go.uber.org/mock/gomock"
)

func TestDeleteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := map[string]struct {
		url          string
		setupUsecase func(*mock_role.MockDeleteUsecase)
		expectedCode int
		expectedBody string
	}{
		"when delete usecase executes successfully": {
			url: "/role/role-123",
			setupUsecase: func(mockUsecase *mock_role.MockDeleteUsecase) {
				mockUsecase.EXPECT().Execute(gomock.Any(), "role-123").Return(nil)
			},
			expectedCode: 200,
			expectedBody: `{"message":"role deleted successfully"}`,
		},
		"when delete usecase for different role ID": {
			url: "/role/role-456",
			setupUsecase: func(mockUsecase *mock_role.MockDeleteUsecase) {
				mockUsecase.EXPECT().Execute(gomock.Any(), "role-456").Return(nil)
			},
			expectedCode: 200,
			expectedBody: `{"message":"role deleted successfully"}`,
		},
		"when delete usecase returns error": {
			url: "/role/role-123",
			setupUsecase: func(mockUsecase *mock_role.MockDeleteUsecase) {
				mockUsecase.EXPECT().Execute(gomock.Any(), "role-123").Return(assert.AnError)
			},
			expectedCode: 500,
			expectedBody: `{"message":"` + assert.AnError.Error() + `"}`,
		},
		"when role ID is missing": {
			url: "/role/",
			setupUsecase: func(mockUsecase *mock_role.MockDeleteUsecase) {},
			expectedCode: 400,
			expectedBody: `{"message":"role ID is required"}`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mock_role.NewMockDeleteUsecase(ctrl)
			if tc.setupUsecase != nil {
				tc.setupUsecase(mockUsecase)
			}

			handler := role_handler.NewDeleteHandler(mockUsecase)

			req := httptest.NewRequest(http.MethodDelete, tc.url, nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Set the URL parameter
			if tc.url == "/role/role-123" {
				c.Params = []gin.Param{
					{Key: "id", Value: "role-123"},
				}
			} else if tc.url == "/role/role-456" {
				c.Params = []gin.Param{
					{Key: "id", Value: "role-456"},
				}
			} else if tc.url == "/role/" {
				c.Params = []gin.Param{
					{Key: "id", Value: ""},
				}
			}

			handler(c)

			assert.Equal(t, tc.expectedCode, w.Code)
			if tc.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tc.expectedBody)
			}
		})
	}
}