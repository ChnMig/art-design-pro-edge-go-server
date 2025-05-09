package role

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"api-server/api/middleware"
	"api-server/api/response"
	"api-server/db/pgdb/system"
)

func GetRoleList(c *gin.Context) {
	params := &struct {
		Name   string `json:"name" form:"name"`
		Status uint   `json:"status" form:"status"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	// 获取分页参数
	page := middleware.GetPage(c)
	pageSize := middleware.GetPageSize(c)

	role := system.SystemRole{
		Name:   params.Name,
		Status: params.Status,
	}

	// 调用带分页的查询函数
	roles, total, err := system.FindRoleList(&role, page, pageSize)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "获取角色列表失败")
		return
	}

	// 返回带总数的结果
	response.ReturnDataWithCount(c, int(total), roles)
}

func AddRole(c *gin.Context) {
	params := &struct {
		Name   string `json:"name" form:"name" binding:"required"`
		Status int    `json:"status" form:"status" binding:"required"`
		Desc   string `json:"desc" form:"desc"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	role := system.SystemRole{
		Name:   params.Name,
		Status: uint(params.Status),
		Desc:   params.Desc,
	}
	err := system.AddRole(&role)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "添加角色失败")
		return
	}
	response.ReturnData(c, role)
}

func UpdateRole(c *gin.Context) {
	params := &struct {
		ID     uint   `json:"id" form:"id" binding:"required"`
		Name   string `json:"name" form:"name" binding:"required"`
		Status int    `json:"status" form:"status" binding:"required"`
		Desc   string `json:"desc" form:"desc"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	role := system.SystemRole{
		Model:  gorm.Model{ID: params.ID},
		Name:   params.Name,
		Status: uint(params.Status),
		Desc:   params.Desc,
	}
	err := system.UpdateRole(&role)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "更新角色失败")
		return
	}
	response.ReturnData(c, role)
}

func DeleteRole(c *gin.Context) {
	params := &struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	role := system.SystemRole{
		Model: gorm.Model{ID: params.ID},
	}
	err := system.DeleteRole(&role)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "删除角色失败")
		return
	}
	response.ReturnData(c, role)
}
