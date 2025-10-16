package role

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api-be/internal/usecases/role"
)

type UpdateInput struct {
	Name string `json:"name" binding:"required"`
}

func NewUpdateHandler(usecase role.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "role ID is required",
			})
			return
		}

		var input UpdateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request body",
				"error":   err.Error(),
			})
			return
		}

		output, err := usecase.Execute(c, id, role.UpdateInput{
			Name: input.Name,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, output)
	}
}