package system

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/config"
	"api-server/db/pgdb"
)

func UpdateRole(role *SystemRole) error {
	if role.ID == 0 || role.TenantID == 0 {
		return gorm.ErrInvalidData
	}
	result := pgdb.GetClient().Model(&SystemRole{}).
		Where("id = ? AND tenant_id = ?", role.ID, role.TenantID).
		Select("name", "desc", "status", "is_tenant_admin").
		Updates(role)
	if result.Error != nil {
		err := result.Error
		zap.L().Error("failed to update role", zap.Error(err))
		return err
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func AddRole(role *SystemRole) error {
	if err := pgdb.GetClient().Create(&role).Error; err != nil {
		zap.L().Error("failed to create role", zap.Error(err))
		return err
	}
	return nil
}

func DeleteRole(role *SystemRole) error {
	if role.ID == 0 || role.TenantID == 0 {
		return gorm.ErrInvalidData
	}
	return pgdb.GetClient().Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM system_roles__system_menus WHERE system_role_id = ?", role.ID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM system_roles__system_auths WHERE system_role_id = ?", role.ID).Error; err != nil {
			return err
		}
		result := tx.Where("id = ? AND tenant_id = ?", role.ID, role.TenantID).Delete(&SystemRole{})
		if result.Error != nil {
			zap.L().Error("failed to delete role", zap.Error(result.Error))
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

// GetRole 获取单个角色信息
func GetRole(role *SystemRole) error {
	if err := pgdb.GetClient().Where(role).First(role).Error; err != nil {
		zap.L().Error("failed to get role", zap.Error(err))
		return err
	}
	return nil
}

// CountUsersByRoleID 统计租户内仍在使用指定角色的用户数量。
func CountUsersByRoleID(tenantID, roleID uint) (int64, error) {
	if tenantID == 0 || roleID == 0 {
		return 0, gorm.ErrInvalidData
	}
	var count int64
	if err := pgdb.GetClient().Model(&SystemUser{}).
		Where("tenant_id = ? AND role_id = ?", tenantID, roleID).
		Count(&count).Error; err != nil {
		zap.L().Error("failed to count users by role", zap.Error(err))
		return 0, err
	}
	return count, nil
}

// FindAllRoles 查询所有角色
func FindAllRoles(roles *[]SystemRole) error {
	if err := pgdb.GetClient().Find(roles).Error; err != nil {
		zap.L().Error("failed to find all roles", zap.Error(err))
		return err
	}
	return nil
}

// FindRoleList 查询角色列表(带分页)
func FindRoleList(role *SystemRole, page, pageSize int) ([]SystemRole, int64, error) {
	var roles []SystemRole
	var total int64
	db := pgdb.GetClient()

	// 构建基础查询
	query := db.Model(&SystemRole{})

	// 应用过滤条件
	if role.Name != "" {
		query = query.Where("name LIKE ?", "%"+role.Name+"%")
	}
	if role.Status != 0 {
		query = query.Where("status = ?", role.Status)
	}
	if role.TenantID != 0 {
		query = query.Where("tenant_id = ?", role.TenantID)
	}

	// 获取符合条件的总记录数
	if err := query.Count(&total).Error; err != nil {
		zap.L().Error("failed to count role list", zap.Error(err))
		return nil, 0, err
	}

	// 构建排序和预加载
	queryWithPreload := query.Preload("SystemUsers", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "role_id", "created_at", "updated_at") // 只选择需要的字段
	}).Order("id DESC")

	// 判断是否需要分页
	if page == config.CancelPage && pageSize == config.CancelPageSize {
		// 不分页，获取所有数据
		if err := queryWithPreload.Find(&roles).Error; err != nil {
			zap.L().Error("failed to find all role list", zap.Error(err))
			return nil, 0, err
		}
	} else {
		// 应用分页并获取数据
		if err := queryWithPreload.Offset((page - 1) * pageSize).
			Limit(pageSize).
			Find(&roles).Error; err != nil {
			zap.L().Error("failed to find role list with pagination", zap.Error(err))
			return nil, 0, err
		}
	}

	return roles, total, nil
}
