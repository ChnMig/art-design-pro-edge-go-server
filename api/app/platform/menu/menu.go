package menu

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"api-server/api/middleware"
	"api-server/api/response"
	commonmenu "api-server/common/menu"
	"api-server/db/pgdb/system"
)

func GetMenuList(c *gin.Context) {
	menus, permissions, err := system.GetMenuData()
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询菜单失败")
		return
	}
	menuTree := commonmenu.BuildMenuTree(menus, permissions, true)
	response.ReturnData(c, menuTree)
}

func AddMenu(c *gin.Context) {
	params := &struct {
		Path          string `json:"path" form:"path" binding:"required"`
		Name          string `json:"name" form:"name" binding:"required"`
		Component     string `json:"component" form:"component"`
		Title         string `json:"title" form:"title" binding:"required"`
		Icon          string `json:"icon" form:"icon"`
		ShowBadge     uint   `json:"showBadge" form:"showBadge"`
		ShowTextBadge string `json:"showTextBadge" form:"showTextBadge"`
		IsHide        uint   `json:"isHide" form:"isHide" binding:"required"`
		IsHideTab     uint   `json:"isHideTab" form:"isHideTab" binding:"required"`
		Link          string `json:"link" form:"link"`
		IsIframe      uint   `json:"isIframe" form:"isIframe" binding:"required"`
		KeepAlive     uint   `json:"keepAlive" form:"keepAlive" binding:"required"`
		IsFirstLevel  uint   `json:"isFirstLevel" form:"isFirstLevel" binding:"required"`
		Status        uint   `json:"status" form:"status" binding:"required"`
		ParentID      uint   `json:"parentId" form:"parentId"`
		Sort          uint   `json:"sort" form:"sort"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	if params.ShowBadge == 0 {
		params.ShowBadge = 2
	}
	var level uint = 1
	if params.ParentID != 0 {
		parentMenu := system.SystemMenu{Model: gorm.Model{ID: params.ParentID}}
		if err := system.GetMenu(&parentMenu); err != nil {
			response.ReturnError(c, response.DATA_LOSS, "父级菜单不存在")
			return
		}
		if parentMenu.Status != system.StatusEnabled {
			response.ReturnError(c, response.DATA_LOSS, "父级菜单已禁用")
			return
		}
		level = parentMenu.Level + 1
	}
	menu := system.SystemMenu{
		Path:          params.Path,
		Name:          params.Name,
		Component:     params.Component,
		Title:         params.Title,
		Icon:          params.Icon,
		ShowBadge:     params.ShowBadge,
		ShowTextBadge: params.ShowTextBadge,
		IsHide:        params.IsHide,
		IsHideTab:     params.IsHideTab,
		Link:          params.Link,
		IsIframe:      params.IsIframe,
		KeepAlive:     params.KeepAlive,
		IsFirstLevel:  params.IsFirstLevel,
		Status:        params.Status,
		Level:         level,
		ParentID:      params.ParentID,
		Sort:          params.Sort,
	}
	if err := system.AddMenu(&menu); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "添加菜单失败")
		return
	}
	response.ReturnData(c, menu)
}

func UpdateMenu(c *gin.Context) {
	params := &struct {
		ID            uint   `json:"id" form:"id" binding:"required"`
		Path          string `json:"path" form:"path" binding:"required"`
		Name          string `json:"name" form:"name" binding:"required"`
		Component     string `json:"component" form:"component"`
		Title         string `json:"title" form:"title" binding:"required"`
		Icon          string `json:"icon" form:"icon"`
		ShowBadge     uint   `json:"showBadge" form:"showBadge"`
		ShowTextBadge string `json:"showTextBadge" form:"showTextBadge"`
		IsHide        uint   `json:"isHide" form:"isHide" binding:"required"`
		IsHideTab     uint   `json:"isHideTab" form:"isHideTab" binding:"required"`
		Link          string `json:"link" form:"link"`
		IsIframe      uint   `json:"isIframe" form:"isIframe" binding:"required"`
		KeepAlive     uint   `json:"keepAlive" form:"keepAlive" binding:"required"`
		IsFirstLevel  uint   `json:"isFirstLevel" form:"isFirstLevel" binding:"required"`
		Status        uint   `json:"status" form:"status" binding:"required"`
		ParentID      uint   `json:"parentId" form:"parentId"`
		Sort          uint   `json:"sort" form:"sort"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	if params.ShowBadge == 0 {
		params.ShowBadge = 2
	}
	var level uint = 1
	if params.ParentID != 0 {
		parent := system.SystemMenu{Model: gorm.Model{ID: params.ParentID}}
		if err := system.GetMenu(&parent); err != nil {
			response.ReturnError(c, response.DATA_LOSS, "父级菜单不存在")
			return
		}
		if parent.Status != system.StatusEnabled {
			response.ReturnError(c, response.DATA_LOSS, "父级菜单已禁用")
			return
		}
		level = parent.Level + 1
	}
	if params.Status == system.StatusDisabled {
		children, _, err := system.FindMenuList(&system.SystemMenu{ParentID: params.ID}, -1, -1)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			response.ReturnError(c, response.DATA_LOSS, "查询子菜单失败")
			return
		}
		for _, child := range children {
			if child.Status == system.StatusEnabled {
				response.ReturnError(c, response.DATA_LOSS, "请先禁用子菜单")
				return
			}
		}
	}
	menu := system.SystemMenu{
		Model:         gorm.Model{ID: params.ID},
		Path:          params.Path,
		Name:          params.Name,
		Component:     params.Component,
		Title:         params.Title,
		Icon:          params.Icon,
		ShowBadge:     params.ShowBadge,
		ShowTextBadge: params.ShowTextBadge,
		IsHide:        params.IsHide,
		IsHideTab:     params.IsHideTab,
		Link:          params.Link,
		IsIframe:      params.IsIframe,
		KeepAlive:     params.KeepAlive,
		IsFirstLevel:  params.IsFirstLevel,
		Status:        params.Status,
		Level:         level,
		ParentID:      params.ParentID,
		Sort:          params.Sort,
	}
	if err := system.UpdateMenu(&menu); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "更新菜单失败")
		return
	}
	response.ReturnData(c, menu)
}

