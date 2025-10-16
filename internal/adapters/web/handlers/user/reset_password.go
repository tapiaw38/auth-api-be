package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

func NewResetPasswordHandler(usecase user.ResetPasswordUsecase) func(c *gin.Context) {
	return func(c *gin.Context) {
		var input user.ResetPasswordInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request format",
				"error":   err.Error(),
			})
			return
		}

		// Validate token
		if input.Token == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Password reset token is required",
			})
			return
		}

		// Validate password basic requirements
		if input.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Password is required",
			})
			return
		}

		if len(input.Password) < 8 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Password must be at least 8 characters long",
			})
			return
		}

		output, err := usecase.Execute(c, input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
