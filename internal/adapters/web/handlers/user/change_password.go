
package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func NewChangePasswordHandler(usecase user.ChangePasswordUsecase) func(c *gin.Context) {
	return func(c *gin.Context) {
		var request ChangePasswordRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		username := c.Request.Context().Value("userID").(string)

		input := user.ChangePasswordInput{
			Username:    username,
			OldPassword: request.OldPassword,
			NewPassword: request.NewPassword,
		}

		if err := usecase.Execute(c, input); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
