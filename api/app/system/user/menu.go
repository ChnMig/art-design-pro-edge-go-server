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

	// 提取用户拥有的菜单ID（因为用户已经有所有这些角色菜单的权限）
	var roleMenuIds []uint
	for _, m := range roleMenus {
		roleMenuIds = append(roleMenuIds, m.ID)
	}

	// 提取用户拥有的权限ID
	var roleAuthIds []uint
	for _, a := range rolePermissions {
		roleAuthIds = append(roleAuthIds, a.ID)
	}

	// 构建菜单树 - 使用带权限标记的菜单树构建函数
	menuTree := menu.BuildMenuTreeWithPermission(roleMenus, rolePermissions, roleMenuIds, roleAuthIds, false)
	response.ReturnData(c, menuTree)
}
