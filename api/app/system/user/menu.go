package user

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"api-server/api/middleware"
	"api-server/api/response"
	"api-server/db/pgdb/system"
)

func GetMenuList(c *gin.Context) {
	uID := c.GetString(middleware.JWTDataKey)
	if uID == "" {
		response.ReturnError(c, response.UNAUTHENTICATED, "未携带 token")
		return
	}
	id, err := strconv.ParseUint(uID, 10, 64)
	if err != nil {
		response.ReturnError(c, response.UNAUTHENTICATED, "无效的用户ID")
		return
	}
	// 查询用户菜单
	ms, err := system.GetMenuList(uint(id))
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询用户菜单失败")
		return
	}
	response.ReturnOk(c, ms)
}
