package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
)

func GetPage(c *gin.Context) int {
	p, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		return defaultPage
	}
	return p
}

func GetPageSize(c *gin.Context) int {
	ps, err := strconv.Atoi(c.Query("page_size"))
	if err != nil {
		return defaultPageSize
	}
	return ps
}
