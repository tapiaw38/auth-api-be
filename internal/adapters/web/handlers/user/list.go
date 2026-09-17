package user

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

func NewListHandler(usecase user.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		filters := user.ListFilterOptions{
			RoleName: c.Query("role"),
			RoleID:   c.Query("role_id"),
			Search:   strings.TrimSpace(c.Query("search")),
		}

		if limitStr := c.Query("limit"); limitStr != "" {
			if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
				filters.Limit = limit
			}
		}

		if offsetStr := c.Query("offset"); offsetStr != "" {
			if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
				filters.Offset = offset
			}
		}

		if filters.Search != "" && (filters.Limit == 0 || filters.Limit > 20) {
			filters.Limit = 20
		}

		users, appErr := usecase.Execute(c, filters)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, user.ListOutput{Data: users})
	}
}
