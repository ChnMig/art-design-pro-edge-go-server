package user

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"api-server/api/middleware"
	"api-server/api/response"
	"api-server/db/pgdb/system"
	"api-server/db/rdb/systemuser"
)

func FindUserByCache(c *gin.Context) {
	params := &struct {
		Username string `json:"username" form:"username"`
		Name     string `json:"name" form:"name"`
		ID       uint   `json:"id" form:"id"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	// 如果提供了ID，则获取单个用户信息
	if params.ID > 0 {
		userInfo, err := systemuser.GetUserFromCache(params.ID)
		if err != nil {
			response.ReturnError(c, response.DATA_LOSS, "获取用户缓存数据失败")
			return
		}
		response.ReturnOk(c, userInfo)
		return
	}

	// 获取分页参数
	page := middleware.GetPage(c)
	pageSize := middleware.GetPageSize(c)

	// 获取所有用户列表
	userList, err := systemuser.GetAllUsersFromCache()
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "获取用户缓存列表失败")
		return
	}

	// 过滤结果
	var filteredList []systemuser.UserCacheInfo
	if params.Username != "" || params.Name != "" {
		for _, user := range userList {
			// 如果提供了用户名，且不匹配，则跳过
			if params.Username != "" && !strings.Contains(user.Username, params.Username) {
				continue
			}

			// 如果提供了名称，且不匹配，则跳过
			if params.Name != "" && !strings.Contains(user.Name, params.Name) {
				continue
			}

			// 所有条件都匹配，添加到结果列表
			filteredList = append(filteredList, user)
		}
	} else {
		filteredList = userList
	}

	// 计算总数
	total := len(filteredList)

	// 应用分页
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= total {
		// 如果起始位置超出了总数，返回空列表
		response.ReturnOkWithCount(c, total, []systemuser.UserCacheInfo{})
		return
	}
	if end > total {
		end = total
	}

	pagedList := filteredList[start:end]
	response.ReturnOkWithCount(c, total, pagedList)
}

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

	// 使用索引方式清空密码
	for i := range usersWithRelations {
		usersWithRelations[i].SystemUser.Password = ""
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
