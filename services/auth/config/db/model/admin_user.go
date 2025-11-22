package model

import (
	mJSON "mynute-go/services/auth/config/db/json"
)

type AdminUser struct {
	BaseModel
	Email    string         `gorm:"type:varchar(255);uniqueIndex:idx_admin_email;not null" json:"email"`
	Password string         `gorm:"type:varchar(255);not null" json:"-"`
	Verified bool           `gorm:"default:false" json:"verified"`
	Meta     mJSON.UserMeta `gorm:"type:jsonb" json:"meta"`
	Roles    []AdminRole    `gorm:"many2many:admin_role_admins;joinForeignKey:admin_id;joinReferences:role_admin_id" json:"roles"`
}
