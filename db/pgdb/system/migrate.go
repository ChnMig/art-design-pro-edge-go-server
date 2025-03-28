package system

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/config"
)

func migrateTable(db *gorm.DB) error {
	err := db.AutoMigrate(&Department{}, &Role{}, &Menu{}, &MenuAuth{}, &User{})
	if err != nil {
		zap.L().Error("failed to migrate system model", zap.Error(err))
		return err
	}
	return nil
}

func migrateData(db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		// 创建菜单
		menus := []Menu{
			{Model: gorm.Model{ID: 1}, Path: "/dashboard", Name: "Dashboard", Component: "/layout/index", Title: "仪表盘", Icon: "&#xe721;", KeepAlive: 2, Status: 1, Level: 1, ParentID: 0, Sort: 99},
			{Model: gorm.Model{ID: 2}, Path: "/system", Name: "System", Component: "/layout/index", Title: "系统管理", Icon: "&#xe72b;", KeepAlive: 2, Status: 1, Level: 1, ParentID: 0, Sort: 20},
			{Model: gorm.Model{ID: 3}, Path: "menu", Name: "Menu", Component: "/system/menu/index", Title: "菜单管理", KeepAlive: 1, Status: 1, Level: 2, ParentID: 2, Sort: 99},
			{Model: gorm.Model{ID: 4}, Path: "role", Name: "Role", Component: "/system/role/index", Title: "角色管理", KeepAlive: 1, Status: 1, Level: 2, ParentID: 2, Sort: 88},
			{Model: gorm.Model{ID: 5}, Path: "department", Name: "Department", Component: "/system/department/index", Title: "部门管理", KeepAlive: 1, Status: 1, Level: 2, ParentID: 2, Sort: 77},
			{Model: gorm.Model{ID: 6}, Path: "user", Name: "User", Component: "/system/user/index", Title: "用户管理", KeepAlive: 1, Status: 1, Level: 2, ParentID: 2, Sort: 66},
			{Model: gorm.Model{ID: 7}, Path: "console", Name: "Console", Component: "/dashboard/console", Title: "工作台", Icon: "", KeepAlive: 1, Status: 1, Level: 2, ParentID: 1, Sort: 99},
			{Model: gorm.Model{ID: 8}, Path: "analysis", Name: "Analysis", Component: "/dashboard/analysis", Title: "分析页", Icon: "", KeepAlive: 1, Status: 1, Level: 2, ParentID: 1, Sort: 88},
			{Model: gorm.Model{ID: 9}, Path: "myInfo", Name: "MyInfo", Component: "/system/myInfo/index", Title: "个人中心", Icon: "", KeepAlive: 1, Status: 1, Level: 2, ParentID: 2, Sort: 55},
		}
		err := db.Create(&menus).Error
		if err != nil {
			zap.L().Error("failed to create menu", zap.Error(err))
			return err
		}
		// 创建角色
		roles := []Role{
			{Model: gorm.Model{ID: 1}, Name: "超级管理员", Desc: "拥有所有权限", Status: 1},
			{Model: gorm.Model{ID: 2}, Name: "普通用户", Desc: "普通用户", Status: 1},
		}
		err = db.Create(&roles).Error
		if err != nil {
			zap.L().Error("failed to create role", zap.Error(err))
			return err
		}
		// 为角色分配菜单权限
		// 超级管理员拥有所有菜单权限
		adminRole := Role{}
		err = db.First(&adminRole, 1).Error
		if err != nil {
			zap.L().Error("failed to find admin role", zap.Error(err))
			return err
		}
		// 为超级管理员分配所有菜单
		var allMenus []Menu
		err = db.Find(&allMenus).Error
		if err != nil {
			zap.L().Error("failed to find menus", zap.Error(err))
			return err
		}
		err = db.Model(&adminRole).Association("Menus").Append(&allMenus)
		if err != nil {
			zap.L().Error("failed to associate menus with admin role", zap.Error(err))
			return err
		}
		// 为普通用户分配首页菜单
		normalRole := Role{}
		err = db.First(&normalRole, 2).Error
		if err != nil {
			zap.L().Error("failed to find normal role", zap.Error(err))
			return err
		}
		var dashboardMenu Menu
		err = db.First(&dashboardMenu, 1).Error
		if err != nil {
			zap.L().Error("failed to find dashboard menu", zap.Error(err))
			return err
		}
		err = db.Model(&normalRole).Association("Menus").Append(&dashboardMenu)
		if err != nil {
			zap.L().Error("failed to associate dashboard menu with normal role", zap.Error(err))
			return err
		}
		// 创建菜单按钮权限
		menuAuths := []MenuAuth{
			{Model: gorm.Model{ID: 1}, MenuID: 3, Mark: "add", Title: "新增", Roles: []Role{{Model: gorm.Model{ID: 1}}}},
			{Model: gorm.Model{ID: 2}, MenuID: 3, Mark: "edit", Title: "修改", Roles: []Role{{Model: gorm.Model{ID: 1}}}},
			{Model: gorm.Model{ID: 3}, MenuID: 3, Mark: "delete", Title: "删除", Roles: []Role{{Model: gorm.Model{ID: 1}}}},
		}
		err = db.Create(&menuAuths).Error
		if err != nil {
			zap.L().Error("failed to create menu permission", zap.Error(err))
			return err
		}
		// 创建部门
		departments := []Department{
			{Model: gorm.Model{ID: 1}, Name: "管理中心", Sort: 1, Status: 1},
		}
		err = db.Create(&departments).Error
		if err != nil {
			zap.L().Error("failed to create department", zap.Error(err))
			return err
		}
		// 创建用户
		pwd := encryptionPWD(config.AdminPassword)
		users := []User{
			{Model: gorm.Model{ID: 1}, DepartmentID: 1, RoleID: 1, Name: "超级管理员", Username: "admin", Password: pwd, Status: 1, Gender: 1},
		}
		err = db.Create(&users).Error
		if err != nil {
			zap.L().Error("failed to create user", zap.Error(err))
			return err
		}
		return nil
	})
	return err
}

func resetSequences(db *gorm.DB) error {
	tables := []string{
		"menus", "menu_auths", "roles", "departments", "users",
	}

	for _, table := range tables {
		seqName := table + "_id_seq"
		query := fmt.Sprintf("SELECT setval('%s', (SELECT COALESCE(MAX(id), 0) FROM %s));", seqName, table)
		if err := db.Exec(query).Error; err != nil {
			zap.L().Error("failed to reset sequence", zap.String("sequence", seqName), zap.Error(err))
			return err
		}
		zap.L().Info("sequence reset successfully", zap.String("sequence", seqName))
	}

	return nil
}

func Migrate(db *gorm.DB) error {
	err := migrateTable(db)
	if err != nil {
		return err
	}
	err = migrateData(db)
	if err != nil {
		return err
	}
	// 添加序列重置操作
	err = resetSequences(db)
	if err != nil {
		return err
	}
	return nil
}
