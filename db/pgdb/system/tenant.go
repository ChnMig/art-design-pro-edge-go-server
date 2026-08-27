package system

import (
	"errors"
	"fmt"
	"strings"

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

// TenantBootstrap 描述新租户首次开户所需的完整数据。
type TenantBootstrap struct {
	Tenant        *SystemTenant
	AdminName     string
	AdminAccount  string
	AdminPassword string
}

// CreateTenantWithAdmin 在同一事务中创建租户、默认部门、管理员角色、首个管理员和权限范围。
func CreateTenantWithAdmin(bootstrap TenantBootstrap) error {
	if bootstrap.Tenant == nil || bootstrap.Tenant.Code == "" || bootstrap.Tenant.Name == "" ||
		bootstrap.AdminAccount == "" || bootstrap.AdminPassword == "" {
		return gorm.ErrInvalidData
	}
	hashedPassword, err := HashPassword(bootstrap.AdminPassword)
	if err != nil {
		return err
	}
	return pgdb.GetClient().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(bootstrap.Tenant).Error; err != nil {
			return err
		}

		department := SystemDepartment{
			TenantID: bootstrap.Tenant.ID,
			Name:     "管理中心",
			Sort:     1,
			Status:   StatusEnabled,
		}
		if err := tx.Create(&department).Error; err != nil {
			return err
		}

		adminRole := SystemRole{
			TenantID:      bootstrap.Tenant.ID,
			Name:          "租户管理员",
			Desc:          "负责当前租户的用户、角色和部门管理",
			Status:        StatusEnabled,
			IsTenantAdmin: true,
		}
		if err := tx.Create(&adminRole).Error; err != nil {
			return err
		}

		menus, auths, err := loadTenantAssignablePermissions(tx)
		if err != nil {
			return err
		}
		if len(menus) > 0 {
			menuScopes := make([]SystemTenantMenuScope, 0, len(menus))
			for _, menu := range menus {
				menuScopes = append(menuScopes, SystemTenantMenuScope{TenantID: bootstrap.Tenant.ID, MenuID: menu.ID})
			}
			if err := tx.Create(&menuScopes).Error; err != nil {
				return err
			}
			if err := tx.Model(&adminRole).Association("SystemMenus").Append(&menus); err != nil {
				return err
			}
		}
		if len(auths) > 0 {
			authScopes := make([]SystemTenantAuthScope, 0, len(auths))
			for _, auth := range auths {
				authScopes = append(authScopes, SystemTenantAuthScope{TenantID: bootstrap.Tenant.ID, AuthID: auth.ID})
			}
			if err := tx.Create(&authScopes).Error; err != nil {
				return err
			}
			if err := tx.Model(&adminRole).Association("SystemMenuAuths").Append(&auths); err != nil {
				return err
			}
		}

		adminName := bootstrap.AdminName
		if adminName == "" {
			adminName = "租户管理员"
		}
		admin := SystemUser{
			TenantID:     bootstrap.Tenant.ID,
			DepartmentID: department.ID,
			RoleID:       adminRole.ID,
			Name:         adminName,
			Username:     bootstrap.AdminAccount,
			Account:      bootstrap.AdminAccount,
			Password:     hashedPassword,
			Status:       StatusEnabled,
			Gender:       1,
		}
		return tx.Create(&admin).Error
	})
}

