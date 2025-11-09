package health

import "github.com/gin-gonic/gin"

// RegisterOpenRoutes 提供健康检查开放接口
// GET /api/v1/open/health
func RegisterOpenRoutes(open *gin.RouterGroup) {
	if open == nil {
		return
	}
	open.GET("/health", Status)
}
