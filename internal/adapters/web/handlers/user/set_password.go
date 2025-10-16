package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

type SetPasswordRequest struct {
	NewPassword string `json:"new_password"`
}

func NewSetPasswordHandler(usecase user.SetPasswordUsecase) func(c *gin.Context) {
	return func(c *gin.Context) {
		var request SetPasswordRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		username := c.Request.Context().Value("userID").(string)

		input := user.SetPasswordInput{
			Username:    username,
			NewPassword: request.NewPassword,
		}

		if err := usecase.Execute(c, input); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
