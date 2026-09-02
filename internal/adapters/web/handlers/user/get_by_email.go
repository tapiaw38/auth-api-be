package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

// NewGetByEmailHandler exists for trusted server-to-server lookups (e.g.
// practiq-campus-be checking whether a shared identity already exists
// before creating one) — superadmin-only, same as user/list.
func NewGetByEmailHandler(usecase user.GetUsecase) func(c *gin.Context) {
	return func(c *gin.Context) {
		email := c.Query("email")
		if email == "" {
			appErr := apperrors.NewApplicationError(mappings.InvalidParamsError, nil)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		filter := user.GetFilterOptions{
			Email: email,
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
