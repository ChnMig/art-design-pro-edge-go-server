package menu

import (
	"github.com/gin-gonic/gin"

	"api-server/api/response"
	"api-server/common/menu"
	"api-server/db/pgdb/system"
)

func GetMenuList(c *gin.Context) {
	// 查询菜单数据
	menus, menup, err := system.GetMenuData()
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询菜单失败")
		return
	}
	// 构建菜单树
	menuTree := menu.BuildMenuTree(menus, menup)
	response.ReturnOk(c, menuTree)
}
