package system

import (
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/config"
	"api-server/db/pgdb"
	"api-server/util/encryption"
)

// encryptionPWD 已移除MD5支持，新安装只支持bcrypt

// HashPassword 使用bcrypt安全哈希密码
func HashPassword(password string) (string, error) {
	return encryption.HashPasswordWithBcrypt(password)
}

// VerifyPassword 验证密码，新安装只支持bcrypt格式
func VerifyPassword(password, hashedPassword, passwordType string) bool {
	return encryption.VerifyBcryptPassword(password, hashedPassword)
}

// VerifyUser 验证用户登录（多租户版本）
func VerifyUser(tenantCode, account, password string) (SystemUser, SystemTenant, error) {
	user := SystemUser{}
	tenant := SystemTenant{}

	// 首先验证租户
	tenant, err := GetTenantByCode(tenantCode)
	if err != nil {
		zap.L().Error("failed to get tenant", zap.Error(err))
		return user, tenant, err
	}
	if tenant.ID == 0 {
		// 租户不存在
		return user, tenant, nil
	}

	// 验证租户状态
	if err := ValidateTenant(&tenant); err != nil {
		zap.L().Error("tenant validation failed", zap.Error(err))
		return user, tenant, err
	}

	// 根据租户ID和账号查找用户
	err = pgdb.GetClient().Where("tenant_id = ? AND account = ?", tenant.ID, account).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, tenant, nil // 返回空用户，ID为0表示未找到
		}
		zap.L().Error("failed to get user", zap.Error(err))
		return user, tenant, err
	}

	// 验证密码
	if !VerifyPassword(password, user.Password, "bcrypt") {
		// 密码错误，返回空用户
		return SystemUser{}, tenant, nil
	}

	return user, tenant, nil
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
	
	// 租户过滤（必须）
	if user.TenantID != 0 {
		baseQuery = baseQuery.Where("system_users.tenant_id = ?", user.TenantID)
	}
	
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
	// 新用户使用bcrypt加密
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		zap.L().Error("failed to hash user password", zap.Error(err))
		return err
	}

	user.Password = hashedPassword

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
	// 如果更新密码，使用bcrypt加密
	if user.Password != "" {
		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			zap.L().Error("failed to hash user password for update", zap.Error(err))
			return err
		}
		user.Password = hashedPassword
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
