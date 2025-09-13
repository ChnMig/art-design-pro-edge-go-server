package system

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"api-server/internal/domain/system"
)

// UserService 用户服务接口
type UserService interface {
	// GetUserInfo 获取用户信息
	GetUserInfo(ctx context.Context, userID uint) (*UserInfoResponse, error)
	// UpdateUserInfo 更新用户信息
	UpdateUserInfo(ctx context.Context, userID uint, req *UpdateUserInfoRequest) error
	// FindUserList 查询用户列表
	FindUserList(ctx context.Context, req *FindUserListRequest) (*FindUserListResponse, error)
	// CreateUser 创建用户
	CreateUser(ctx context.Context, req *CreateUserRequest) error
	// UpdateUser 更新用户
	UpdateUser(ctx context.Context, req *UpdateUserRequest) error
	// DeleteUser 删除用户
	DeleteUser(ctx context.Context, userID uint) error
}

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	User *system.SystemUser `json:"user"`
}

// UpdateUserInfoRequest 更新用户信息请求
type UpdateUserInfoRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// FindUserListRequest 查询用户列表请求
type FindUserListRequest struct {
	TenantID     uint   `json:"tenant_id"`
	Username     string `json:"username"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	RoleID       uint   `json:"role_id"`
	DepartmentID uint   `json:"department_id"`
	Page         int    `json:"page"`
	PageSize     int    `json:"page_size"`
}

// FindUserListResponse 查询用户列表响应
type FindUserListResponse struct {
	Users []system.UserWithRelations `json:"users"`
	Total int64                      `json:"total"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	TenantID     uint   `json:"tenant_id" binding:"required"`
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Phone        string `json:"phone"`
	RoleID       uint   `json:"role_id" binding:"required"`
	DepartmentID uint   `json:"department_id"`
	Status       uint   `json:"status"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	ID           uint   `json:"id" binding:"required"`
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password"`
	Name         string `json:"name" binding:"required"`
	Phone        string `json:"phone"`
	RoleID       uint   `json:"role_id" binding:"required"`
	DepartmentID uint   `json:"department_id"`
	Status       uint   `json:"status"`
}

// userService 用户服务实现
type userService struct {
	logger *zap.Logger
}

// NewUserService 创建用户服务
func NewUserService() UserService {
	return &userService{
		logger: zap.L(),
	}
}

// GetUserInfo 获取用户信息
func (s *userService) GetUserInfo(ctx context.Context, userID uint) (*UserInfoResponse, error) {
	user := system.SystemUser{}
	user.ID = userID

	if err := system.GetUser(&user); err != nil {
		s.logger.Error("获取用户信息失败", zap.Uint("user_id", userID), zap.Error(err))
		return nil, errors.New("获取用户信息失败")
	}

	// 清空敏感信息
	user.Password = ""

	return &UserInfoResponse{
		User: &user,
	}, nil
}

// UpdateUserInfo 更新用户信息
func (s *userService) UpdateUserInfo(ctx context.Context, userID uint, req *UpdateUserInfoRequest) error {
	user := system.SystemUser{
		Name:  req.Name,
		Phone: req.Phone,
	}
	user.ID = userID

	if err := system.UpdateUser(&user); err != nil {
		s.logger.Error("更新用户信息失败", zap.Uint("user_id", userID), zap.Error(err))
		return errors.New("更新用户信息失败")
	}

	s.logger.Info("用户信息更新成功", zap.Uint("user_id", userID))
	return nil
}

// FindUserList 查询用户列表
func (s *userService) FindUserList(ctx context.Context, req *FindUserListRequest) (*FindUserListResponse, error) {
	user := system.SystemUser{
		TenantID:     req.TenantID,
		Username:     req.Username,
		Name:         req.Name,
		Phone:        req.Phone,
		RoleID:       req.RoleID,
		DepartmentID: req.DepartmentID,
	}

	users, total, err := system.FindUserList(&user, req.Page, req.PageSize)
	if err != nil {
		s.logger.Error("查询用户列表失败", zap.Error(err))
		return nil, errors.New("查询用户列表失败")
	}

	return &FindUserListResponse{
		Users: users,
		Total: total,
	}, nil
}

// CreateUser 创建用户
func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) error {
	// TODO: 添加用户名唯一性检查

	user := system.SystemUser{
		TenantID:     req.TenantID,
		Username:     req.Username,
		Password:     req.Password,
		Name:         req.Name,
		Phone:        req.Phone,
		RoleID:       req.RoleID,
		DepartmentID: req.DepartmentID,
		Status:       req.Status,
	}

	if user.Status == 0 {
		user.Status = 1 // 默认状态为启用
	}

	if err := system.AddUser(&user); err != nil {
		s.logger.Error("创建用户失败", zap.Error(err))
		return errors.New("创建用户失败")
	}

	s.logger.Info("用户创建成功",
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username),
	)

	return nil
}

// UpdateUser 更新用户
func (s *userService) UpdateUser(ctx context.Context, req *UpdateUserRequest) error {
	user := system.SystemUser{
		Username:     req.Username,
		Password:     req.Password,
		Name:         req.Name,
		Phone:        req.Phone,
		RoleID:       req.RoleID,
		DepartmentID: req.DepartmentID,
		Status:       req.Status,
	}
	user.ID = req.ID

	if err := system.UpdateUser(&user); err != nil {
		s.logger.Error("更新用户失败", zap.Uint("user_id", req.ID), zap.Error(err))
		return errors.New("更新用户失败")
	}

	s.logger.Info("用户更新成功", zap.Uint("user_id", req.ID))
	return nil
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(ctx context.Context, userID uint) error {
	user := system.SystemUser{}
	user.ID = userID

	if err := system.DeleteUser(&user); err != nil {
		s.logger.Error("删除用户失败", zap.Uint("user_id", userID), zap.Error(err))
		return errors.New("删除用户失败")
	}

	s.logger.Info("用户删除成功", zap.Uint("user_id", userID))
	return nil
}