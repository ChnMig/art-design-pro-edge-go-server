package user

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"api-server/api/middleware"
	"api-server/api/response"
	"api-server/db/pgdb/system"
)

func GetUserInfo(c *gin.Context) {
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
	u, err := system.GetUser(uint(id))
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询用户失败")
		return
	}
	response.ReturnOk(c, gin.H{
		"id":       u.ID,
		"name":     u.Username,
		"username": u.Username,
		"avatar":   "",
		"email":    "",
	})
}
