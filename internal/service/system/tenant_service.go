package system

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"api-server/internal/domain/system"
)

// TenantService 租户服务接口
type TenantService interface {
	// FindTenantList 查询租户列表
	FindTenantList(ctx context.Context, req *FindTenantListRequest) (*FindTenantListResponse, error)
	// CreateTenant 创建租户
	CreateTenant(ctx context.Context, req *CreateTenantRequest) error
	// UpdateTenant 更新租户
	UpdateTenant(ctx context.Context, req *UpdateTenantRequest) error
	// DeleteTenant 删除租户
	DeleteTenant(ctx context.Context, tenantID uint) error
	// GetTenantByCode 根据编号获取租户
	GetTenantByCode(ctx context.Context, code string) (*system.SystemTenant, error)
}

// FindTenantListRequest 查询租户列表请求
type FindTenantListRequest struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status uint   `json:"status"`
	Page   int    `json:"page"`
	PageSize int  `json:"page_size"`
}

// FindTenantListResponse 查询租户列表响应
type FindTenantListResponse struct {
	Tenants []system.SystemTenant `json:"tenants"`
	Total   int64                 `json:"total"`
}

// CreateTenantRequest 创建租户请求
type CreateTenantRequest struct {
	Code      string     `json:"code" binding:"required"`
	Name      string     `json:"name" binding:"required"`
	Contact   string     `json:"contact"`
	Phone     string     `json:"phone"`
	Email     string     `json:"email"`
	Address   string     `json:"address"`
	Status    uint       `json:"status" binding:"required"`
	ExpiredAt *time.Time `json:"expired_at"`
	MaxUsers  uint       `json:"max_users"`
}

// UpdateTenantRequest 更新租户请求
type UpdateTenantRequest struct {
	ID        uint       `json:"id" binding:"required"`
	Code      string     `json:"code" binding:"required"`
	Name      string     `json:"name" binding:"required"`
	Contact   string     `json:"contact"`
	Phone     string     `json:"phone"`
	Email     string     `json:"email"`
	Address   string     `json:"address"`
	Status    uint       `json:"status" binding:"required"`
	ExpiredAt *time.Time `json:"expired_at"`
	MaxUsers  uint       `json:"max_users"`
}

// tenantService 租户服务实现
type tenantService struct {
	logger *zap.Logger
}

// NewTenantService 创建租户服务
func NewTenantService() TenantService {
	return &tenantService{
		logger: zap.L(),
	}
}

// FindTenantList 查询租户列表
func (s *tenantService) FindTenantList(ctx context.Context, req *FindTenantListRequest) (*FindTenantListResponse, error) {
	tenant := system.SystemTenant{
		Code:   req.Code,
		Name:   req.Name,
		Status: req.Status,
	}

	tenants, total, err := system.FindTenantList(&tenant, req.Page, req.PageSize)
	if err != nil {
		s.logger.Error("查询租户列表失败", zap.Error(err))
		return nil, errors.New("查询租户列表失败")
	}

	return &FindTenantListResponse{
		Tenants: tenants,
		Total:   total,
	}, nil
}

// CreateTenant 创建租户
func (s *tenantService) CreateTenant(ctx context.Context, req *CreateTenantRequest) error {
	// 检查租户编号是否已存在
	existingTenant, err := s.GetTenantByCode(ctx, req.Code)
	if err == nil && existingTenant.ID > 0 {
		return errors.New("租户编号已存在")
	}

	tenant := system.SystemTenant{
		Code:      req.Code,
		Name:      req.Name,
		Contact:   req.Contact,
		Phone:     req.Phone,
		Email:     req.Email,
		Address:   req.Address,
		Status:    req.Status,
		ExpiredAt: req.ExpiredAt,
		MaxUsers:  req.MaxUsers,
	}

	if tenant.MaxUsers == 0 {
		tenant.MaxUsers = 100 // 默认最大用户数
	}

	if err := system.AddTenant(&tenant); err != nil {
		s.logger.Error("创建租户失败", zap.Error(err))
		return errors.New("创建租户失败")
	}

	s.logger.Info("租户创建成功",
		zap.Uint("tenant_id", tenant.ID),
		zap.String("tenant_code", tenant.Code),
		zap.String("tenant_name", tenant.Name),
	)

	return nil
}

// UpdateTenant 更新租户
func (s *tenantService) UpdateTenant(ctx context.Context, req *UpdateTenantRequest) error {
	// 检查租户是否存在
	existingTenant := system.SystemTenant{}
	existingTenant.ID = req.ID
	if err := system.GetTenant(&existingTenant); err != nil {
		return errors.New("租户不存在")
	}

	// 如果修改了编号，检查新编号是否已被其他租户使用
	if req.Code != existingTenant.Code {
		conflictTenant, err := s.GetTenantByCode(ctx, req.Code)
		if err == nil && conflictTenant.ID > 0 && conflictTenant.ID != req.ID {
			return errors.New("租户编号已被使用")
		}
	}

	tenant := system.SystemTenant{
		Code:      req.Code,
		Name:      req.Name,
		Contact:   req.Contact,
		Phone:     req.Phone,
		Email:     req.Email,
		Address:   req.Address,
		Status:    req.Status,
		ExpiredAt: req.ExpiredAt,
		MaxUsers:  req.MaxUsers,
	}
	tenant.ID = req.ID

	if err := system.UpdateTenant(&tenant); err != nil {
		s.logger.Error("更新租户失败", zap.Uint("tenant_id", req.ID), zap.Error(err))
		return errors.New("更新租户失败")
	}

	s.logger.Info("租户更新成功", zap.Uint("tenant_id", req.ID))
	return nil
}

// DeleteTenant 删除租户
func (s *tenantService) DeleteTenant(ctx context.Context, tenantID uint) error {
	// 检查租户下是否还有用户
	users := []system.SystemUser{}
	if err := system.FindAllUsers(&users); err != nil {
		s.logger.Error("查询用户列表失败", zap.Error(err))
		return errors.New("无法验证租户状态")
	}

	// 检查该租户下是否还有用户
	for _, user := range users {
		if user.TenantID == tenantID {
			return errors.New("租户下还有用户，无法删除")
		}
	}

	tenant := system.SystemTenant{}
	tenant.ID = tenantID

	if err := system.DeleteTenant(&tenant); err != nil {
		s.logger.Error("删除租户失败", zap.Uint("tenant_id", tenantID), zap.Error(err))
		return errors.New("删除租户失败")
	}

	s.logger.Info("租户删除成功", zap.Uint("tenant_id", tenantID))
	return nil
}

// GetTenantByCode 根据编号获取租户
func (s *tenantService) GetTenantByCode(ctx context.Context, code string) (*system.SystemTenant, error) {
	tenant, err := system.GetTenantByCode(code)
	if err != nil {
		s.logger.Error("根据编号获取租户失败", zap.String("code", code), zap.Error(err))
		return nil, err
	}

	return &tenant, nil
}