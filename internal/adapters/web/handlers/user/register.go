package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

func NewRegisterHandler(usecase user.RegisterUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user user.RegisterInput
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		output, err := usecase.Execute(c, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
