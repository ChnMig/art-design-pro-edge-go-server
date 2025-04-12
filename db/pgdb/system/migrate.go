package system

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/config"
)

func migrateTable(db *gorm.DB) error {
	err := db.AutoMigrate(&SystemDepartment{}, &SystemRole{}, &SystemMenu{}, &SystemMenuAuth{}, &SystemUser{}, &SystemUserTodo{}, &SystemUserTodoComments{}, &SystemUserTodoStep{}, &SystemUserTodoLog{})
	if err != nil {
		zap.L().Error("failed to migrate system model", zap.Error(err))
		return err
	}
	return nil
}

func migrateData(db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		// 创建菜单
		menus := []SystemMenu{
			{Model: gorm.Model{ID: 1}, Path: "/dashboard", Name: "Dashboard", Component: "/layout/index", Title: "仪表盘", Icon: "&#xe721;", KeepAlive: 2, Status: 1, Level: 1, ParentID: 0, Sort: 99},
			{Model: gorm.Model{ID: 2}, Path: "/system", Name: "System", Component: "/layout/index", Title: "系统管理", Icon: "&#xe72b;", KeepAlive: 2, Status: 1, Level: 1, ParentID: 0, Sort: 20},
			{Model: gorm.Model{ID: 3}, Path: "menu", Name: "SystemMenu", Component: "/system/menu/index", Title: "菜单管理", KeepAlive: 1, Status: 1, Level: 2, ParentID: 2, Sort: 99},
			{Model: gorm.Model{ID: 4}, Path: "role", Name: "SystemRole", Component: "/system/role/index", Title: "角色管理", KeepAlive: 1, Status: 1, Level: 2, ParentID: 2, Sort: 88},
			{Model: gorm.Model{ID: 5}, Path: "department", Name: "SystemDepartment", Component: "/system/department/index", Title: "部门管理", KeepAlive: 1, Status: 1, Level: 2, ParentID: 2, Sort: 77},
			{Model: gorm.Model{ID: 6}, Path: "user", Name: "SystemUser", Component: "/system/user/index", Title: "用户管理", KeepAlive: 1, Status: 1, Level: 2, ParentID: 2, Sort: 66},
			{Model: gorm.Model{ID: 7}, Path: "console", Name: "DashboardConsole", Component: "/dashboard/console/index", Title: "工作台", Icon: "", KeepAlive: 1, Status: 1, Level: 2, ParentID: 1, Sort: 99},
			{Model: gorm.Model{ID: 8}, Path: "analysis", Name: "DashboardAnalysis", Component: "/dashboard/analysis/index", Title: "分析页", Icon: "", KeepAlive: 1, Status: 1, Level: 2, ParentID: 1, Sort: 88},
		}
		err := db.Create(&menus).Error
		if err != nil {
			zap.L().Error("failed to create menu", zap.Error(err))
			return err
		}
		// 创建角色
		roles := []SystemRole{
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
		adminRole := SystemRole{}
		err = db.First(&adminRole, 1).Error
		if err != nil {
			zap.L().Error("failed to find admin role", zap.Error(err))
			return err
		}
		// 为超级管理员分配所有菜单
		var allMenus []SystemMenu
		err = db.Find(&allMenus).Error
		if err != nil {
			zap.L().Error("failed to find menus", zap.Error(err))
			return err
		}
		err = db.Model(&adminRole).Association("SystemMenus").Append(&allMenus)
		if err != nil {
			zap.L().Error("failed to associate menus with admin role", zap.Error(err))
			return err
		}
		// 为普通用户分配首页菜单
		normalRole := SystemRole{}
		err = db.First(&normalRole, 2).Error
		if err != nil {
			zap.L().Error("failed to find normal role", zap.Error(err))
			return err
		}
		// 为普通用户分配工作台和分析页菜单
		var consoleMenu, analysisMenu, dashboardMenu SystemMenu
		err = db.First(&dashboardMenu, 1).Error
		if err != nil {
			zap.L().Error("failed to find dashboard menu", zap.Error(err))
			return err
		}
		err = db.First(&consoleMenu, 7).Error
		if err != nil {
			zap.L().Error("failed to find console menu", zap.Error(err))
			return err
		}
		err = db.First(&analysisMenu, 8).Error
		if err != nil {
			zap.L().Error("failed to find analysis menu", zap.Error(err))
			return err
		}
		err = db.Model(&normalRole).Association("SystemMenus").Append([]SystemMenu{dashboardMenu, consoleMenu, analysisMenu})
		if err != nil {
			zap.L().Error("failed to associate console and analysis menus with normal role", zap.Error(err))
			return err
		}
		// 创建部门
		departments := []SystemDepartment{
			{Model: gorm.Model{ID: 1}, Name: "管理中心", Sort: 1, Status: 1},
		}
		err = db.Create(&departments).Error
		if err != nil {
			zap.L().Error("failed to create department", zap.Error(err))
			return err
		}
		// 创建用户
		pwd := encryptionPWD(config.AdminPassword)
		users := []SystemUser{
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
		"system_menus", "system_roles", "system_departments", "system_users",
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
