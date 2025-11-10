package model

import (
	"github.com/google/uuid"
)

type TenantRole struct {
	BaseModel
	TenantID    uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_tenant_role_name;not null" json:"tenant_id"`
	Name        string    `gorm:"type:varchar(50);uniqueIndex:idx_tenant_role_name;not null" json:"name"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
}
