package systemuser

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/db/pgdb/system"
	"api-server/db/rdb"
)

const (
	// UserInfoKey 用户信息缓存键前缀
	UserInfoKey = "system:user:info:"
	// UserListKey 用户列表缓存键
	UserListKey = "system:user:list"
	// UserCacheExpiration 用户缓存过期时间（12小时）
	UserCacheExpiration = 12 * time.Hour
)

// UserCacheInfo 用户缓存信息
type UserCacheInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"` // 昵称
	Name     string `json:"name"`     // 姓名
	RoleID   uint   `json:"role_id"`
	RoleName string `json:"role_name"`
}

// CacheAllUsers 缓存所有用户信息到Redis
func CacheAllUsers() error {
	client := rdb.GetClient()

	// 获取所有用户信息
	var users []system.SystemUser
	if err := system.FindAllUsers(&users); err != nil {
		zap.L().Error("获取所有用户信息失败", zap.Error(err))
		return err
	}

	// 获取所有角色信息，用于映射角色名称
	var roles []system.SystemRole
	if err := system.FindAllRoles(&roles); err != nil {
		zap.L().Error("获取所有角色信息失败", zap.Error(err))
		return err
	}

	// 创建角色映射表
	roleMap := make(map[uint]string)
	for _, role := range roles {
		roleMap[role.ID] = role.Name
	}

	// 使用管道批量操作，提高效率
	pipe := client.Pipeline()

	// 缓存用户列表
	var userList []UserCacheInfo
	for _, user := range users {
		// 创建用户缓存对象
		userCache := UserCacheInfo{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
			RoleID:   user.RoleID,
			RoleName: roleMap[user.RoleID],
		}

		// 将用户信息添加到列表
		userList = append(userList, userCache)

		// 单独缓存每个用户信息，方便通过ID快速获取
		userJSON, err := json.Marshal(userCache)
		if err != nil {
			zap.L().Error("序列化用户信息失败", zap.Error(err))
			continue
		}
		pipe.Set(UserInfoKey+strconv.FormatUint(uint64(user.ID), 10), userJSON, UserCacheExpiration)
	}

	// 缓存完整用户列表
	listJSON, err := json.Marshal(userList)
	if err != nil {
		zap.L().Error("序列化用户列表失败", zap.Error(err))
		return err
	}
	pipe.Set(UserListKey, listJSON, UserCacheExpiration)

	// 执行管道操作
	_, err = pipe.Exec()
	if err != nil {
		zap.L().Error("缓存用户信息到Redis失败", zap.Error(err))
		return err
	}

	zap.L().Info("已成功缓存所有用户信息到Redis", zap.Int("用户数量", len(users)))
	return nil
}

// GetUserFromCache 从缓存中获取用户信息
func GetUserFromCache(userID uint) (*UserCacheInfo, error) {
	client := rdb.GetClient()

	// 从Redis获取用户信息
	val, err := client.Get(UserInfoKey + strconv.FormatUint(uint64(userID), 10)).Result()
	if err != nil {
		if err == redis.Nil {
			// 缓存未命中，尝试单独获取并缓存该用户
			return cacheUserByID(userID)
		}
		zap.L().Error("从Redis获取用户信息失败", zap.Error(err))
		return nil, err
	}

	// 反序列化用户信息
	var userCache UserCacheInfo
	if err = json.Unmarshal([]byte(val), &userCache); err != nil {
		zap.L().Error("反序列化用户信息失败", zap.Error(err))
		return nil, err
	}

	return &userCache, nil
}

// GetAllUsersFromCache 从缓存中获取所有用户列表
func GetAllUsersFromCache() ([]UserCacheInfo, error) {
	client := rdb.GetClient()

	// 从Redis获取用户列表
	val, err := client.Get(UserListKey).Result()
	if err != nil {
		if err == redis.Nil {
			// 缓存未命中，重新缓存所有用户
			if err = CacheAllUsers(); err != nil {
				return nil, err
			}
			// 再次尝试获取
			val, err = client.Get(UserListKey).Result()
			if err != nil {
				zap.L().Error("从Redis获取用户列表失败", zap.Error(err))
				return nil, err
			}
		} else {
			zap.L().Error("从Redis获取用户列表失败", zap.Error(err))
			return nil, err
		}
	}

	// 反序列化用户列表
	var userList []UserCacheInfo
	if err = json.Unmarshal([]byte(val), &userList); err != nil {
		zap.L().Error("反序列化用户列表失败", zap.Error(err))
		return nil, err
	}

	return userList, nil
}

// cacheUserByID 单独获取并缓存指定ID的用户
func cacheUserByID(userID uint) (*UserCacheInfo, error) {
	client := rdb.GetClient()

	// 获取用户信息
	user := system.SystemUser{Model: gorm.Model{ID: userID}}
	if err := system.GetUser(&user); err != nil {
		zap.L().Error("获取用户信息失败", zap.Error(err))
		return nil, err
	}

	// 获取角色信息
	role := system.SystemRole{Model: gorm.Model{ID: user.RoleID}}
	if err := system.GetRole(&role); err != nil {
		zap.L().Error("获取角色信息失败", zap.Error(err))
		// 继续执行，只是角色名称可能为空
	}

	// 创建用户缓存对象
	userCache := UserCacheInfo{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		RoleID:   user.RoleID,
		RoleName: role.Name,
	}

	// 序列化用户信息
	userJSON, err := json.Marshal(userCache)
	if err != nil {
		zap.L().Error("序列化用户信息失败", zap.Error(err))
		return nil, err
	}

	// 缓存用户信息
	if err = client.Set(UserInfoKey+strconv.FormatUint(uint64(user.ID), 10), userJSON, UserCacheExpiration).Err(); err != nil {
		zap.L().Error("缓存用户信息到Redis失败", zap.Error(err))
		return nil, err
	}

	return &userCache, nil
}
