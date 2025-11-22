package model

import (
	mJSON "mynute-go/services/auth/config/db/json"

	"github.com/google/uuid"
)

type TenantUser struct {
	BaseModel
	TenantID uuid.UUID      `gorm:"type:uuid;uniqueIndex:idx_tenant_email;not null" json:"tenant_id"`
	Email    string         `gorm:"type:varchar(255);uniqueIndex:idx_tenant_email;not null" json:"email"`
	Password string         `gorm:"type:varchar(255);not null" json:"-"`
	Verified bool           `gorm:"default:false" json:"verified"`
	Meta     mJSON.UserMeta `gorm:"type:jsonb" json:"meta"`
	Roles    []TenantRole   `gorm:"many2many:tenant_user_roles;" json:"roles"`
}
