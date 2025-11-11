package user

import (
    "github.com/gin-gonic/gin"

    "api-server/api/middleware"
    "api-server/api/response"
    "api-server/common/menu"
    "api-server/db/pgdb/system"
)

func GetUserMenuList(c *gin.Context) {
	// 获取用户ID
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		response.ReturnError(c, response.UNAUTHENTICATED, "未携带 token")
		return
	}

    // 根据用户ID获取用户角色的菜单与按钮权限
    roleMenus, rolePermissions, err := system.GetUserMenuData(userID)
    if err != nil {
        response.ReturnError(c, response.DATA_LOSS, "查询用户菜单失败")
        return
    }

    // 所有用户（包括超级管理员）均受“租户菜单范围/按钮范围”限制；
    // 超级管理员之所以通常看到全部菜单，是因为默认为其租户配置了全量范围。
    tenantID := middleware.GetTenantID(c)
    if tenantID == 0 {
        response.ReturnError(c, response.UNAUTHENTICATED, "租户信息缺失")
        return
    }
    scopeIDs, err := system.GetTenantMenuScopeIDs(tenantID)
    if err != nil {
        response.ReturnError(c, response.DATA_LOSS, "获取菜单范围失败")
        return
    }
    // 在范围内过滤用户可见菜单与按钮权限
    roleMenus, rolePermissions = system.FilterMenusByIDs(roleMenus, rolePermissions, scopeIDs)

    // 进一步按“租户按钮权限范围”过滤按钮
    authScopeIDs, err := system.GetTenantAuthScopeIDs(tenantID)
    if err != nil {
        response.ReturnError(c, response.DATA_LOSS, "获取按钮权限范围失败")
        return
    }
    if len(authScopeIDs) > 0 {
        // 构造集合以便过滤 rolePermissions
        allowed := make(map[uint]struct{}, len(authScopeIDs))
        for _, id := range authScopeIDs { allowed[id] = struct{}{} }
        filtered := make([]system.SystemMenuAuth, 0, len(rolePermissions))
        for _, a := range rolePermissions {
            if _, ok := allowed[a.ID]; ok { filtered = append(filtered, a) }
        }
        rolePermissions = filtered
    } else {
        // 未配置按钮范围则全部按钮不勾选
        rolePermissions = []system.SystemMenuAuth{}
    }

    // 提取（过滤后的）菜单ID与按钮权限ID
    var roleMenuIds []uint
    for _, m := range roleMenus {
        roleMenuIds = append(roleMenuIds, m.ID)
    }
    var roleAuthIds []uint
    for _, a := range rolePermissions {
        roleAuthIds = append(roleAuthIds, a.ID)
    }

    // 构建带权限标记的菜单树
    menuTree := menu.BuildMenuTreeWithPermission(roleMenus, rolePermissions, roleMenuIds, roleAuthIds, false)
    response.ReturnData(c, menuTree)
}
