package user

import (
	"github.com/gin-gonic/gin"

	"api-server/api/middleware"
	"api-server/api/response"
	menudomain "api-server/domain/admin/menu"
)

func GetUserMenuList(c *gin.Context) {
	// 获取用户ID
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		response.ReturnError(c, response.UNAUTHENTICATED, "未携带 token")
		return
	}

	// 所有用户（包括超级管理员）均受“租户菜单范围/按钮范围”限制；
	// 超级管理员之所以通常看到全部菜单，是因为默认为其租户配置了全量范围。
	tenantID := middleware.GetTenantID(c)
	if tenantID == 0 {
		response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
		return
	}

	// 构建带权限标记的菜单树
	menuTree, err := menudomain.GetUserMenuTree(userID, tenantID)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询用户菜单失败")
		return
	}
	response.ReturnData(c, menuTree)
}
