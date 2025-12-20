package user

import (
	"errors"

	"github.com/gin-gonic/gin"

	"api-server/api/middleware"
	"api-server/api/response"
	userdomain "api-server/domain/admin/user"
)

func FindUserByCache(c *gin.Context) {
	params := &struct {
		Username string `json:"username" form:"username"` // 昵称
		Name     string `json:"name" form:"name"`         // 姓名
		ID       uint   `json:"id" form:"id"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	// 如果提供了ID，则获取单个用户信息
	if params.ID > 0 {
		userInfo, err := userdomain.GetUserFromCache(params.ID)
		if err != nil {
			response.ReturnError(c, response.DATA_LOSS, "获取用户缓存数据失败")
			return
		}
		response.ReturnData(c, userInfo)
		return
	}

	// 获取分页参数
	page := middleware.GetPage(c)
	pageSize := middleware.GetPageSize(c)

	userList, total, err := userdomain.ListUsersFromCache(userdomain.CacheFilter{
		Username: params.Username,
		Name:     params.Name,
	}, page, pageSize)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "获取用户缓存列表失败")
		return
	}
	response.ReturnDataWithTotal(c, total, userList)
}

func FindUser(c *gin.Context) {
	params := &struct {
		Username     string `json:"username" form:"username"` // 昵称
		Name         string `json:"name" form:"name"`         // 姓名
		Phone        string `json:"phone" form:"phone"`
		DepartmentID uint   `json:"department_id" form:"department_id"`
		RoleID       uint   `json:"role_id" form:"role_id"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	// 获取当前租户ID
	tenantID := middleware.GetTenantID(c)
	if tenantID == 0 {
		response.ReturnError(c, response.UNAUTHENTICATED, "Invalid tenant context")
		return
	}

	page := middleware.GetPage(c)
	pageSize := middleware.GetPageSize(c)
	usersWithRelations, total, err := userdomain.FindUserList(tenantID, userdomain.FindUserQuery{
		Username:     params.Username,
		Name:         params.Name,
		Phone:        params.Phone,
		RoleID:       params.RoleID,
		DepartmentID: params.DepartmentID,
	}, page, pageSize)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, "查询用户失败")
		return
	}

	type userListItem struct {
		ID             uint   `json:"id"`
		TenantID       uint   `json:"tenant_id"`
		DepartmentID   uint   `json:"department_id"`
		RoleID         uint   `json:"role_id"`
		Name           string `json:"name"`
		Username       string `json:"username"`
		Account        string `json:"account"`
		Phone          string `json:"phone"`
		Gender         uint   `json:"gender"`
		Status         uint   `json:"status"`
		CreatedAt      int64  `json:"created_at"`
		UpdatedAt      int64  `json:"updated_at"`
		RoleName       string `json:"role_name"`
		RoleDesc       string `json:"role_desc"`
		DepartmentName string `json:"department_name"`
	}

	items := make([]userListItem, len(usersWithRelations))
	for i, item := range usersWithRelations {
		items[i] = userListItem{
			ID:             item.SystemUser.ID,
			TenantID:       item.SystemUser.TenantID,
			DepartmentID:   item.SystemUser.DepartmentID,
			RoleID:         item.SystemUser.RoleID,
			Name:           item.SystemUser.Name,
			Username:       item.SystemUser.Username,
			Account:        item.SystemUser.Account,
			Phone:          item.SystemUser.Phone,
			Gender:         item.SystemUser.Gender,
			Status:         item.SystemUser.Status,
			CreatedAt:      item.SystemUser.CreatedAt.Unix(),
			UpdatedAt:      item.SystemUser.UpdatedAt.Unix(),
			RoleName:       item.RoleName,
			RoleDesc:       item.RoleDesc,
			DepartmentName: item.DepartmentName,
		}
	}

	response.ReturnDataWithTotal(c, int(total), items)
}

func AddUser(c *gin.Context) {
	params := &struct {
		Name         string `json:"name" form:"name" binding:"required"`         // 姓名
		Username     string `json:"username" form:"username" binding:"required"` // 昵称
		Account      string `json:"account" form:"account" binding:"required"`   // 登录账号
		Password     string `json:"password" form:"password" binding:"required"`
		Phone        string `json:"phone" form:"phone" binding:"required"`
		Gender       uint   `json:"gender" form:"gender" binding:"required"`
		Status       uint   `json:"status" form:"status" binding:"required"`
		RoleID       uint   `json:"role_id" form:"role_id" binding:"required"`
		DepartmentID uint   `json:"department_id" form:"department_id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	// 获取当前租户ID
	tenantID := middleware.GetTenantID(c)
	if tenantID == 0 {
		response.ReturnError(c, response.UNAUTHENTICATED, "Invalid tenant context")
		return
	}
	if err := userdomain.AddUser(tenantID, userdomain.AddUserInput{
		Name:         params.Name,
		Username:     params.Username,
		Account:      params.Account,
		Password:     params.Password,
		Phone:        params.Phone,
		Gender:       params.Gender,
		Status:       params.Status,
		RoleID:       params.RoleID,
		DepartmentID: params.DepartmentID,
	}); err != nil {
		if errors.Is(err, userdomain.ErrRoleNotInTenant) {
			response.ReturnError(c, response.PERMISSION_DENIED, "角色不存在或不属于当前租户")
			return
		}
		response.ReturnError(c, response.DATA_LOSS, "添加用户失败")
		return
	}
	response.ReturnData(c, nil)
}

func UpdateUser(c *gin.Context) {
	params := &struct {
		ID           uint   `json:"id" form:"id" binding:"required"`
		Name         string `json:"name" form:"name" binding:"required"`         // 姓名
		Username     string `json:"username" form:"username" binding:"required"` // 昵称
		Account      string `json:"account" form:"account" binding:"required"`   // 登录账号
		Password     string `json:"password" form:"password"`
		Phone        string `json:"phone" form:"phone" binding:"required"`
		Gender       uint   `json:"gender" form:"gender" binding:"required"`
		Status       uint   `json:"status" form:"status" binding:"required"`
		RoleID       uint   `json:"role_id" form:"role_id" binding:"required"`
		DepartmentID uint   `json:"department_id" form:"department_id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	// 获取当前租户ID
	tenantID := middleware.GetTenantID(c)
	if tenantID == 0 {
		response.ReturnError(c, response.UNAUTHENTICATED, "Invalid tenant context")
		return
	}
	if err := userdomain.UpdateUser(tenantID, userdomain.UpdateUserInput{
		ID:           params.ID,
		Name:         params.Name,
		Username:     params.Username,
		Account:      params.Account,
		Password:     params.Password,
		Phone:        params.Phone,
		Gender:       params.Gender,
		Status:       params.Status,
		RoleID:       params.RoleID,
		DepartmentID: params.DepartmentID,
	}); err != nil {
		if errors.Is(err, userdomain.ErrRoleNotInTenant) {
			response.ReturnError(c, response.PERMISSION_DENIED, "角色不存在或不属于当前租户")
			return
		}
		response.ReturnError(c, response.DATA_LOSS, "更新用户失败")
		return
	}
	response.ReturnData(c, nil)
}

func DeleteUser(c *gin.Context) {
	params := &struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}
	if err := userdomain.DeleteUser(params.ID); err != nil {
		if errors.Is(err, userdomain.ErrCannotDeleteSuperAdmin) {
			response.ReturnError(c, response.DATA_LOSS, "不能删除超级管理员")
			return
		}
		response.ReturnError(c, response.DATA_LOSS, "删除用户失败")
		return
	}
	response.ReturnData(c, nil)
}
