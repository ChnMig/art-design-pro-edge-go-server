package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"api-server/api/middleware"
	"api-server/api/response"
	"api-server/db/pgdb/system"
)

func FindUser(c *gin.Context) {
	params := &struct {
		Username     string `json:"username" form:"username"`
		Name         string `json:"name" form:"name"`
		Phone        string `json:"phone" form:"phone"`
		DepartmentID uint   `json:"department_id" form:"department_id"`
		RoleID       uint   `json:"role_id" form:"role_id"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	page := middleware.GetPage(c)
	pageSize := middleware.GetPageSize(c)
	u := system.SystemUser{
		Username:     params.Username,
		Name:         params.Name,
		Phone:        params.Phone,
		RoleID:       params.RoleID,
		DepartmentID: params.DepartmentID,
	}
	usersWithRelations, total, err := system.FindUserList(&u, page, pageSize)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询用户失败")
		return
	}
	for _, v := range usersWithRelations {
		v.Password = ""
	}
	response.ReturnOkWithCount(c, int(total), usersWithRelations)
}

func AddUser(c *gin.Context) {
	params := &struct {
		Username     string `json:"username" form:"username" binding:"required"`
		Password     string `json:"password" form:"password" binding:"required"`
		Name         string `json:"name" form:"name" binding:"required"`
		Phone        string `json:"phone" form:"phone" binding:"required"`
		Gender       uint   `json:"gender" form:"gender" binding:"required"`
		Status       uint   `json:"status" form:"status" binding:"required"`
		RoleID       uint   `json:"role_id" form:"role_id" binding:"required"`
		DepartmentID uint   `json:"department_id" form:"department_id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	u := system.SystemUser{
		Username:     params.Username,
		Password:     params.Password,
		Name:         params.Name,
		Phone:        params.Phone,
		Gender:       params.Gender,
		Status:       params.Status,
		RoleID:       params.RoleID,
		DepartmentID: params.DepartmentID,
	}
	if err := system.AddUser(&u); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "添加用户失败")
		return
	}
	response.ReturnOk(c, nil)
}

func UpdateUser(c *gin.Context) {
	params := &struct {
		ID           uint   `json:"id" form:"id" binding:"required"`
		Username     string `json:"username" form:"username" binding:"required"`
		Password     string `json:"password" form:"password"`
		Name         string `json:"name" form:"name" binding:"required"`
		Phone        string `json:"phone" form:"phone" binding:"required"`
		Gender       uint   `json:"gender" form:"gender" binding:"required"`
		Status       uint   `json:"status" form:"status" binding:"required"`
		RoleID       uint   `json:"role_id" form:"role_id" binding:"required"`
		DepartmentID uint   `json:"department_id" form:"department_id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	u := system.SystemUser{
		Model:        gorm.Model{ID: params.ID},
		Username:     params.Username,
		Name:         params.Name,
		Phone:        params.Phone,
		Gender:       params.Gender,
		Status:       params.Status,
		RoleID:       params.RoleID,
		DepartmentID: params.DepartmentID,
	}
	if params.Password != "" {
		u.Password = params.Password
	}
	if err := system.UpdateUser(&u); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "更新用户失败")
		return
	}
	response.ReturnOk(c, nil)
}

func DeleteUser(c *gin.Context) {
	params := &struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	if params.ID == 1 {
		response.ReturnError(c, response.DATA_LOSS, "不能删除超级管理员")
		return
	}
	u := system.SystemUser{
		Model: gorm.Model{ID: params.ID},
	}
	if err := system.DeleteUser(&u); err != nil {
		response.ReturnError(c, response.DATA_LOSS, "删除用户失败")
		return
	}
	response.ReturnOk(c, nil)
}
