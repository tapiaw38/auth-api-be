package role_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	role_handler "github.com/tapiaw38/auth-api-be/internal/adapters/web/handlers/role"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	roleUsecase "github.com/tapiaw38/auth-api-be/internal/usecases/role"
	mock_role "github.com/tapiaw38/auth-api-be/internal/usecases/role/mocks"
	"go.uber.org/mock/gomock"
)

func TestUpdateHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := map[string]struct {
		url          string
		body         string
		setupUsecase func(*mock_role.MockUpdateUsecase)
		expectedCode int
		expectedBody string
	}{
		"when update usecase executes successfully": {
			url:  "/role/role-123",
			body: `{"name":"user"}`,
			setupUsecase: func(mockUsecase *mock_role.MockUpdateUsecase) {
				mockUsecase.EXPECT().Execute(gomock.Any(), "role-123", roleUsecase.UpdateInput{Name: "user"}).Return(&roleUsecase.UpdateOutput{
					Data: roleUsecase.RoleOutputData{
						ID:   "role-123",
						Name: "user",
					},
				}, nil)
			},
			expectedCode: 200,
			expectedBody: `{"data":{"id":"role-123","name":"user"}}`,
		},
		"when update usecase updates to admin role": {
			url:  "/role/role-456",
			body: `{"name":"admin"}`,
			setupUsecase: func(mockUsecase *mock_role.MockUpdateUsecase) {
				mockUsecase.EXPECT().Execute(gomock.Any(), "role-456", roleUsecase.UpdateInput{Name: "admin"}).Return(&roleUsecase.UpdateOutput{
					Data: roleUsecase.RoleOutputData{
						ID:   "role-456",
						Name: "admin",
					},
				}, nil)
			},
			expectedCode: 200,
			expectedBody: `{"data":{"id":"role-456","name":"admin"}}`,
		},
		"when update usecase returns error": {
			url:  "/role/role-123",
			body: `{"name":"user"}`,
			setupUsecase: func(mockUsecase *mock_role.MockUpdateUsecase) {
				mockUsecase.EXPECT().Execute(gomock.Any(), "role-123", roleUsecase.UpdateInput{Name: "user"}).Return(nil, apperrors.NewApplicationError(mappings.RoleUpdateQueryError, errors.New("database error")))
			},
			expectedCode: 500,
			expectedBody: `{"code":"role:update:query-error","message":"failed to update role"}`,
		},
		"when role ID is missing": {
			url:  "/role/",
			body: `{"name":"user"}`,
			setupUsecase: func(mockUsecase *mock_role.MockUpdateUsecase) {
				mockUsecase.EXPECT().Execute(gomock.Any(), "", roleUsecase.UpdateInput{Name: "user"}).
					Return(nil, apperrors.NewApplicationError(mappings.RoleUpdateIDRequiredError, nil))
			},
			expectedCode: 400,
			expectedBody: `{"code":"role:update:id-required","message":"role ID is required"}`,
		},
		"when request body is invalid": {
			url:  "/role/role-123",
			body: `{"invalid":}`,
			setupUsecase: func(mockUsecase *mock_role.MockUpdateUsecase) {},
			expectedCode: 400,
			expectedBody: `{"code":"common:request-body-parsing-error","message":"invalid request body format"}`,
		},
		"when name is missing": {
			url:  "/role/role-123",
			body: `{}`,
			setupUsecase: func(mockUsecase *mock_role.MockUpdateUsecase) {},
			expectedCode: 400,
			expectedBody: `{"code":"common:request-body-parsing-error","message":"invalid request body format"}`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mock_role.NewMockUpdateUsecase(ctrl)
			if tc.setupUsecase != nil {
				tc.setupUsecase(mockUsecase)
			}

			handler := role_handler.NewUpdateHandler(mockUsecase)

			req := httptest.NewRequest(http.MethodPut, tc.url, bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
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
				if tc.expectedCode == 400 {
					// For 400 errors, we want to check that the response contains the expected message
					if tc.url == "/role/" {
						// This case is for when role ID is missing
						assert.Contains(t, w.Body.String(), "role ID is required")
					} else {
						// Other 400 cases should contain the new error format
						assert.Contains(t, w.Body.String(), "common:request-body-parsing-error")
						assert.Contains(t, w.Body.String(), "invalid request body format")
					}
				} else {
					assert.Contains(t, w.Body.String(), tc.expectedBody)
				}
			}
		})
	}
}