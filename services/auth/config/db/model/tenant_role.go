package model

import (
	"mynute-go/services/auth/api/lib"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type TenantRole struct {
	Role
	TenantID uuid.UUID `gorm:"type:uuid;index;not null" json:"tenant_id"` // Optional tenant ID for multi-tenant setups
}

func (r *TenantRole) BeforeCreate(tx *gorm.DB) error {
	if is, err := r.isRoleNameReserved(tx); err != nil {
		return err
	} else if is {
		return lib.Error.Role.NameReserved
	}
	return nil
}

func (r *TenantRole) BeforeUpdate(tx *gorm.DB) error {
	if nameIsReserved, err := r.isRoleNameReserved(tx); err != nil {
		return err
	} else if nameIsReserved {
		var existing Role
		if err := tx.
			Where("id = ?", r.ID).
			First(&existing).Error; err != nil {
			return err
		}

		// Só bloqueia se o nome foi alterado
		if existing.Name != r.Name {
			return lib.Error.Role.NameReserved
		}
	}
	return nil
}

func (r *TenantRole) isRoleNameReserved(tx *gorm.DB) (bool, error) {
	var count int64
	if err := tx.
		Model(&Role{}).
		Where("name = ? AND tenant_id = ?", r.Name, r.TenantID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
