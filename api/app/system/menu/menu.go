package menu

import (
	"github.com/gin-gonic/gin"

	"api-server/api/middleware"
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

func AddMenu(c *gin.Context) {
	params := &struct {
		Path              string `json:"path" form:"path" binding:"required"`
		Name              string `json:"name" form:"name" binding:"required"`
		Compent           string `json:"compent" form:"compent" binding:"required"`
		Title             string `json:"title" form:"title" binding:"required"`
		Icon              string `json:"icon" form:"icon" binding:"required"`
		ShowBadge         uint   `json:"show_badge" form:"show_badge" binding:"required"`
		ShowTextBadge     string `json:"show_text_badge" form:"show_text_badge" binding:"required"`
		IsHide            uint   `json:"is_hide" form:"is_hide" binding:"required"`
		IsHideTab         uint   `json:"is_hide_tab" form:"is_hide_tab" binding:"required"`
		Link              string `json:"link" form:"link" binding:"required"`
		IsIframe          uint   `json:"is_iframe" form:"is_iframe" binding:"required"`
		KeepAlive         uint   `json:"keep_alive" form:"keep_alive" binding:"required"`
		IsInMainContainer uint   `json:"is_in_main_container" form:"is_in_main_container" binding:"required"`
		Status            uint   `json:"status" form:"status" binding:"required"`
		ParentID          uint   `json:"parent_id" form:"parent_id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	var level uint = 1
	// 如果有父级ID, 则查询父级ID是否存在
	if params.ParentID != 0 {
		parentMenu, err := system.GetMenuByID(params.ParentID)
		if err != nil {
			response.ReturnError(c, response.DATA_LOSS, "父级菜单不存在")
			return
		}
		if parentMenu.Status != 1 {
			response.ReturnError(c, response.DATA_LOSS, "父级菜单已禁用")
			return
		}
		level = parentMenu.Level + 1
	}
	menu := system.Menu{
		Path:              params.Path,
		Name:              params.Name,
		Component:         params.Compent,
		Title:             params.Title,
		Icon:              params.Icon,
		ShowBadge:         params.ShowBadge,
		ShowTextBadge:     params.ShowTextBadge,
		IsHide:            params.IsHide,
		IsHideTab:         params.IsHideTab,
		Link:              params.Link,
		IsIframe:          params.IsIframe,
		KeepAlive:         params.KeepAlive,
		IsInMainContainer: params.IsInMainContainer,
		Status:            params.Status,
		Level:             level,
		ParentID:          params.ParentID,
	}
	if err := system.AddMenu(&menu); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "添加菜单失败")
		return
	}
	response.ReturnOk(c, menu)
}
