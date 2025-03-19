package user

import (
	"github.com/gin-gonic/gin"

	"api-server/api/response"
	"api-server/common/menu"
	"api-server/db/pgdb/system"
)

func GetUserMenuList(c *gin.Context) {
	// 查询用户菜单数据
	menus, permissions, err := system.GetMenuData()
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询用户菜单失败")
		return
	}
	// 构建菜单树
	menuTree := menu.BuildMenuTree(menus, permissions)
	response.ReturnOk(c, menuTree)
}
