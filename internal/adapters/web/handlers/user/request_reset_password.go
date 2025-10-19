package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

type RequestResetPasswordInput struct {
	Email string `json:"email"`
}

func NewRequestResetPasswordHandler(usecase user.RequestResetPasswordUsecase) func(c *gin.Context) {
	return func(c *gin.Context) {
		var input RequestResetPasswordInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request format",
				"error":   err.Error(),
			})
			return
		}

		if input.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Email is required",
			})
			return
		}

		output, err := usecase.Execute(c, input.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
