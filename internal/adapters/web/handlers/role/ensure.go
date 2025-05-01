package role

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tapiaw38/auth-api-be/internal/usecases/role"
)

func NewEnsureHandler(usecase role.EnsureUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := usecase.Execute(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	}
}
