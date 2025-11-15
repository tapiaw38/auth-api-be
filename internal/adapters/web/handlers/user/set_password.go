package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

type SetPasswordRequest struct {
	NewPassword string `json:"new_password"`
}

func NewSetPasswordHandler(usecase user.SetPasswordUsecase) func(c *gin.Context) {
	return func(c *gin.Context) {
		var request SetPasswordRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			appErr := apperrors.NewApplicationError(mappings.RequestBodyParsingError, err)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		username := c.Request.Context().Value("userID").(string)

		input := user.SetPasswordInput{
			Username:    username,
			NewPassword: request.NewPassword,
		}

		if appErr := usecase.Execute(c, input); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
