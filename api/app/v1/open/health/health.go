package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

// Status 返回服务健康状态
func Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"ready":  true,
		"uptime": time.Since(startTime).String(),
	})
}
