package v1

import (
	"github.com/gin-gonic/gin"

	"api-server/api/app/v1/open"
	"api-server/api/app/v1/private"
)

// RegisterRoutes 在 /api/v1 下注册 open/private 等分组
func RegisterRoutes(v1 *gin.RouterGroup) {
	if v1 == nil {
		return
	}

	openGroup := v1.Group("/open")
	open.RegisterRoutes(openGroup)

	privateGroup := v1.Group("/private")
	private.RegisterRoutes(privateGroup)
}
