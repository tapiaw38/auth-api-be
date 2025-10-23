package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

func NewChangePasswordHandler(usecase user.ChangePasswordUsecase) func(c *gin.Context) {
	return func(c *gin.Context) {
		var request user.ChangePasswordInput

		if err := c.ShouldBindJSON(&request); err != nil {
			appErr := apperrors.NewApplicationError(mappings.RequestBodyParsingError, err)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		username, ok := c.Request.Context().Value("userID").(string)
		if !ok || username == "" {
			appErr := apperrors.NewApplicationError(mappings.AuthUnauthorizedError, nil)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		if appErr := usecase.Execute(c, request, username); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
