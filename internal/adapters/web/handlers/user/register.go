package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

func NewRegisterHandler(usecase user.RegisterUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user user.RegisterInput
		if err := c.ShouldBindJSON(&user); err != nil {
			appErr := apperrors.NewApplicationError(mappings.RequestBodyParsingError, err)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		output, appErr := usecase.Execute(c, user)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
