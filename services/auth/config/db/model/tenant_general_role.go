package model

import (
	"mynute-go/services/auth/api/lib"

	"gorm.io/gorm"
)

type TenantGeneralRole struct {
	Role
}

func (r *TenantGeneralRole) BeforeCreate(tx *gorm.DB) error {
	if is, err := r.isRoleNameReserved(tx); err != nil {
		return err
	} else if is {
		return lib.Error.Role.NameReserved
	}
	return nil
}

func (r *TenantGeneralRole) BeforeUpdate(tx *gorm.DB) error {
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

func (r *TenantGeneralRole) isRoleNameReserved(tx *gorm.DB) (bool, error) {
	var count int64
	if err := tx.
		Model(&Role{}).
		Where("name = ?", r.Name).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
