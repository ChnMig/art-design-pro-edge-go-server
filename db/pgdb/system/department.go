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
	if department.TenantID == 0 {
		return nil, 0, gorm.ErrInvalidData
	}
	query, err := buildDepartmentListQuery(pgdb.GetClient(), department)
	if err != nil {
		return nil, 0, err
	}

	// 获取符合条件的总记录数
	if err := query.Count(&total).Error; err != nil {
		zap.L().Error("failed to count department list", zap.Error(err))
		return nil, 0, err
	}

	// 构建排序
	queryOrder := query.Order("sort ASC, id ASC")

	// 判断是否需要分页
	if page == config.CancelPage && pageSize == config.CancelPageSize {
		// 不分页，获取所有数据
		if err := queryOrder.Find(&departments).Error; err != nil {
			zap.L().Error("failed to find all department list", zap.Error(err))
			return nil, 0, err
		}
	} else {
		// 应用分页并获取数据
		if err := queryOrder.Offset((page - 1) * pageSize).
			Limit(pageSize).
			Find(&departments).Error; err != nil {
			zap.L().Error("failed to find department list with pagination", zap.Error(err))
			return nil, 0, err
		}
	}

	return departments, total, nil
}

func buildDepartmentListQuery(db *gorm.DB, department *SystemDepartment) (*gorm.DB, error) {
	if department.TenantID == 0 {
		return nil, gorm.ErrInvalidData
	}
	query := db.Model(&SystemDepartment{}).Where("tenant_id = ?", department.TenantID)
	if department.Name != "" {
		query = query.Where("name LIKE ?", "%"+department.Name+"%")
	}
	if department.Status != 0 {
		query = query.Where("status = ?", department.Status)
	}
	return query, nil
}

// GetDepartment 查询单个部门
func GetDepartment(department *SystemDepartment) error {
	if department.ID == 0 || department.TenantID == 0 {
		return gorm.ErrInvalidData
	}
	if err := pgdb.GetClient().Where(department).First(department).Error; err != nil {
		zap.L().Error("failed to get department", zap.Error(err))
		return err
	}
	return nil
}

func AddDepartment(department *SystemDepartment) error {
	if department.TenantID == 0 {
		return gorm.ErrInvalidData
	}
	if err := pgdb.GetClient().Create(department).Error; err != nil {
		zap.L().Error("failed to create department", zap.Error(err))
		return err
	}
	return nil
}

func UpdateDepartment(department *SystemDepartment) error {
	if department.ID == 0 || department.TenantID == 0 {
		return gorm.ErrInvalidData
	}
	result := scopedDepartmentUpdate(pgdb.GetClient(), department)
	if result.Error != nil {
		zap.L().Error("failed to update department", zap.Error(result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func scopedDepartmentUpdate(db *gorm.DB, department *SystemDepartment) *gorm.DB {
	return db.Model(&SystemDepartment{}).
		Where("id = ? AND tenant_id = ?", department.ID, department.TenantID).
		Select("name", "status", "sort").
		Updates(department)
}

func DeleteDepartment(department *SystemDepartment) error {
	if department.ID == 0 || department.TenantID == 0 {
		return gorm.ErrInvalidData
	}
	result := scopedDepartmentDelete(pgdb.GetClient(), department)
	if result.Error != nil {
		zap.L().Error("failed to delete department", zap.Error(result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func scopedDepartmentDelete(db *gorm.DB, department *SystemDepartment) *gorm.DB {
	return db.Where("id = ? AND tenant_id = ?", department.ID, department.TenantID).Delete(&SystemDepartment{})
}

// CountUsersByDepartmentID 统计指定部门下的用户数量
func CountUsersByDepartmentID(tenantID, departmentID uint, count *int64) error {
	if tenantID == 0 || departmentID == 0 {
		return gorm.ErrInvalidData
	}
	if err := pgdb.GetClient().Model(&SystemUser{}).
		Where("tenant_id = ? AND department_id = ?", tenantID, departmentID).
		Count(count).Error; err != nil {
		zap.L().Error("failed to count users by department id", zap.Error(err))
		return err
	}
	return nil
}
