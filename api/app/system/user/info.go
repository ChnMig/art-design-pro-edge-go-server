package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

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
	user := system.User{Model: gorm.Model{ID: uint(id)}}
	if err := system.GetUser(&user); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询用户失败")
		return
	}
	response.ReturnOk(c, gin.H{
		"id":       user.ID,
		"name":     user.Name,
		"username": user.Username,
		"avatar":   "",
		"email":    "",
	})
}
