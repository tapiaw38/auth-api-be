package role_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	role_handler "github.com/tapiaw38/auth-api-be/internal/adapters/web/handlers/role"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	mock_role "github.com/tapiaw38/auth-api-be/internal/usecases/role/mocks"
	"go.uber.org/mock/gomock"
)

func TestEnsureHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := map[string]struct {
		setupUsecase func(*mock_role.MockEnsureUseCase)
		expectedCode int
		expectedBody string
	}{
		"when ensure usecase executes successfully": {
			setupUsecase: func(mockUsecase *mock_role.MockEnsureUseCase) {
				mockUsecase.EXPECT().Execute(gomock.Any()).Return(nil)
			},
			expectedCode: 200,
			expectedBody: `{"message":"ok"}`,
		},
		"when ensure usecase returns error": {
			setupUsecase: func(mockUsecase *mock_role.MockEnsureUseCase) {
				mockUsecase.EXPECT().Execute(gomock.Any()).Return(apperrors.NewApplicationError(mappings.RoleEnsureDefaultRoleNotFoundError, errors.New("default role not found")))
			},
			expectedCode: 500,
			expectedBody: `{"code":"role:ensure:default-role-not-found","message":"default user role not found in the system"}`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mock_role.NewMockEnsureUseCase(ctrl)
			if tc.setupUsecase != nil {
				tc.setupUsecase(mockUsecase)
			}

			handler := role_handler.NewEnsureHandler(mockUsecase)

			req := httptest.NewRequest(http.MethodGet, "/ensure", nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			handler(c)

			assert.Equal(t, tc.expectedCode, w.Code)
			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, w.Body.String())
			}
		})
	}
}
