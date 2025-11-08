package model

import (
	"mynute-go/services/auth/api/lib"

	"gorm.io/gorm"
)

type GeneralRole struct {
	Role
}

func (r *GeneralRole) BeforeCreate(tx *gorm.DB) error {
	if is, err := r.isRoleNameReserved(tx); err != nil {
		return err
	} else if is {
		return lib.Error.Role.NameReserved
	}
	return nil
}

func (r *GeneralRole) BeforeUpdate(tx *gorm.DB) error {
	if nameIsReserved, err := r.isRoleNameReserved(tx); err != nil {
		return err
	} else if nameIsReserved {
		var existing Role
		if err := tx.
			Where("id = ?", r.ID).
			First(&existing).Error; err != nil {
			return err
		}

		// Only block if the name was changed
		if existing.Name != r.Name {
			return lib.Error.Role.NameReserved
		}
	}
	return nil
}

func (r *GeneralRole) isRoleNameReserved(tx *gorm.DB) (bool, error) {
	var count int64
	if err := tx.
		Model(&GeneralRole{}).
		Where("name = ?", r.Name).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
