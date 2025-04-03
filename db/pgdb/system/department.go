package system

import (
	"go.uber.org/zap"

	"gorm.io/gorm"

	"api-server/db/pgdb"
)

// FindDepartmentList 查询部门列表（包含有限的用户信息）
func FindDepartmentList(department *SystemDepartment) ([]SystemDepartment, error) {
	var departments []SystemDepartment
	db := pgdb.GetClient()

	// 构建查询条件
	query := db.Preload("SystemUsers", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "department_id", "created_at", "updated_at") // 只选择需要的字段
	})

	// 如果提供了名称，使用模糊查询
	if department.Name != "" {
		query = query.Where("name LIKE ?", "%s"+department.Name+"%s")
	}

	if err := query.Find(&departments).Error; err != nil {
		zap.L().Error("failed to find department list", zap.Error(err))
		return nil, err
	}
	return departments, nil
}

// GetDepartment 查询单个部门（包含有限的用户信息）
func GetDepartment(department *SystemDepartment) error {
	if err := pgdb.GetClient().Preload("SystemUsers", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "department_id", "created_at", "updated_at") // 只选择需要的字段
	}).Where(department).First(department).Error; err != nil {
		zap.L().Error("failed to get department", zap.Error(err))
		return err
	}
	return nil
}

func AddDepartment(department *SystemDepartment) error {
	if err := pgdb.GetClient().Create(&department).Error; err != nil {
		zap.L().Error("failed to create department", zap.Error(err))
		return err
	}
	return nil
}

func UpdateDepartment(department *SystemDepartment) error {
	if err := pgdb.GetClient().Save(&department).Error; err != nil {
		zap.L().Error("failed to update department", zap.Error(err))
		return err
	}
	return nil
}

func DeleteDepartment(department *SystemDepartment) error {
	if err := pgdb.GetClient().Delete(&department).Error; err != nil {
		zap.L().Error("failed to delete department", zap.Error(err))
		return err
	}
	return nil
}
