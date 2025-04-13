package system

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/config"
	"api-server/db/pgdb"
)

// FindDepartmentList 查询部门列表(带分页)
func FindDepartmentList(department *SystemDepartment, page, pageSize int) ([]SystemDepartment, int64, error) {
	var departments []SystemDepartment
	var total int64
	db := pgdb.GetClient()

	// 构建基础查询
	query := db.Model(&SystemDepartment{})

	// 应用过滤条件
	if department.Name != "" {
		query = query.Where("name LIKE ?", "%"+department.Name+"%")
	}
	if department.Status != 0 {
		query = query.Where("status = ?", department.Status)
	}

	// 获取符合条件的总记录数
	if err := query.Count(&total).Error; err != nil {
		zap.L().Error("failed to count department list", zap.Error(err))
		return nil, 0, err
	}

	// 构建排序和预加载
	queryWithPreload := query.Preload("SystemUsers", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "department_id", "created_at", "updated_at") // 只选择需要的字段
	}).Order("sort DESC, id DESC")

	// 判断是否需要分页
	if page == config.CancelPage && pageSize == config.CancelPageSize {
		// 不分页，获取所有数据
		if err := queryWithPreload.Find(&departments).Error; err != nil {
			zap.L().Error("failed to find all department list", zap.Error(err))
			return nil, 0, err
		}
	} else {
		// 应用分页并获取数据
		if err := queryWithPreload.Offset((page - 1) * pageSize).
			Limit(pageSize).
			Find(&departments).Error; err != nil {
			zap.L().Error("failed to find department list with pagination", zap.Error(err))
			return nil, 0, err
		}
	}

	return departments, total, nil
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
