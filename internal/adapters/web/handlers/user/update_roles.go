package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

type updateRolesInput struct {
	Roles []string `json:"roles"`
}

func NewUpdateRolesHandler(usecase user.UpdateRolesUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, ok := c.Request.Context().Value("userID").(string)
		if !ok || username == "" {
			appErr := apperrors.NewApplicationError(mappings.AuthUnauthorizedError, nil)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		var input updateRolesInput
		if err := c.ShouldBindJSON(&input); err != nil {
			appErr := apperrors.NewApplicationError(mappings.RequestBodyParsingError, err)
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		output, appErr := usecase.Execute(c, user.UpdateRolesInput{
			TargetID:      c.Param("id"),
			ActorUsername: username,
			Roles:         input.Roles,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
