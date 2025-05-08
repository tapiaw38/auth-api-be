package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

func NewMeHandler(usecase user.GetUsecase) func(c *gin.Context) {
	return func(c *gin.Context) {

		username := c.Request.Context().Value("userID").(string)

		userOutput, err := usecase.Execute(c, user.GetFilterOptions{
			Username: username,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, userOutput)
	}
}
