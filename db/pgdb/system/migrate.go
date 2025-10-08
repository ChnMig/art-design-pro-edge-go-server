package system

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/config"
)

func migrateTable(db *gorm.DB) error {
	err := db.AutoMigrate(
		&SystemTenant{},
		&SystemDepartment{},
		&SystemRole{},
		&SystemMenu{},
        &SystemMenuAuth{},
        &SystemUser{},
        &SystemUserLoginLog{},
        &SystemTenantMenuScope{},
        &SystemTenantAuthScope{},
	)
	if err != nil {
		zap.L().Error("failed to migrate system model", zap.Error(err))
		return err
	}
	return nil
}

func migrateData(db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		// 检查是否已有租户数据
		var tenantCount int64
		tx.Model(&SystemTenant{}).Count(&tenantCount)
		if tenantCount == 0 {
			// 创建默认租户
			defaultTenant := SystemTenant{
				Model:  gorm.Model{ID: 1},
				Code:   config.DefaultTenantCode,
				Name:   "平台管理",
				Status: StatusEnabled,
			}
			err := tx.Create(&defaultTenant).Error
			if err != nil {
				zap.L().Error("failed to create default tenant", zap.Error(err))
				return err
			}
		}

		// 检查是否已有数据，如果有则跳过初始化
		var count int64
		tx.Model(&SystemMenu{}).Count(&count)
		if count > 0 {
			zap.L().Info("menu data already exists, skipping initial data creation")
			return nil
		}

        // 创建菜单（适配当前前端路由结构）
        menus := []SystemMenu{
            // 顶部一级：仪表盘
            {Model: gorm.Model{ID: 1}, Path: "/dashboard", Name: "Dashboard", Component: "/index/index", Title: "仪表盘", Icon: "&#xe721;", KeepAlive: 2, Status: StatusEnabled, Level: 1, ParentID: 0, Sort: 1},
            // 顶部一级：平台管理
            {Model: gorm.Model{ID: 2}, Path: "/platform", Name: "Platform", Component: "/index/index", Title: "平台管理", Icon: "&#xe72b;", KeepAlive: 2, Status: StatusEnabled, Level: 1, ParentID: 0, Sort: 2},
            // 平台管理二级（去除“菜单范围”页面）
            {Model: gorm.Model{ID: 3}, Path: "tenant", Name: "PlatformTenant", Component: "/platform/tenant/index", Title: "租户管理", KeepAlive: 2, Status: StatusEnabled, Level: 2, ParentID: 2, Sort: 1},
            {Model: gorm.Model{ID: 4}, Path: "menu", Name: "PlatformMenu", Component: "/platform/menu/index", Title: "菜单管理", KeepAlive: 2, Status: StatusEnabled, Level: 2, ParentID: 2, Sort: 2},
            // 顶部一级：系统管理
            {Model: gorm.Model{ID: 5}, Path: "/system", Name: "System", Component: "/index/index", Title: "系统管理", Icon: "&#xe72b;", KeepAlive: 2, Status: StatusEnabled, Level: 1, ParentID: 0, Sort: 3},
            // 系统管理二级（不包含“菜单管理”，菜单仅平台管理可操作）
            {Model: gorm.Model{ID: 6}, Path: "role", Name: "TenantRole", Component: "/system/role/index", Title: "角色管理", KeepAlive: 2, Status: StatusEnabled, Level: 2, ParentID: 5, Sort: 1},
            {Model: gorm.Model{ID: 7}, Path: "department", Name: "SystemDepartment", Component: "/system/department/index", Title: "部门管理", KeepAlive: 2, Status: StatusEnabled, Level: 2, ParentID: 5, Sort: 2},
            {Model: gorm.Model{ID: 8}, Path: "user", Name: "SystemUser", Component: "/system/user/index", Title: "用户管理", KeepAlive: 2, Status: StatusEnabled, Level: 2, ParentID: 5, Sort: 3},
            // 仪表盘子页
            {Model: gorm.Model{ID: 9}, Path: "console", Name: "DashboardConsole", Component: "/dashboard/console/index", Title: "工作台", Icon: "", KeepAlive: 2, Status: StatusEnabled, Level: 2, ParentID: 1, Sort: 1},
            {Model: gorm.Model{ID: 10}, Path: "analysis", Name: "DashboardAnalysis", Component: "/dashboard/analysis/index", Title: "分析页", Icon: "", KeepAlive: 2, Status: StatusEnabled, Level: 2, ParentID: 1, Sort: 2},
            // 隐藏页
            {Model: gorm.Model{ID: 11}, Path: "/private", Name: "Private", Component: "/index/index", Title: "隐藏页面", Icon: "", KeepAlive: 2, Status: StatusEnabled, Level: 1, ParentID: 0, Sort: 99, IsHide: 1},
        }
		err := tx.Create(&menus).Error
		if err != nil {
			zap.L().Error("failed to create menu", zap.Error(err))
			return err
		}

		// 检查是否已有角色数据
		tx.Model(&SystemRole{}).Count(&count)
		if count > 0 {
			zap.L().Info("role data already exists, skipping role creation")
			return nil
		}

        // 创建角色（默认租户）仅保留“超级管理员”
        roles := []SystemRole{
            {Model: gorm.Model{ID: 1}, TenantID: 1, Name: "超级管理员", Desc: "拥有所有权限", Status: StatusEnabled},
        }
        err = tx.Create(&roles).Error
        if err != nil {
            zap.L().Error("failed to create role", zap.Error(err))
            return err
        }

		// 为角色分配菜单权限
		// 超级管理员拥有所有菜单权限
		adminRole := SystemRole{}
		err = tx.First(&adminRole, 1).Error
		if err != nil {
			zap.L().Error("failed to find admin role", zap.Error(err))
			return err
		}
		// 为超级管理员分配所有菜单
		var allMenus []SystemMenu
		err = tx.Find(&allMenus).Error
		if err != nil {
			zap.L().Error("failed to find menus", zap.Error(err))
			return err
		}
		err = tx.Model(&adminRole).Association("SystemMenus").Append(&allMenus)
		if err != nil {
			zap.L().Error("failed to associate menus with admin role", zap.Error(err))
			return err
		}
        // 移除默认“普通用户”角色及其权限分配，留给超级管理员自定义创建
        // 默认租户可访问的菜单范围：授权所有页面
        allMenuScopes := make([]SystemTenantMenuScope, 0, len(allMenus))
        for _, m := range allMenus {
            allMenuScopes = append(allMenuScopes, SystemTenantMenuScope{TenantID: 1, MenuID: m.ID})
        }
        if err := tx.Create(&allMenuScopes).Error; err != nil {
            zap.L().Error("failed to create default tenant full menu scope", zap.Error(err))
            return err
        }

        // 检查是否已有部门数据
		tx.Model(&SystemDepartment{}).Count(&count)
		if count > 0 {
			zap.L().Info("department data already exists, skipping department creation")
			return nil
		}

		// 创建部门（默认租户）
		departments := []SystemDepartment{
			{Model: gorm.Model{ID: 1}, TenantID: 1, Name: "管理中心", Sort: 1, Status: StatusEnabled},
		}
		err = tx.Create(&departments).Error
		if err != nil {
			zap.L().Error("failed to create department", zap.Error(err))
			return err
		}

		// 检查是否已有用户数据
		tx.Model(&SystemUser{}).Count(&count)
		if count > 0 {
			zap.L().Info("user data already exists, skipping user creation")
			return nil
		}

		// 创建用户
		pwd, hashErr := HashPassword(config.AdminPassword)
		if hashErr != nil {
			zap.L().Error("failed to hash admin password", zap.Error(hashErr))
			return hashErr
		}
		users := []SystemUser{
			{Model: gorm.Model{ID: 1}, TenantID: 1, DepartmentID: 1, RoleID: 1, Name: "超级管理员", Username: "admin", Account: "admin", Password: pwd, Status: StatusEnabled, Gender: 1},
		}
		err = tx.Create(&users).Error
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
        "system_tenants", "system_menus", "system_roles", "system_departments", "system_users",
        "system_menu_auths", "system_tenant_menu_scopes", "system_tenant_auth_scopes",
    }

	for _, table := range tables {
		seqName := table + "_id_seq"
		query := fmt.Sprintf("SELECT setval('%s', (SELECT COALESCE(MAX(id), 1) FROM %s));", seqName, table)
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
