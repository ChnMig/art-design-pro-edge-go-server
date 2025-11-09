package private

import (
	"github.com/gin-gonic/gin"

	"api-server/api/app/v1/private/admin"
)

// RegisterRoutes 将所有内部业务接口挂载到 /api/v1/private
func RegisterRoutes(private *gin.RouterGroup) {
	if private == nil {
		return
	}

	adminGroup := private.Group("/admin")
	admin.RegisterRoutes(adminGroup)
}
