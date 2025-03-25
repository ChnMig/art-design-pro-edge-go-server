package menu

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"api-server/api/middleware"
	"api-server/api/response"
	"api-server/db/pgdb/system"
)

func AddMenuAuth(c *gin.Context) {
	params := &struct {
		MenuID uint   `json:"menu_id"`
		Mark   string `json:"mark"` // 标识
		Title  string `json:"title"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	auth := system.MenuAuth{
		MenuID: params.MenuID,
		Mark:   params.Mark,
		Title:  params.Title,
	}
	if err := system.AddMenuAuth(&auth); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "添加菜单权限失败")
		return
	}
	response.ReturnOk(c, auth)
}

func UpdateMenuAuth(c *gin.Context) {
	params := &struct {
		ID     uint   `json:"id"`
		Title  string `json:"title"`
		Mark   string `json:"mark"` // 标识
		MenuID uint   `json:"menu_id"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	auth := system.MenuAuth{
		Model:  gorm.Model{ID: params.ID},
		Title:  params.Title,
		Mark:   params.Mark,
		MenuID: params.MenuID,
	}
	if err := system.UpdateMenuAuth(&auth); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "更新菜单权限失败")
		return
	}
	response.ReturnOk(c, auth)
}

func DeleteMenuAuth(c *gin.Context) {
	params := &struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	auth := system.MenuAuth{Model: gorm.Model{ID: params.ID}}
	if err := system.DeleteMenuAuth(&auth); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "删除菜单权限失败")
		return
	}
	response.ReturnOk(c, auth)
}

func GetMenuAuthList(c *gin.Context) {
	params := &struct {
		MenuID uint `json:"menu_id" form:"menu_id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	auth := system.MenuAuth{MenuID: params.MenuID}
	auths, err := system.FindMenuAuthList(&auth)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询菜单权限失败")
		return
	}
	response.ReturnOk(c, auths)
}
