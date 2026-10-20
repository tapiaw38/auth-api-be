package user

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api-be/internal/usecases/user"
)

func NewBatchHandler(usecase user.BatchUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ids []string
		for _, id := range strings.Split(c.Query("ids"), ",") {
			id = strings.TrimSpace(id)
			if id != "" {
				ids = append(ids, id)
			}
		}

		output, appErr := usecase.Execute(c, user.BatchInput{IDs: ids})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
