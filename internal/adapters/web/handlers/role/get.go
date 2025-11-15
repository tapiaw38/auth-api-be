package role

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api-be/internal/usecases/role"
)

func NewGetHandler(usecase role.GetUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		name := c.Query("name")

		filters := role.GetFilterOptions{}
		if id != "" {
			filters.ID = id
		}
		if name != "" {
			filters.Name = name
		}

		output, appErr := usecase.Execute(c, filters)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		if output == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "role not found",
			})
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
