package system

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/db/pgdb"
)

func FindRoleList(role *SystemRole) ([]SystemRole, error) {
	var roles []SystemRole
	db := pgdb.GetClient()
	// 构建查询条件
	query := db.Preload("SystemUsers", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "role_id", "created_at", "updated_at") // 只选择需要的用户字段
	})
	// 如果提供了名称，使用模糊查询
	if role.Name != "" {
		query = query.Where("name LIKE ?", "%s"+role.Name+"%s")
	}
	if err := query.Find(&roles).Error; err != nil {
		zap.L().Error("failed to find role list", zap.Error(err))
		return nil, err
	}
	return roles, nil
}

func UpdateRole(role *SystemRole) error {
	if err := pgdb.GetClient().Save(&role).Error; err != nil {
		zap.L().Error("failed to update role", zap.Error(err))
		return err
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
	if err := pgdb.GetClient().Delete(&role).Error; err != nil {
		zap.L().Error("failed to delete role", zap.Error(err))
		return err
	}
	return nil
}

// GetRole 获取单个角色信息
func GetRole(role *SystemRole) error {
	if err := pgdb.GetClient().Where(role).First(role).Error; err != nil {
		zap.L().Error("failed to get role", zap.Error(err))
		return err
	}
	return nil
}

// FindAllRoles 查询所有角色
func FindAllRoles(roles *[]SystemRole) error {
	if err := pgdb.GetClient().Find(roles).Error; err != nil {
		zap.L().Error("failed to find all roles", zap.Error(err))
		return err
	}
	return nil
}
