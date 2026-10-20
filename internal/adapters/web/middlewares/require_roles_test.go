package middlewares_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/middlewares"
	"github.com/tapiaw38/auth-api-be/internal/domain"
	"github.com/tapiaw38/auth-api-be/internal/platform/auth"
)

func TestRequireRoles(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := map[string]struct {
		roles        []auth.RoleClaim
		expectedCode int
	}{
		"superadmin passes": {
			roles:        []auth.RoleClaim{{ID: "1", Name: "superadmin"}},
			expectedCode: http.StatusOK,
		},
		"admin is blocked": {
			roles:        []auth.RoleClaim{{ID: "1", Name: "admin"}},
			expectedCode: http.StatusForbidden,
		},
		"user is blocked": {
			roles:        []auth.RoleClaim{{ID: "1", Name: "user"}},
			expectedCode: http.StatusForbidden,
		},
		"one matching role among several passes": {
			roles: []auth.RoleClaim{
				{ID: "1", Name: "user"},
				{ID: "2", Name: "superadmin"},
			},
			expectedCode: http.StatusOK,
		},
		"no roles in context is blocked": {
			roles:        nil,
			expectedCode: http.StatusForbidden,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest(http.MethodGet, "/user/list", nil)
			if tc.roles != nil {
				ctx := context.WithValue(req.Context(), "userRoles", tc.roles)
				req = req.WithContext(ctx)
			}
			c.Request = req

			middlewares.RequireRoles(domain.RoleSuperAdmin)(c)
			if !c.IsAborted() {
				c.Status(http.StatusOK)
			}

			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}
