package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

func NewMeHandler(usecase user.GetUsecase) func(c *gin.Context) {
	return func(c *gin.Context) {

		username, ok := c.Request.Context().Value("userID").(string)
		if !ok || username == "" {
			appErr := apperrors.NewApplicationError(mappings.AuthUnauthorizedError, nil)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		filter := user.GetFilterOptions{
			Username: username,
		}

		userOutput, appErr := usecase.Execute(c, filter)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, userOutput)
	}
}
