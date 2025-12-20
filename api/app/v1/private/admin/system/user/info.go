package user

import (
	"github.com/gin-gonic/gin"

	"api-server/api/middleware"
	"api-server/api/response"
	userdomain "api-server/domain/admin/user"
)

func UpdateUserInfo(c *gin.Context) {
	params := &struct {
		Password string `json:"password" form:"password"`
		Username string `json:"username" form:"username" binding:"required"`
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
	if err := userdomain.UpdateUserProfile(userdomain.UpdateProfileInput{
		UserID:   userID,
		Password: params.Password,
		Username: params.Username,
		Phone:    params.Phone,
		Gender:   params.Gender,
	}); err != nil {
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
	user, err := userdomain.GetUserProfile(userID)
	if err != nil {
		ReturnDomainError(c, err, "查询用户失败")
		return
	}
	response.ReturnData(c, user)
}
