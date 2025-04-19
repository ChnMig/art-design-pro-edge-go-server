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

// 记录用户登录日志
func CreateLoginLog(log *SystemUserLoginLog) error {
	if err := pgdb.GetClient().Create(log).Error; err != nil {
		zap.L().Error("failed to record login log", zap.Error(err))
		return err
	}
	return nil
}

// FindLoginLogList 查询登录日志列表，支持分页和按用户名、IP查询
func FindLoginLogList(loginLog *SystemUserLoginLog, page, pageSize int) ([]SystemUserLoginLog, int64, error) {
	var loginLogs []SystemUserLoginLog
	var total int64
	db := pgdb.GetClient()

	// 构建基础查询
	baseQuery := db.Model(&SystemUserLoginLog{}).Where("deleted_at IS NULL")

	// 使用模糊查询
	if loginLog.UserName != "" {
		baseQuery = baseQuery.Where("user_name LIKE ?", "%"+loginLog.UserName+"%")
	}
	if loginLog.IP != "" {
		baseQuery = baseQuery.Where("ip LIKE ?", "%"+loginLog.IP+"%")
	}

	// 获取符合条件的总记录数
	baseQuery.Count(&total)

	// 判断是否需要分页
	if page == config.CancelPage && pageSize == config.CancelPageSize {
		// 不分页，获取所有数据
		if err := baseQuery.Order("created_at DESC").Find(&loginLogs).Error; err != nil {
			zap.L().Error("failed to find all login logs", zap.Error(err))
			return nil, 0, err
		}
	} else {
		// 应用分页并获取数据
		if err := baseQuery.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&loginLogs).Error; err != nil {
			zap.L().Error("failed to find login logs", zap.Error(err))
			return nil, 0, err
		}
	}

	return loginLogs, total, nil
}

// UserWithRelations 包含用户及其关联的角色和部门信息
type UserWithRelations struct {
	SystemUser     `json:"User"`
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
		Joins("left join system_roles on system_users.role_id = system_roles.id").
		Joins("left join system_departments on system_users.department_id = system_departments.id").
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

	// 准备查询对象，添加选择字段和排序
	query := baseQuery.Select("system_users.*, system_roles.name as role_name, system_roles.desc as role_desc, system_departments.name as department_name")

	// 判断是否需要分页
	if page == config.CancelPage && pageSize == config.CancelPageSize {
		// 不分页，获取所有数据
		if err := query.Find(&usersWithRelations).Error; err != nil {
			zap.L().Error("failed to find all user list", zap.Error(err))
			return nil, 0, err
		}
	} else {
		// 应用分页并获取数据
		if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&usersWithRelations).Error; err != nil {
			zap.L().Error("failed to find user list", zap.Error(err))
			return nil, 0, err
		}
	}

	return usersWithRelations, total, nil
}

func AddUser(user *SystemUser) error {
	user.Password = encryptionPWD(user.Password)
	if err := pgdb.GetClient().Create(user).Error; err != nil {
		zap.L().Error("failed to add user", zap.Error(err))
		return err
	}
	return nil
}

func GetUser(user *SystemUser) error {
	if err := pgdb.GetClient().Where(user).First(user).Error; err != nil {
		zap.L().Error("failed to get user", zap.Error(err))
		return err
	}
	return nil
}

func UpdateUser(user *SystemUser) error {
	if user.Password != "" {
		user.Password = encryptionPWD(user.Password)
	}
	if err := pgdb.GetClient().Updates(user).Error; err != nil {
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

// FindAllUsers 查询所有用户
func FindAllUsers(users *[]SystemUser) error {
	if err := pgdb.GetClient().Find(users).Error; err != nil {
		zap.L().Error("failed to find all users", zap.Error(err))
		return err
	}
	return nil
}