func loadTenantAssignablePermissions(tx *gorm.DB) ([]SystemMenu, []SystemMenuAuth, error) {
	var menus []SystemMenu
	if err := tx.Where("status = ?", StatusEnabled).Order("id ASC").Find(&menus).Error; err != nil {
		return nil, nil, err
	}
	excluded := make(map[uint]struct{})
	for _, menu := range menus {
		if menu.ParentID == 0 && (menu.Path == "/platform" || menu.Name == "Platform") {
			excluded[menu.ID] = struct{}{}
		}
	}
	for changed := true; changed; {
		changed = false
		for _, menu := range menus {
			if _, found := excluded[menu.ID]; found {
				continue
			}
			if _, parentExcluded := excluded[menu.ParentID]; parentExcluded {
				excluded[menu.ID] = struct{}{}
				changed = true
			}
		}
	}
	allowedMenus := make([]SystemMenu, 0, len(menus))
	allowedMenuIDs := make([]uint, 0, len(menus))
	for _, menu := range menus {
		if _, found := excluded[menu.ID]; !found {
			allowedMenus = append(allowedMenus, menu)
			allowedMenuIDs = append(allowedMenuIDs, menu.ID)
		}
	}
	if len(allowedMenuIDs) == 0 {
		return allowedMenus, []SystemMenuAuth{}, nil
	}
	var auths []SystemMenuAuth
	if err := tx.Where("menu_id IN ?", allowedMenuIDs).Order("id ASC").Find(&auths).Error; err != nil {
		return nil, nil, fmt.Errorf("查询租户可分配按钮权限失败: %w", err)
	}
	return allowedMenus, auths, nil
}

// UpdateTenant 更新租户
func UpdateTenant(tenant *SystemTenant) error {
	if tenant.ID == 0 {
		return gorm.ErrInvalidData
	}
	result := pgdb.GetClient().Model(&SystemTenant{}).
		Where("id = ?", tenant.ID).
		Select("code", "name", "contact", "phone", "email", "status").
		Updates(tenant)
	if result.Error != nil {
		err := result.Error
		zap.L().Error("failed to update tenant", zap.Error(err))
		return err
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeleteTenant 删除租户
func DeleteTenant(tenant *SystemTenant) error {
	if tenant.ID == 0 {
		return gorm.ErrInvalidData
	}
	return pgdb.GetClient().Transaction(func(tx *gorm.DB) error {
		var roleIDs []uint
		if err := tx.Model(&SystemRole{}).Where("tenant_id = ?", tenant.ID).Pluck("id", &roleIDs).Error; err != nil {
			return err
		}
		if len(roleIDs) > 0 {
			if err := tx.Exec("DELETE FROM system_roles__system_menus WHERE system_role_id IN ?", roleIDs).Error; err != nil {
				return err
			}
			if err := tx.Exec("DELETE FROM system_roles__system_auths WHERE system_role_id IN ?", roleIDs).Error; err != nil {
				return err
			}
		}
		if err := tx.Unscoped().Where("tenant_id = ?", tenant.ID).Delete(&SystemTenantMenuScope{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("tenant_id = ?", tenant.ID).Delete(&SystemTenantAuthScope{}).Error; err != nil {
			return err
		}
		if err := tx.Where("tenant_id = ?", tenant.ID).Delete(&SystemUser{}).Error; err != nil {
			return err
		}
		if err := tx.Where("tenant_id = ?", tenant.ID).Delete(&SystemRole{}).Error; err != nil {
			return err
		}
		if err := tx.Where("tenant_id = ?", tenant.ID).Delete(&SystemDepartment{}).Error; err != nil {
			return err
		}
		result := tx.Where("id = ?", tenant.ID).Delete(&SystemTenant{})
		if result.Error != nil {
			zap.L().Error("failed to delete tenant", zap.Error(result.Error))
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
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
	if err := baseQuery.Count(&total).Error; err != nil {
		zap.L().Error("failed to count tenants", zap.Error(err))
		return nil, 0, err
	}

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
	escapedCode := strings.NewReplacer("!", "!!", "%", "!%", "_", "!_").Replace(code)
	// 仅查询启用、未删除的租户，并把用户输入按字面量执行前缀匹配。
	err := db.Model(&SystemTenant{}).
		Where("deleted_at IS NULL AND status = 1 AND code LIKE ? ESCAPE '!'", escapedCode+"%").
		Order("created_at DESC").
		Limit(limit).
		Find(&tenants).Error
	if err != nil {
		zap.L().Error("failed to suggest tenant by code", zap.Error(err))
		return nil, err
	}
	return tenants, nil
}
