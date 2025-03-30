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
func VerifyUser(userName, password string) (User, error) {
	user := User{}
	err := pgdb.GetClient().Where(&User{Username: userName, Password: encryptionPWD(password)}).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, nil
		}
		zap.L().Error("failed to get user", zap.Error(err))
		return user, err
	}
	return user, nil
}

func GetUser(user *User) error {
	if err := pgdb.GetClient().Where(user).First(user).Error; err != nil {
		zap.L().Error("failed to get user", zap.Error(err))
		return err
	}
	return nil
}

// UserWithRelations 包含用户及其关联的角色和部门信息
type UserWithRelations struct {
	User
	RoleName       string `json:"role_name"`
	RoleDesc       string `json:"role_desc"`
	DepartmentName string `json:"department_name"`
}

func FindUserList(user *User, page, pageSize int) ([]UserWithRelations, int64, error) {
	var usersWithRelations []UserWithRelations
	var total int64
	db := pgdb.GetClient()
	// 构建基础查询
	baseQuery := db.Table("users").
		Joins("left join roles on users.role_id = roles.id").
		Joins("left join departments on users.department_id = departments.id").
		Where("users.deleted_at IS NULL")
	// 使用模糊查询
	if user.Username != "" {
		baseQuery = baseQuery.Where("users.username LIKE ?", "%"+user.Username+"%")
	}
	if user.Name != "" {
		baseQuery = baseQuery.Where("users.name LIKE ?", "%"+user.Name+"%")
	}
	if user.Phone != "" {
		baseQuery = baseQuery.Where("users.phone LIKE ?", "%"+user.Phone+"%")
	}
	if user.RoleID != 0 {
		baseQuery = baseQuery.Where("users.role_id = ?", user.RoleID)
	}
	if user.DepartmentID != 0 {
		baseQuery = baseQuery.Where("users.department_id = ?", user.DepartmentID)
	}
	// 获取符合条件的总记录数
	baseQuery.Count(&total)
	// 应用分页并获取数据
	query := baseQuery.Select("users.*, roles.name as role_name, roles.desc as role_desc, departments.name as department_name")
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&usersWithRelations).Error; err != nil {
		zap.L().Error("failed to find user list", zap.Error(err))
		return nil, 0, err
	}
	return usersWithRelations, total, nil
}

func AddUser(user *User) error {
	if err := pgdb.GetClient().Create(user).Error; err != nil {
		zap.L().Error("failed to add user", zap.Error(err))
		return err
	}
	return nil
}

func UpdateUser(user *User) error {
	if user.Password != "" {
		user.Password = encryptionPWD(user.Password)
	}
	if err := pgdb.GetClient().Save(user).Error; err != nil {
		zap.L().Error("failed to update user", zap.Error(err))
		return err
	}
	return nil
}

func DeleteUser(user *User) error {
	if err := pgdb.GetClient().Delete(user).Error; err != nil {
		zap.L().Error("failed to delete user", zap.Error(err))
		return err
	}
	return nil
}
