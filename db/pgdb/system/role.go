package system

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/db/pgdb"
)

func FindRoleList(role *Role) ([]Role, error) {
	var roles []Role
	db := pgdb.GetClient()

	// 构建查询条件
	query := db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "role_id", "created_at", "updated_at") // 只选择需要的用户字段
	})

	// 如果提供了名称，使用模糊查询
	if role.Name != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", role.Name))
	}

	// 应用其他过滤条件
	if role.Status > 0 {
		query = query.Where("status = ?", role.Status)
	}

	if role.ID > 0 {
		query = query.Where("id = ?", role.ID)
	}

	if err := query.Find(&roles).Error; err != nil {
		zap.L().Error("failed to find role list", zap.Error(err))
		return nil, err
	}
	return roles, nil
}

func UpdateRole(role *Role) error {
	if err := pgdb.GetClient().Save(&role).Error; err != nil {
		zap.L().Error("failed to update role", zap.Error(err))
		return err
	}
	return nil
}

func AddRole(role *Role) error {
	if err := pgdb.GetClient().Create(&role).Error; err != nil {
		zap.L().Error("failed to create role", zap.Error(err))
		return err
	}
	return nil
}

func DeleteRole(role *Role) error {
	// 先检查角色下是否有用户
	var count int64
	if err := pgdb.GetClient().Model(&User{}).Where("role_id = ?", role.ID).Count(&count).Error; err != nil {
		zap.L().Error("failed to check role users", zap.Error(err))
		return err
	}

	// 如果有用户，不允许删除
	if count > 0 {
		errMsg := fmt.Sprintf("角色[ID:%d]下有%d个用户，请先删除或转移用户", role.ID, count)
		zap.L().Warn(errMsg)
		return fmt.Errorf(errMsg)
	}

	if err := pgdb.GetClient().Delete(&role).Error; err != nil {
		zap.L().Error("failed to delete role", zap.Error(err))
		return err
	}
	return nil
}
