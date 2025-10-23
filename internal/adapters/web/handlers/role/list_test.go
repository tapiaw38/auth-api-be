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
	roleUsecase "github.com/tapiaw38/auth-api-be/internal/usecases/role"
	mock_role "github.com/tapiaw38/auth-api-be/internal/usecases/role/mocks"
	"go.uber.org/mock/gomock"
)

func TestListHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := map[string]struct {
		queryParams  string
		setupUsecase func(*mock_role.MockListUsecase)
		expectedCode int
		expectedBody string
	}{
		"when list usecase executes successfully with no filters": {
			queryParams: "",
			setupUsecase: func(mockUsecase *mock_role.MockListUsecase) {
				mockUsecase.EXPECT().Execute(gomock.Any(), roleUsecase.ListFilterOptions{}).Return([]roleUsecase.RoleOutputData{
					{
						ID:   "role-123",
						Name: "admin",
					},
					{
						ID:   "role-456",
						Name: "user",
					},
				}, nil)
			},
			expectedCode: 200,
			expectedBody: `[{"id":"role-123","name":"admin"},{"id":"role-456","name":"user"}]`,
		},
		"when list usecase executes successfully with name filter": {
			queryParams: "?name=admin",
			setupUsecase: func(mockUsecase *mock_role.MockListUsecase) {
				mockUsecase.EXPECT().Execute(gomock.Any(), roleUsecase.ListFilterOptions{Name: "admin"}).Return([]roleUsecase.RoleOutputData{
					{
						ID:   "role-123",
						Name: "admin",
					},
				}, nil)
			},
			expectedCode: 200,
			expectedBody: `[{"id":"role-123","name":"admin"}]`,
		},
		"when list usecase returns empty results": {
			queryParams: "",
			setupUsecase: func(mockUsecase *mock_role.MockListUsecase) {
				mockUsecase.EXPECT().Execute(gomock.Any(), roleUsecase.ListFilterOptions{}).Return([]roleUsecase.RoleOutputData{}, nil)
			},
			expectedCode: 200,
			expectedBody: `[]`,
		},
		"when list usecase returns error": {
			queryParams: "",
			setupUsecase: func(mockUsecase *mock_role.MockListUsecase) {
				mockUsecase.EXPECT().Execute(gomock.Any(), roleUsecase.ListFilterOptions{}).Return(nil, apperrors.NewApplicationError(mappings.RoleListQueryError, errors.New("database error")))
			},
			expectedCode: 500,
			expectedBody: `{"code":"role:list:query-error","message":"failed to retrieve role list"}`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mock_role.NewMockListUsecase(ctrl)
			if tc.setupUsecase != nil {
				tc.setupUsecase(mockUsecase)
			}

			handler := role_handler.NewListHandler(mockUsecase)

			url := "/list"
			if tc.queryParams != "" {
				url += tc.queryParams
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
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