package system

import (
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/config"
	"api-server/db/pgdb"
	"api-server/util/encryption"
)

func encryptionPWD(password string) string {
	return encryption.MD5WithSalt(config.PWDSalt + password)
}

// 查询用户
func VerifyUser(userName, password string) (SystemUser, error) {
	user := SystemUser{}
	err := pgdb.GetClient().Where(&SystemUser{Username: userName, Password: encryptionPWD(password)}).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, nil
		}
		zap.L().Error("failed to get user", zap.Error(err))
		return user, err
	}
	return user, nil
}

func GetUser(user *SystemUser) error {
	if err := pgdb.GetClient().Where(user).First(user).Error; err != nil {
		zap.L().Error("failed to get user", zap.Error(err))
		return err
	}
	return nil
}

// UserWithRelations 包含用户及其关联的角色和部门信息
type UserWithRelations struct {
	SystemUser
	RoleName       string `json:"role_name"`
	RoleDesc       string `json:"role_desc"`
	DepartmentName string `json:"department_name"`
}

func FindUserList(user *SystemUser, page, pageSize int) ([]UserWithRelations, int64, error) {
	var usersWithRelations []UserWithRelations
	var total int64
	db := pgdb.GetClient()
	// 构建基础查询
	baseQuery := db.Table("system_users").
		Joins("left join roles on system_users.role_id = roles.id").
		Joins("left join departments on system_users.department_id = departments.id").
		Where("system_users.deleted_at IS NULL")
	// 使用模糊查询
	if user.Username != "" {
		baseQuery = baseQuery.Where("system_users.username LIKE ?", "%"+user.Username+"%")
	}
	if user.Name != "" {
		baseQuery = baseQuery.Where("system_users.name LIKE ?", "%"+user.Name+"%")
	}
	if user.Phone != "" {
		baseQuery = baseQuery.Where("system_users.phone LIKE ?", "%"+user.Phone+"%")
	}
	if user.RoleID != 0 {
		baseQuery = baseQuery.Where("system_users.role_id = ?", user.RoleID)
	}
	if user.DepartmentID != 0 {
		baseQuery = baseQuery.Where("system_users.department_id = ?", user.DepartmentID)
	}
	// 获取符合条件的总记录数
	baseQuery.Count(&total)
	// 应用分页并获取数据
	query := baseQuery.Select("system_users.*, system_roles.name as role_name, system_roles.desc as role_desc, system_departments.name as department_name")
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&usersWithRelations).Error; err != nil {
		zap.L().Error("failed to find user list", zap.Error(err))
		return nil, 0, err
	}
	return usersWithRelations, total, nil
}

func AddUser(user *SystemUser) error {
	if err := pgdb.GetClient().Create(user).Error; err != nil {
		zap.L().Error("failed to add user", zap.Error(err))
		return err
	}
	return nil
}

func UpdateUser(user *SystemUser) error {
	if user.Password != "" {
		user.Password = encryptionPWD(user.Password)
	}
	if err := pgdb.GetClient().Save(user).Error; err != nil {
		zap.L().Error("failed to update user", zap.Error(err))
		return err
	}
	return nil
}

func DeleteUser(user *SystemUser) error {
	if err := pgdb.GetClient().Delete(user).Error; err != nil {
		zap.L().Error("failed to delete user", zap.Error(err))
		return err
	}
	return nil
}
