package system

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/config"
	"api-server/db/pgdb"
)

// GetUserMenuData 获取用户菜单数据
func GetUserMenuData(userID uint) ([]SystemMenu, []SystemMenuAuth, error) {
	// 获取用户信息及其角色
	var user SystemUser
	if err := pgdb.GetClient().Where(&SystemUser{Model: gorm.Model{ID: userID}}).First(&user).Error; err != nil {
		zap.L().Error("failed to get user", zap.Error(err))
		return nil, nil, err
	}
	// 获取该角色关联的所有菜单(包括权限)
	var role SystemRole
	if err := pgdb.GetClient().Preload("SystemMenus").
		Preload("SystemMenuAuths").
		Where("id = ?", user.RoleID).
		First(&role).Error; err != nil {
		zap.L().Error("failed to get role", zap.Error(err))
		return nil, nil, err
	}
	return role.SystemMenus, role.SystemMenuAuths, nil
}

// 获取菜单树(不带分页)
func GetMenuData() ([]SystemMenu, []SystemMenuAuth, error) {
	var menus []SystemMenu
	if err := pgdb.GetClient().Find(&menus).Error; err != nil {
		zap.L().Error("failed to get menus", zap.Error(err))
		return nil, nil, err
	}
	var Auths []SystemMenuAuth
	if err := pgdb.GetClient().Find(&Auths).Error; err != nil {
		zap.L().Error("failed to get menu Auths", zap.Error(err))
		return nil, nil, err
	}
	return menus, Auths, nil
}

// GetMenuDataByRoleID 获取指定角色ID的菜单和权限数据
func GetMenuDataByRoleID(roleID uint) ([]SystemMenu, []SystemMenuAuth, []uint, []uint, error) {
	// 获取所有菜单
	var allMenus []SystemMenu
	if err := pgdb.GetClient().Find(&allMenus).Error; err != nil {
		zap.L().Error("failed to get all menus", zap.Error(err))
		return nil, nil, nil, nil, err
	}
	// 获取所有权限
	var allAuths []SystemMenuAuth
	if err := pgdb.GetClient().Find(&allAuths).Error; err != nil {
		zap.L().Error("failed to get all menu auths", zap.Error(err))
		return nil, nil, nil, nil, err
	}
	// 获取角色拥有的菜单ID列表
	var role SystemRole
	if err := pgdb.GetClient().Preload("SystemMenus").
		Preload("SystemMenuAuths").
		Where("id = ?", roleID).
		First(&role).Error; err != nil {
		zap.L().Error("failed to get role with menus", zap.Error(err))
		return nil, nil, nil, nil, err
	}
	// 提取角色拥有的菜单ID和权限ID
	var roleMenuIds []uint
	var roleAuthIds []uint
	for _, m := range role.SystemMenus {
		roleMenuIds = append(roleMenuIds, m.ID)
	}
	for _, a := range role.SystemMenuAuths {
		roleAuthIds = append(roleAuthIds, a.ID)
	}
	return allMenus, allAuths, roleMenuIds, roleAuthIds, nil
}

// 新增一个菜单
func AddMenu(menu *SystemMenu) error {
	if err := pgdb.GetClient().Create(menu).Error; err != nil {
		zap.L().Error("failed to create menu", zap.Error(err))
		return err
	}
	return nil
}

// 删除一个菜单
func DeleteMenu(menu *SystemMenu) error {
	if err := pgdb.GetClient().Delete(menu).Error; err != nil {
		zap.L().Error("failed to delete menu", zap.Error(err))
		return err
	}
	return nil
}

func UpdateMenu(menu *SystemMenu) error {
	if err := pgdb.GetClient().Omit("created_at").Save(menu).Error; err != nil {
		zap.L().Error("failed to update menu", zap.Error(err))
		return err
	}
	return nil
}

func GetMenu(menu *SystemMenu) error {
	if err := pgdb.GetClient().Where(menu).First(menu).Error; err != nil {
		zap.L().Error("failed to get menu", zap.Error(err))
		return err
	}
	return nil
}

