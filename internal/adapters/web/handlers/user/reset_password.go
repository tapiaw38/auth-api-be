package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

func NewResetPasswordHandler(usecase user.ResetPasswordUsecase) func(c *gin.Context) {
	return func(c *gin.Context) {
		var input user.ResetPasswordInput
		if err := c.ShouldBindJSON(&input); err != nil {
			appErr := apperrors.NewApplicationError(mappings.RequestBodyParsingError, err)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		// Validate token
		if input.Token == "" {
			appErr := apperrors.NewApplicationError(mappings.UserResetPasswordTokenRequiredError, nil)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		// Validate password basic requirements
		if input.Password == "" {
			appErr := apperrors.NewApplicationError(mappings.UserResetPasswordPasswordRequiredError, nil)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		if len(input.Password) < 8 {
			appErr := apperrors.NewApplicationError(mappings.UserResetPasswordPasswordTooShortError, nil)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		output, appErr := usecase.Execute(c, input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
