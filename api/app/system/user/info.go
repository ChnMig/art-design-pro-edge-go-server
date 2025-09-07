package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"api-server/api/middleware"
	"api-server/api/response"
	"api-server/db/pgdb/system"
)

func UpdateUserInfo(c *gin.Context) {
	params := &struct {
		Password string `json:"password" form:"password"`
		Name     string `json:"name" form:"name" binding:"required"`
		Phone    string `json:"phone" form:"phone" binding:"required"`
		Gender   uint   `json:"gender" form:"gender" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		response.ReturnError(c, response.UNAUTHENTICATED, "未携带 token")
		return
	}
	u := system.SystemUser{
		Model:  gorm.Model{ID: userID},
		Name:   params.Name,
		Phone:  params.Phone,
		Gender: params.Gender,
	}
	if params.Password != "" {
		u.Password = params.Password
	}
	if err := system.UpdateUser(&u); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "更新用户失败")
		return
	}
	response.ReturnData(c, "更新用户成功")
}

func GetUserInfo(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		response.ReturnError(c, response.UNAUTHENTICATED, "未携带 token")
		return
	}
	user := system.SystemUser{Model: gorm.Model{ID: userID}}
	if err := system.GetUser(&user); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询用户失败")
		return
	}
	user.Password = "" // 不返回密码
	response.ReturnData(c, user)
}
