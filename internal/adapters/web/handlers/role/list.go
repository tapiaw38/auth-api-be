package role

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api-be/internal/usecases/role"
)

func NewListHandler(usecase role.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		filter := parseListFilter(c.Request.URL.Query())
		output, err := usecase.Execute(c, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

func parseListFilter(queries url.Values) role.ListFilterOptions {
	return role.ListFilterOptions{
		Name: queries.Get("name"),
	}
}
