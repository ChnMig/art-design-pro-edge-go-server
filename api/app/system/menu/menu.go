package menu

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

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
	menuTree := menu.BuildMenuTree(menus, menup, true)
	response.ReturnOk(c, menuTree)
}

func DeleteMenu(c *gin.Context) {
	params := &struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	menu := system.Menu{Model: gorm.Model{ID: params.ID}}
	err := system.GetMenu(&menu)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ReturnError(c, response.DATA_LOSS, "菜单不存在")
			return
		}
		response.ReturnError(c, response.DATA_LOSS, "查询菜单失败")
		return
	}
	children, err := system.FindMenuList(&system.Menu{ParentID: menu.ID})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			response.ReturnError(c, response.DATA_LOSS, "查询子菜单失败")
			return
		}
	}
	if len(children) > 0 {
		response.ReturnError(c, response.DATA_LOSS, "请先删除子菜单")
		return
	}
	if err := system.DeleteMenu(&menu); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "删除菜单失败")
		return
	}
	response.ReturnOk(c, menu)
}

func AddMenu(c *gin.Context) {
	params := &struct {
		Path              string `json:"path" form:"path" binding:"required"`
		Name              string `json:"name" form:"name" binding:"required"`
		Component         string `json:"component" form:"component"`
		Title             string `json:"title" form:"title" binding:"required"`
		Icon              string `json:"icon" form:"icon"`
		ShowBadge         uint   `json:"showBadge" form:"showBadge"`
		ShowTextBadge     string `json:"showTextBadge" form:"showTextBadge"`
		IsHide            uint   `json:"isHide" form:"isHide" binding:"required"`
		IsHideTab         uint   `json:"isHideTab" form:"isHideTab" binding:"required"`
		Link              string `json:"link" form:"link"`
		IsIframe          uint   `json:"isIframe" form:"isIframe" binding:"required"`
		KeepAlive         uint   `json:"keepAlive" form:"keepAlive" binding:"required"`
		IsInMainContainer uint   `json:"isInMainContainer" form:"isInMainContainer" binding:"required"`
		Status            uint   `json:"status" form:"status" binding:"required"`
		ParentID          uint   `json:"parentId" form:"parentId"`
		Sort              uint   `json:"sort" form:"sort"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	if params.ShowBadge == 0 {
		params.ShowBadge = 2
	}
	var level uint = 1
	// 如果有父级ID, 则查询父级ID是否存在
	if params.ParentID != 0 {
		parentMenu := &system.Menu{}
		err := system.GetMenu(&system.Menu{Model: gorm.Model{ID: params.ParentID}})
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
		Component:         params.Component,
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
		Sort:              params.Sort,
	}
	if err := system.AddMenu(&menu); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "添加菜单失败")
		return
	}
	response.ReturnOk(c, menu)
}

func UpdateMenu(c *gin.Context) {
	params := &struct {
		ID                uint   `json:"id" form:"id" binding:"required"`
		Path              string `json:"path" form:"path" binding:"required"`
		Name              string `json:"name" form:"name" binding:"required"`
		Component         string `json:"component" form:"component"`
		Title             string `json:"title" form:"title" binding:"required"`
		Icon              string `json:"icon" form:"icon"`
		ShowBadge         uint   `json:"showBadge" form:"showBadge"`
		ShowTextBadge     string `json:"showTextBadge" form:"showTextBadge"`
		IsHide            uint   `json:"isHide" form:"isHide" binding:"required"`
		IsHideTab         uint   `json:"isHideTab" form:"isHideTab" binding:"required"`
		Link              string `json:"link" form:"link"`
		IsIframe          uint   `json:"isIframe" form:"isIframe" binding:"required"`
		KeepAlive         uint   `json:"keepAlive" form:"keepAlive" binding:"required"`
		IsInMainContainer uint   `json:"isInMainContainer" form:"isInMainContainer" binding:"required"`
		Status            uint   `json:"status" form:"status" binding:"required"`
		ParentID          uint   `json:"parentId" form:"parentId"`
		Sort              uint   `json:"sort" form:"sort"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	if params.ShowBadge == 0 {
		params.ShowBadge = 2
	}
	var level uint = 1
	// 如果有父级ID, 则查询父级ID是否存在
	if params.ParentID != 0 {
		parent := system.Menu{Model: gorm.Model{ID: params.ParentID}}
		err := system.GetMenu(&parent)
		if err != nil {
			response.ReturnError(c, response.DATA_LOSS, "父级菜单不存在")
			return
		}
		if parent.Status != 1 {
			response.ReturnError(c, response.DATA_LOSS, "父级菜单已禁用")
			return
		}
		level = parent.Level + 1
	}
	if params.Status == 2 {
		// 判断子菜单是否是禁用状态
		children, err := system.FindMenuList(&system.Menu{ParentID: params.ID})
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				response.ReturnError(c, response.DATA_LOSS, "查询子菜单失败")
				return
			}
		}
		if len(children) > 0 {
			for _, v := range children {
				if v.Status == 1 {
					response.ReturnError(c, response.DATA_LOSS, "请先禁用子菜单")
					return
				}
			}
		}
	}
	menu := system.Menu{
		Model:             gorm.Model{ID: params.ID},
		Path:              params.Path,
		Name:              params.Name,
		Component:         params.Component,
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
		Sort:              params.Sort,
	}
	if err := system.UpdateMenu(&menu); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "更新菜单失败")
		return
	}
	response.ReturnOk(c, menu)
}
