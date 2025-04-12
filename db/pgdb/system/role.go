package system

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/db/pgdb"
)

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
	if page == -1 && pageSize == -1 {
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
