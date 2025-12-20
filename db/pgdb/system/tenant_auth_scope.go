package system

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/db/pgdb"
)

// GetTenantAuthScopeIDs 获取租户已授权的按钮权限ID集合
func GetTenantAuthScopeIDs(tenantID uint) ([]uint, error) {
	if tenantID == 0 {
		return nil, nil
	}
	var scopes []SystemTenantAuthScope
	if err := pgdb.GetClient().Where("tenant_id = ?", tenantID).Find(&scopes).Error; err != nil {
		zap.L().Error("failed to get tenant auth scope", zap.Uint("tenantID", tenantID), zap.Error(err))
		return nil, err
	}
	ids := make([]uint, 0, len(scopes))
	for _, s := range scopes {
		ids = append(ids, s.AuthID)
	}
	return ids, nil
}

// SaveTenantAuthScope 保存租户的按钮权限范围（全量覆盖）
func SaveTenantAuthScope(tenantID uint, authIDs []uint) error {
	if tenantID == 0 {
		return fmt.Errorf("tenant id is required")
	}
	return pgdb.GetClient().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("tenant_id = ?", tenantID).Delete(&SystemTenantAuthScope{}).Error; err != nil {
			zap.L().Error("failed to clear tenant auth scope", zap.Uint("tenantID", tenantID), zap.Error(err))
			return err
		}
		if len(authIDs) == 0 {
			return nil
		}
		// 校验权限ID有效性
		var count int64
		if err := tx.Model(&SystemMenuAuth{}).Where("id IN ?", authIDs).Count(&count).Error; err != nil {
			zap.L().Error("failed to validate auth ids", zap.Error(err))
			return err
		}
		if count != int64(len(authIDs)) {
			return fmt.Errorf("invalid auth ids provided")
		}
		items := make([]SystemTenantAuthScope, len(authIDs))
		for i, id := range authIDs {
			items[i] = SystemTenantAuthScope{TenantID: tenantID, AuthID: id}
		}
		if err := tx.Create(&items).Error; err != nil {
			zap.L().Error("failed to create tenant auth scope", zap.Uint("tenantID", tenantID), zap.Error(err))
			return err
		}
		return nil
	})
}
