package system

import (
	"github.com/gin-gonic/gin"

	"api-server/internal/service/system"
	"api-server/internal/transport/http/middleware"
	"api-server/internal/transport/http/response"
)

// UserHandler 用户HTTP处理器
type UserHandler struct {
	userService system.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService system.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Login 用户登录处理器
func (h *UserHandler) Login(c *gin.Context) {
	// 1. 参数绑定和验证
	req := &system.LoginRequest{}
	if !middleware.CheckParam(req, c) {
		return
	}

	// 2. 设置客户端IP
	req.IP = c.ClientIP()

	// 3. 调用认证服务
	authService := system.NewAuthService()
	result, err := authService.Login(c.Request.Context(), req)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, err.Error())
		return
	}

	// 4. 返回成功结果
	response.ReturnData(c, result)
}

// GetUserInfo 获取用户信息处理器
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	// 1. 获取当前用户ID
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未认证")
		return
	}

	// 2. 调用用户服务
	result, err := h.userService.GetUserInfo(c.Request.Context(), userID)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, err.Error())
		return
	}

	// 3. 返回结果
	response.ReturnData(c, result.User)
}

// UpdateUserInfo 更新用户信息处理器
func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	// 1. 参数绑定和验证
	req := &system.UpdateUserInfoRequest{}
	if !middleware.CheckParam(req, c) {
		return
	}

	// 2. 获取当前用户ID
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		response.ReturnError(c, response.UNAUTHENTICATED, "用户未认证")
		return
	}

	// 3. 调用用户服务
	err := h.userService.UpdateUserInfo(c.Request.Context(), userID, req)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, err.Error())
		return
	}

	// 4. 返回成功
	response.ReturnData(c, nil)
}

// FindUser 查询用户列表处理器
func (h *UserHandler) FindUser(c *gin.Context) {
	// 1. 参数绑定
	req := &system.FindUserListRequest{}
	if !middleware.CheckParam(req, c) {
		return
	}

	// 2. 获取分页参数
	req.Page = middleware.GetPage(c)
	req.PageSize = middleware.GetPageSize(c)

	// 3. 获取当前用户的租户ID（多租户隔离）
	tenantID := middleware.GetCurrentTenantID(c)
	req.TenantID = tenantID

	// 4. 调用用户服务
	result, err := h.userService.FindUserList(c.Request.Context(), req)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, err.Error())
		return
	}

	// 5. 返回结果
	response.ReturnDataWithTotal(c, int(result.Total), result.Users)
}

// AddUser 添加用户处理器
func (h *UserHandler) AddUser(c *gin.Context) {
	// 1. 参数绑定和验证
	req := &system.CreateUserRequest{}
	if !middleware.CheckParam(req, c) {
		return
	}

	// 2. 设置租户ID（多租户隔离）
	req.TenantID = middleware.GetCurrentTenantID(c)

	// 3. 调用用户服务
	err := h.userService.CreateUser(c.Request.Context(), req)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, err.Error())
		return
	}

	// 4. 返回成功
	response.ReturnData(c, nil)
}

// UpdateUser 更新用户处理器
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// 1. 参数绑定和验证
	req := &system.UpdateUserRequest{}
	if !middleware.CheckParam(req, c) {
		return
	}

	// 2. 调用用户服务
	err := h.userService.UpdateUser(c.Request.Context(), req)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, err.Error())
		return
	}

	// 3. 返回成功
	response.ReturnData(c, nil)
}

// DeleteUser 删除用户处理器
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// 1. 参数绑定
	params := &struct {
		ID uint `json:"id" form:"id" binding:"required"`
	}{}
	if !middleware.CheckParam(params, c) {
		return
	}

	// 2. 调用用户服务
	err := h.userService.DeleteUser(c.Request.Context(), params.ID)
	if err != nil {
		response.ReturnError(c, response.DATA_LOSS, err.Error())
		return
	}

	// 3. 返回成功
	response.ReturnData(c, nil)
}