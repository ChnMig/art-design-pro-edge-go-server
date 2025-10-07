package system

import (
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/config"
	"api-server/db/pgdb"
)

// GetTenantByCode 根据企业编号获取租户信息
func GetTenantByCode(code string) (SystemTenant, error) {
	var tenant SystemTenant
	err := pgdb.GetClient().Where("code = ? AND status = 1", code).First(&tenant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tenant, nil // 返回空租户，ID为0表示未找到
		}
		zap.L().Error("failed to get tenant by code", zap.Error(err))
		return tenant, err
	}
	return tenant, nil
}

// GetTenant 获取租户信息
func GetTenant(tenant *SystemTenant) error {
	if err := pgdb.GetClient().Where(tenant).First(tenant).Error; err != nil {
		zap.L().Error("failed to get tenant", zap.Error(err))
		return err
	}
	return nil
}

// AddTenant 添加租户
func AddTenant(tenant *SystemTenant) error {
	if err := pgdb.GetClient().Create(tenant).Error; err != nil {
		zap.L().Error("failed to add tenant", zap.Error(err))
		return err
	}
	return nil
}

// UpdateTenant 更新租户
func UpdateTenant(tenant *SystemTenant) error {
	if err := pgdb.GetClient().Updates(tenant).Error; err != nil {
		zap.L().Error("failed to update tenant", zap.Error(err))
		return err
	}
	return nil
}

// DeleteTenant 删除租户
func DeleteTenant(tenant *SystemTenant) error {
	if err := pgdb.GetClient().Delete(tenant).Error; err != nil {
		zap.L().Error("failed to delete tenant", zap.Error(err))
		return err
	}
	return nil
}

// FindTenantList 查询租户列表，支持分页
func FindTenantList(tenant *SystemTenant, page, pageSize int) ([]SystemTenant, int64, error) {
	var tenants []SystemTenant
	var total int64
	db := pgdb.GetClient()

	// 构建基础查询
	baseQuery := db.Model(&SystemTenant{}).Where("deleted_at IS NULL")

	// 使用模糊查询
	if tenant.Code != "" {
		baseQuery = baseQuery.Where("code LIKE ?", "%"+tenant.Code+"%")
	}
	if tenant.Name != "" {
		baseQuery = baseQuery.Where("name LIKE ?", "%"+tenant.Name+"%")
	}
	if tenant.Status != 0 {
		baseQuery = baseQuery.Where("status = ?", tenant.Status)
	}

	// 获取符合条件的总记录数
	baseQuery.Count(&total)

	// 判断是否需要分页
	if page == config.CancelPage && pageSize == config.CancelPageSize {
		// 不分页，获取所有数据
		if err := baseQuery.Order("created_at DESC").Find(&tenants).Error; err != nil {
			zap.L().Error("failed to find all tenants", zap.Error(err))
			return nil, 0, err
		}
	} else {
		// 应用分页并获取数据
		if err := baseQuery.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&tenants).Error; err != nil {
			zap.L().Error("failed to find tenants", zap.Error(err))
			return nil, 0, err
		}
	}

	return tenants, total, nil
}

// FindAllTenants 查询所有租户
func FindAllTenants(tenants *[]SystemTenant) error {
	if err := pgdb.GetClient().Find(tenants).Error; err != nil {
		zap.L().Error("failed to find all tenants", zap.Error(err))
		return err
	}
	return nil
}

// ValidateTenant 验证租户状态和权限
func ValidateTenant(tenant *SystemTenant) error {
	if tenant.Status != StatusEnabled {
		return errors.New("tenant is disabled")
	}

	return nil
}

// SuggestTenantByCode 根据代码进行模糊查询，返回前N条启用中的租户
func SuggestTenantByCode(code string, limit int) ([]SystemTenant, error) {
	var tenants []SystemTenant
	if limit <= 0 {
		limit = 10
	}
	db := pgdb.GetClient()
	// 仅查询启用、未删除的租户，按创建时间倒序，模糊匹配code
	err := db.Model(&SystemTenant{}).
		Where("deleted_at IS NULL AND status = 1 AND code LIKE ?", "%"+code+"%").
		Order("created_at DESC").
		Limit(limit).
		Find(&tenants).Error
	if err != nil {
		zap.L().Error("failed to suggest tenant by code", zap.Error(err))
		return nil, err
	}
	return tenants, nil
}
