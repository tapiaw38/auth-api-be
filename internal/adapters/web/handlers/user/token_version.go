package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

func NewGetTokenVersionHandler(usecase user.GetTokenVersionUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		version, err := usecase.Execute(c, c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "user:token-version:error",
				"message": "could not read token version",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": gin.H{"token_version": version}})
	}
}
