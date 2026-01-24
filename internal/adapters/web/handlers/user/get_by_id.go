package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

func NewGetByIDHandler(usecase user.GetUsecase) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			appErr := apperrors.NewApplicationError(mappings.InvalidParamsError, nil)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		filter := user.GetFilterOptions{
			Username: id,
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
