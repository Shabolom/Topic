package tools

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetQueryPagination возвращает значения из мапы query param по ключам page и limit
func GetQueryPagination(c *gin.Context) (int, int, error) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		return 0, 0, err
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		return 0, 0, err
	}

	return page, limit, nil
}