func DeleteMenu(c *gin.Context) {
	params := &struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	menu := system.SystemMenu{Model: gorm.Model{ID: params.ID}}
	if err := system.GetMenu(&menu); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ReturnError(c, response.DATA_LOSS, "菜单不存在")
			return
		}
		response.ReturnError(c, response.DATA_LOSS, "查询菜单失败")
		return
	}
	children, _, err := system.FindMenuList(&system.SystemMenu{ParentID: menu.ID}, -1, -1)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.ReturnError(c, response.DATA_LOSS, "查询子菜单失败")
		return
	}
	if len(children) > 0 {
		response.ReturnError(c, response.DATA_LOSS, "请先删除子菜单")
		return
	}
	if err := system.DeleteMenu(&menu); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "删除菜单失败")
		return
	}
	response.ReturnData(c, menu)
}

func GetMenuAuthList(c *gin.Context) {
	params := &struct {
		MenuID uint `json:"menu_id" form:"menu_id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	auth := system.SystemMenuAuth{MenuID: params.MenuID}
	auths, err := system.FindMenuAuthList(&auth)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询菜单权限失败")
		return
	}
	response.ReturnData(c, auths)
}

func AddMenuAuth(c *gin.Context) {
	params := &struct {
		MenuID uint   `json:"menu_id"`
		Mark   string `json:"mark"`
		Title  string `json:"title"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	auth := system.SystemMenuAuth{
		MenuID: params.MenuID,
		Mark:   params.Mark,
		Title:  params.Title,
	}
	if err := system.AddMenuAuth(&auth); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "添加菜单权限失败")
		return
	}
	response.ReturnData(c, auth)
}

func UpdateMenuAuth(c *gin.Context) {
	params := &struct {
		ID     uint   `json:"id"`
		Title  string `json:"title"`
		Mark   string `json:"mark"`
		MenuID uint   `json:"menu_id"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	auth := system.SystemMenuAuth{
		Model:  gorm.Model{ID: params.ID},
		Title:  params.Title,
		Mark:   params.Mark,
		MenuID: params.MenuID,
	}
	if err := system.UpdateMenuAuth(&auth); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "更新菜单权限失败")
		return
	}
	response.ReturnData(c, auth)
}

func DeleteMenuAuth(c *gin.Context) {
	params := &struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	auth := system.SystemMenuAuth{Model: gorm.Model{ID: params.ID}}
	if err := system.DeleteMenuAuth(&auth); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "删除菜单权限失败")
		return
	}
	response.ReturnData(c, auth)
}
