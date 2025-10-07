package system

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-server/db/pgdb"
)

func GetTenantMenuScopeIDs(tenantID uint) ([]uint, error) {
	if tenantID == 0 {
		return nil, nil
	}
	var scopes []SystemTenantMenuScope
	if err := pgdb.GetClient().Where("tenant_id = ?", tenantID).Find(&scopes).Error; err != nil {
		zap.L().Error("failed to get tenant menu scope", zap.Uint("tenantID", tenantID), zap.Error(err))
		return nil, err
	}
	menuIDs := make([]uint, 0, len(scopes))
	for _, scope := range scopes {
		menuIDs = append(menuIDs, scope.MenuID)
	}
	return menuIDs, nil
}

func SaveTenantMenuScope(tenantID uint, menuIDs []uint) error {
	if tenantID == 0 {
		return fmt.Errorf("tenant id is required")
	}
	return pgdb.GetClient().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("tenant_id = ?", tenantID).Delete(&SystemTenantMenuScope{}).Error; err != nil {
			zap.L().Error("failed to clear tenant menu scope", zap.Uint("tenantID", tenantID), zap.Error(err))
			return err
		}
		if len(menuIDs) == 0 {
			return nil
		}
		var count int64
		if err := tx.Model(&SystemMenu{}).Where("id IN ?", menuIDs).Count(&count).Error; err != nil {
			zap.L().Error("failed to validate menu ids", zap.Error(err))
			return err
		}
		if count != int64(len(menuIDs)) {
			return fmt.Errorf("invalid menu ids provided")
		}
		items := make([]SystemTenantMenuScope, len(menuIDs))
		for i, id := range menuIDs {
			items[i] = SystemTenantMenuScope{
				TenantID: tenantID,
				MenuID:   id,
			}
		}
		if err := tx.Create(&items).Error; err != nil {
			zap.L().Error("failed to create tenant menu scope", zap.Uint("tenantID", tenantID), zap.Error(err))
			return err
		}
		return nil
	})
}
