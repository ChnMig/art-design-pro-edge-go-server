package open

import (
	"github.com/gin-gonic/gin"

	"api-server/api/app/v1/open/health"
)

// RegisterRoutes 注册所有开放路由
func RegisterRoutes(open *gin.RouterGroup) {
	if open == nil {
		return
	}
	health.RegisterOpenRoutes(open)
}
