package department

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"api-server/api/middleware"
	"api-server/api/response"
	"api-server/db/pgdb/system"
)

func AddDepartment(c *gin.Context) {
	params := &struct {
		Name   string `json:"name" form:"name" binding:"required"`
		Status int    `json:"status" form:"status" binding:"required"`
		Sort   int    `json:"sort" form:"sort"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	depatment := system.SystemDepartment{
		Name:   params.Name,
		Status: uint(params.Status),
		Sort:   uint(params.Sort),
	}
	err := system.AddDepartment(&depatment)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "添加部门失败")
		return
	}
	response.ReturnData(c, depatment)
}

func UpdateDepartment(c *gin.Context) {
	params := &struct {
		ID     uint   `json:"id" form:"id" binding:"required"`
		Name   string `json:"name" form:"name" binding:"required"`
		Status int    `json:"status" form:"status" binding:"required"`
		Sort   int    `json:"sort" form:"sort"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	department := system.SystemDepartment{
		Model:  gorm.Model{ID: params.ID},
		Name:   params.Name,
		Status: uint(params.Status),
		Sort:   uint(params.Sort),
	}
	err := system.UpdateDepartment(&department)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "更新部门失败")
		return
	}
	response.ReturnData(c, department)
}

func GetDepartmentList(c *gin.Context) {
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

	department := system.SystemDepartment{
		Name:   params.Name,
		Status: params.Status,
	}

	// 调用带分页的查询函数
	departments, total, err := system.FindDepartmentList(&department, page, pageSize)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询部门失败")
		return
	}

	// 返回带总数的结果
	response.ReturnDataWithTotal(c, int(total), departments)
}

func DeleteDepartment(c *gin.Context) {
	params := &struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	department := system.SystemDepartment{Model: gorm.Model{ID: params.ID}}
	err := system.GetDepartment(&department)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ReturnError(c, response.DATA_LOSS, "部门不存在")
			return
		}
		response.ReturnError(c, response.DATA_LOSS, "查询部门失败")
		return
	}
	if len(department.SystemUsers) > 0 {
		response.ReturnError(c, response.DATA_LOSS, "请先删除部门下的用户")
		return
	}
	err = system.DeleteDepartment(&department)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "删除部门失败")
		return
	}
	response.ReturnData(c, department)
}
