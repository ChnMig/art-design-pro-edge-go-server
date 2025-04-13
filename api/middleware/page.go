package middleware

import (
	"strconv"

	"api-server/config"

	"github.com/gin-gonic/gin"
)

func GetPage(c *gin.Context) int {
	p, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		return config.DefaultPage
	}
	return p
}

func GetPageSize(c *gin.Context) int {
	ps, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		return config.DefaultPageSize
	}
	return ps
}
