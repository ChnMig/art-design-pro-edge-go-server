package system

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/config"
)

func migrateTable(db *gorm.DB) error {
	if err := cleanupLegacyScopeRows(db); err != nil {
		return err
	}
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
	return cleanupRelationalOrphans(db)
}

func cleanupRelationalOrphans(db *gorm.DB) error {
	queries := []string{
		"DELETE FROM system_menu_auths a WHERE NOT EXISTS (SELECT 1 FROM system_menus m WHERE m.id = a.menu_id AND m.deleted_at IS NULL)",
		"DELETE FROM system_roles__system_menus l WHERE NOT EXISTS (SELECT 1 FROM system_roles r WHERE r.id = l.system_role_id AND r.deleted_at IS NULL) OR NOT EXISTS (SELECT 1 FROM system_menus m WHERE m.id = l.system_menu_id AND m.deleted_at IS NULL)",
		"DELETE FROM system_roles__system_auths l WHERE NOT EXISTS (SELECT 1 FROM system_roles r WHERE r.id = l.system_role_id AND r.deleted_at IS NULL) OR NOT EXISTS (SELECT 1 FROM system_menu_auths a WHERE a.id = l.system_menu_auth_id AND a.deleted_at IS NULL)",
		"DELETE FROM system_tenant_menu_scopes s WHERE NOT EXISTS (SELECT 1 FROM system_tenants t WHERE t.id = s.tenant_id AND t.deleted_at IS NULL) OR NOT EXISTS (SELECT 1 FROM system_menus m WHERE m.id = s.menu_id AND m.deleted_at IS NULL)",
		"DELETE FROM system_tenant_auth_scopes s WHERE NOT EXISTS (SELECT 1 FROM system_tenants t WHERE t.id = s.tenant_id AND t.deleted_at IS NULL) OR NOT EXISTS (SELECT 1 FROM system_menu_auths a WHERE a.id = s.auth_id AND a.deleted_at IS NULL)",
	}
	return db.Transaction(func(tx *gorm.DB) error {
		for _, query := range queries {
			if err := tx.Exec(query).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func cleanupLegacyScopeRows(db *gorm.DB) error {
	cleanupQueries := map[interface{}][]string{
		&SystemTenantMenuScope{}: {
			"DELETE FROM system_tenant_menu_scopes WHERE deleted_at IS NOT NULL",
			"DELETE FROM system_tenant_menu_scopes older USING system_tenant_menu_scopes newer WHERE older.id < newer.id AND older.tenant_id = newer.tenant_id AND older.menu_id = newer.menu_id",
		},
		&SystemTenantAuthScope{}: {
			"DELETE FROM system_tenant_auth_scopes WHERE deleted_at IS NOT NULL",
			"DELETE FROM system_tenant_auth_scopes older USING system_tenant_auth_scopes newer WHERE older.id < newer.id AND older.tenant_id = newer.tenant_id AND older.auth_id = newer.auth_id",
		},
	}
	for model, queries := range cleanupQueries {
		if !db.Migrator().HasTable(model) {
			continue
		}
		for _, query := range queries {
			if err := db.Exec(query).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func migrateData(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := ensureDefaultTenant(tx); err != nil {
			return err
		}

		menus := initialMenus()
		for i := range menus {
			if err := restoreOrCreateMenu(tx, &menus[i]); err != nil {
				return err
			}
		}

		adminRole := SystemRole{}
		if err := tx.Unscoped().Where("id = ?", DefaultTenantID).First(&adminRole).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			adminRole = SystemRole{
				Model:         gorm.Model{ID: DefaultTenantID},
				TenantID:      DefaultTenantID,
				Name:          "超级管理员",
				Desc:          "拥有所有权限",
				Status:        StatusEnabled,
				IsTenantAdmin: true,
			}
			if err := tx.Create(&adminRole).Error; err != nil {
				return err
			}
		} else {
			updates := map[string]interface{}{
				"tenant_id": DefaultTenantID, "status": StatusEnabled, "is_tenant_admin": true, "deleted_at": nil,
			}
			if adminRole.Name == "" {
				updates["name"] = "超级管理员"
			}
			if err := tx.Unscoped().Model(&adminRole).Updates(updates).Error; err != nil {
				return err
			}
		}

		department := SystemDepartment{}
		if err := tx.Unscoped().Where("id = ?", DefaultTenantID).First(&department).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			department = SystemDepartment{Model: gorm.Model{ID: DefaultTenantID}, TenantID: DefaultTenantID, Name: "管理中心", Sort: 1, Status: StatusEnabled}
			if err := tx.Create(&department).Error; err != nil {
				return err
			}
		} else {
			updates := map[string]interface{}{
				"tenant_id": DefaultTenantID, "status": StatusEnabled, "deleted_at": nil,
			}
			if department.Name == "" {
				updates["name"] = "管理中心"
			}
			if err := tx.Unscoped().Model(&department).Updates(updates).Error; err != nil {
				return err
			}
		}

		adminUser := SystemUser{}
		if err := tx.Unscoped().Where("id = ?", DefaultTenantID).First(&adminUser).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			password, hashErr := HashPassword(config.AdminPassword)
			if hashErr != nil {
				return hashErr
			}
			adminUser = SystemUser{
				Model: gorm.Model{ID: DefaultTenantID}, TenantID: DefaultTenantID, DepartmentID: department.ID,
				RoleID: adminRole.ID, Name: "超级管理员", Username: "admin", Account: "admin", Password: password,
				Status: StatusEnabled, Gender: 1,
			}
			if err := tx.Create(&adminUser).Error; err != nil {
				return err
			}
		} else {
			updates := map[string]interface{}{
				"tenant_id": DefaultTenantID, "department_id": department.ID, "role_id": adminRole.ID,
				"status": StatusEnabled, "deleted_at": nil,
			}
			if adminUser.Account == "" {
				updates["account"] = "admin"
			}
			if adminUser.Username == "" {
				updates["username"] = "admin"
			}
			if adminUser.Name == "" {
				updates["name"] = "超级管理员"
			}
			if adminUser.Password == "" {
				password, hashErr := HashPassword(config.AdminPassword)
				if hashErr != nil {
					return hashErr
				}
				updates["password"] = password
			}
			if err := tx.Unscoped().Model(&adminUser).Updates(updates).Error; err != nil {
				return err
			}
		}

		var allMenus []SystemMenu
		if err := tx.Find(&allMenus).Error; err != nil {
			return err
		}
		if err := tx.Model(&adminRole).Association("SystemMenus").Replace(&allMenus); err != nil {
			return err
		}
		var allAuths []SystemMenuAuth
		if err := tx.Find(&allAuths).Error; err != nil {
			return err
		}
		if err := tx.Model(&adminRole).Association("SystemMenuAuths").Replace(&allAuths); err != nil {
			return err
		}
		if err := ensureDefaultTenantScopes(tx, allMenus, allAuths); err != nil {
			return err
		}
		return nil
	})
}

func ensureDefaultTenant(tx *gorm.DB) error {
	tenant := SystemTenant{}
	err := tx.Unscoped().Where("id = ?", DefaultTenantID).First(&tenant).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		tenant = SystemTenant{Model: gorm.Model{ID: DefaultTenantID}, Code: config.DefaultTenantCode, Name: "平台管理", Status: StatusEnabled}
		return tx.Create(&tenant).Error
	}
	if err != nil {
		return err
	}
	updates := map[string]interface{}{"status": StatusEnabled, "deleted_at": nil}
	if tenant.Code == "" {
		updates["code"] = config.DefaultTenantCode
	}
	if tenant.Name == "" {
		updates["name"] = "平台管理"
	}
	return tx.Unscoped().Model(&tenant).Updates(updates).Error
}

func initialMenus() []SystemMenu {
	return []SystemMenu{
		{Model: gorm.Model{ID: 1}, Path: "/dashboard", Name: "Dashboard", Component: "/index/index", Title: "仪表盘", Icon: "ri-dashboard-line", KeepAlive: 2, Status: StatusEnabled, Level: 1, Sort: 1},
		{Model: gorm.Model{ID: 2}, Path: "/platform", Name: "Platform", Component: "/index/index", Title: "平台管理", Icon: "ri-settings-line", KeepAlive: 2, Status: StatusEnabled, Level: 1, Sort: 2},
		{Model: gorm.Model{ID: 3}, Path: "tenant", Name: "PlatformTenant", Component: "/platform/tenant/index", Title: "租户管理", KeepAlive: 2, Status: StatusEnabled, Level: 2, ParentID: 2, Sort: 1},
		{Model: gorm.Model{ID: 4}, Path: "menu", Name: "PlatformMenu", Component: "/platform/menu/index", Title: "菜单管理", KeepAlive: 2, Status: StatusEnabled, Level: 2, ParentID: 2, Sort: 2},
		{Model: gorm.Model{ID: 5}, Path: "/system", Name: "System", Component: "/index/index", Title: "系统管理", Icon: "ri-user-line", KeepAlive: 2, Status: StatusEnabled, Level: 1, Sort: 3},
		{Model: gorm.Model{ID: 6}, Path: "role", Name: "TenantRole", Component: "/system/role/index", Title: "角色管理", KeepAlive: 2, Status: StatusEnabled, Level: 2, ParentID: 5, Sort: 1},
		{Model: gorm.Model{ID: 7}, Path: "department", Name: "SystemDepartment", Component: "/system/department/index", Title: "部门管理", KeepAlive: 2, Status: StatusEnabled, Level: 2, ParentID: 5, Sort: 2},
		{Model: gorm.Model{ID: 8}, Path: "user", Name: "SystemUser", Component: "/system/user/index", Title: "用户管理", KeepAlive: 2, Status: StatusEnabled, Level: 2, ParentID: 5, Sort: 3},
		{Model: gorm.Model{ID: 9}, Path: "console", Name: "DashboardConsole", Component: "/dashboard/console/index", Title: "工作台", KeepAlive: 2, Status: StatusEnabled, Level: 2, ParentID: 1, Sort: 1},
		{Model: gorm.Model{ID: 10}, Path: "analysis", Name: "DashboardAnalysis", Component: "/dashboard/analysis/index", Title: "分析页", KeepAlive: 2, Status: StatusEnabled, Level: 2, ParentID: 1, Sort: 2},
		{Model: gorm.Model{ID: 11}, Path: "/private", Name: "Private", Component: "/index/index", Title: "隐藏页面", KeepAlive: 2, Status: StatusEnabled, Level: 1, Sort: 99, IsHide: 1},
	}
}

func restoreOrCreateMenu(tx *gorm.DB, seed *SystemMenu) error {
	existing := SystemMenu{}
	err := tx.Unscoped().Where("id = ?", seed.ID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return tx.Create(seed).Error
	}
	if err != nil {
		return err
	}
	if existing.DeletedAt.Valid {
		return tx.Unscoped().Model(&existing).Update("deleted_at", nil).Error
	}
	return nil
}

func ensureDefaultTenantScopes(tx *gorm.DB, menus []SystemMenu, auths []SystemMenuAuth) error {
	for _, menu := range menus {
		scope := SystemTenantMenuScope{}
		err := tx.Unscoped().Where("tenant_id = ? AND menu_id = ?", DefaultTenantID, menu.ID).First(&scope).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := tx.Create(&SystemTenantMenuScope{TenantID: DefaultTenantID, MenuID: menu.ID}).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else if scope.DeletedAt.Valid {
			if err := tx.Unscoped().Model(&scope).Update("deleted_at", nil).Error; err != nil {
				return err
			}
		}
	}
	for _, auth := range auths {
		scope := SystemTenantAuthScope{}
		err := tx.Unscoped().Where("tenant_id = ? AND auth_id = ?", DefaultTenantID, auth.ID).First(&scope).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := tx.Create(&SystemTenantAuthScope{TenantID: DefaultTenantID, AuthID: auth.ID}).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else if scope.DeletedAt.Valid {
			if err := tx.Unscoped().Model(&scope).Update("deleted_at", nil).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func resetSequences(db *gorm.DB) error {
	tables := []string{
		"system_tenants", "system_menus", "system_roles", "system_departments", "system_users",
		"system_menu_auths", "system_tenant_menu_scopes", "system_tenant_auth_scopes",
	}

	for _, table := range tables {
		seqName := table + "_id_seq"
		query := fmt.Sprintf("SELECT setval(pg_get_serial_sequence('%s', 'id'), COALESCE(MAX(id), 1), MAX(id) IS NOT NULL) FROM %s;", table, table)
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
