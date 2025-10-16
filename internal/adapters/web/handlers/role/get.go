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

		output, err := usecase.Execute(c, filters)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
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
