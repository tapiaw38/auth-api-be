package role

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tapiaw38/auth-api-be/internal/usecases/role"
)

func NewEnsureHandler(usecase role.EnsureUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		if appErr := usecase.Execute(c); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	}
}
