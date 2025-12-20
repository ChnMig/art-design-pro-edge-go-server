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
	keys := []string{"page_size", "pageSize"}
	for _, key := range keys {
		raw := c.Query(key)
		if raw == "" {
			continue
		}
		ps, err := strconv.Atoi(raw)
		if err == nil {
			return ps
		}
	}
	return config.DefaultPageSize
}
