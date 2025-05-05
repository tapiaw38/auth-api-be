package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

func NewResetPasswordHandler(usecase user.ResetPasswordUsecase) func(c *gin.Context) {
	return func(c *gin.Context) {
		var input InputResetPassword
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		if input.Token == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "token is required",
			})
			return
		}

		if input.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "password is required",
			})
			return
		}

		output, err := usecase.Execute(c, input.Token, input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
