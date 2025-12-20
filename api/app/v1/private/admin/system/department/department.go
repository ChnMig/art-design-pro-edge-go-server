package department

import (
	"github.com/gin-gonic/gin"

	"api-server/api/middleware"
	"api-server/api/response"
	departmentdomain "api-server/domain/admin/department"
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
	department, err := departmentdomain.AddDepartment(departmentdomain.AddInput{
		Name:   params.Name,
		Status: uint(params.Status),
		Sort:   uint(params.Sort),
	})
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "添加部门失败")
		return
	}
	response.ReturnData(c, department)
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
	department, err := departmentdomain.UpdateDepartment(departmentdomain.UpdateInput{
		ID:     params.ID,
		Name:   params.Name,
		Status: uint(params.Status),
		Sort:   uint(params.Sort),
	})
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

	// 调用带分页的查询函数
	departments, total, err := departmentdomain.FindDepartmentList(departmentdomain.FindListQuery{
		Name:   params.Name,
		Status: params.Status,
	}, page, pageSize)
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
	department, err := departmentdomain.DeleteDepartment(params.ID)
	if err != nil {
		ReturnDomainError(c, err, "删除部门失败")
		return
	}
	response.ReturnData(c, department)
}