// FindMenuList 查询菜单列表(带分页)
func FindMenuList(menu *SystemMenu, page, pageSize int) ([]SystemMenu, int64, error) {
	var menus []SystemMenu
	var total int64
	db := pgdb.GetClient()

	// 构建基础查询
	query := db.Model(&SystemMenu{})

	// 应用过滤条件
	if menu.Title != "" {
		query = query.Where("title LIKE ?", "%"+menu.Title+"%")
	}
	if menu.Name != "" {
		query = query.Where("name LIKE ?", "%"+menu.Name+"%")
	}
	if menu.Path != "" {
		query = query.Where("path LIKE ?", "%"+menu.Path+"%")
	}
	if menu.ParentID != 0 {
		query = query.Where("parent_id = ?", menu.ParentID)
	}
	if menu.Status != 0 {
		query = query.Where("status = ?", menu.Status)
	}

	// 获取符合条件的总记录数
	if err := query.Count(&total).Error; err != nil {
		zap.L().Error("failed to count menu list", zap.Error(err))
		return nil, 0, err
	}

	// 构建排序
	queryOrder := query.Order("sort ASC, id ASC")

	// 判断是否需要分页
	if page == config.CancelPage && pageSize == config.CancelPageSize {
		// 不分页，获取所有数据
		if err := queryOrder.Find(&menus).Error; err != nil {
			zap.L().Error("failed to find all menu list", zap.Error(err))
			return nil, 0, err
		}
	} else {
		// 应用分页并获取数据
		if err := queryOrder.Offset((page - 1) * pageSize).
			Limit(pageSize).
			Find(&menus).Error; err != nil {
			zap.L().Error("failed to find menu list with pagination", zap.Error(err))
			return nil, 0, err
		}
	}

	return menus, total, nil
}

func FilterMenusByIDs(allMenus []SystemMenu, allPermissions []SystemMenuAuth, allowedMenuIDs []uint) ([]SystemMenu, []SystemMenuAuth) {
	if len(allowedMenuIDs) == 0 {
		return []SystemMenu{}, []SystemMenuAuth{}
	}
	menuMap := make(map[uint]SystemMenu, len(allMenus))
	for _, menu := range allMenus {
		menuMap[menu.ID] = menu
	}
	allowedSet := make(map[uint]struct{}, len(allowedMenuIDs))
	for _, id := range allowedMenuIDs {
		if _, exists := menuMap[id]; exists {
			allowedSet[id] = struct{}{}
		}
	}
	// 补充父级菜单
	for id := range allowedSet {
		current := menuMap[id]
		for current.ParentID != 0 {
			if _, exists := allowedSet[current.ParentID]; exists {
				break
			}
			if parent, ok := menuMap[current.ParentID]; ok {
				allowedSet[parent.ID] = struct{}{}
				current = parent
			} else {
				break
			}
		}
	}
	filteredMenus := make([]SystemMenu, 0, len(allowedSet))
	for _, menu := range allMenus {
		if _, ok := allowedSet[menu.ID]; ok {
			filteredMenus = append(filteredMenus, menu)
		}
	}
	filteredPermissions := make([]SystemMenuAuth, 0)
	for _, perm := range allPermissions {
		if _, ok := allowedSet[perm.MenuID]; ok {
			filteredPermissions = append(filteredPermissions, perm)
		}
	}
	return filteredMenus, filteredPermissions
}

func FilterUintIDs(source []uint, allowedIDs []uint) []uint {
	if len(source) == 0 {
		return source
	}
	if len(allowedIDs) == 0 {
		return []uint{}
	}
	allowedSet := make(map[uint]struct{}, len(allowedIDs))
	for _, id := range allowedIDs {
		allowedSet[id] = struct{}{}
	}
	result := make([]uint, 0, len(source))
	for _, id := range source {
		if _, ok := allowedSet[id]; ok {
			result = append(result, id)
		}
	}
	return result
}
