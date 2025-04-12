package user

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"api-server/api/middleware"
	"api-server/api/response"
	"api-server/common/menu"
	"api-server/db/pgdb/system"
)

func GetUserMenuList(c *gin.Context) {
	// 获取用户ID
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

	// 根据用户ID获取用户角色的菜单权限
	roleMenus, rolePermissions, err := system.GetUserMenuData(uint(id))
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询用户菜单失败")
		return
	}

	// 构建菜单树 - 只包含用户有权限的菜单
	menuTree := menu.BuildMenuTree(roleMenus, rolePermissions, false)
	response.ReturnOk(c, menuTree)
}
