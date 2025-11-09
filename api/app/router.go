package app

import (
	"github.com/gin-gonic/gin"

	v1 "api-server/api/app/v1"
)

// RegisterRoutes 在 /api 下挂载各版本路由
func RegisterRoutes(api *gin.RouterGroup) {
	if api == nil {
		return
	}
	v1Group := api.Group("/v1")
	v1.RegisterRoutes(v1Group)
}
