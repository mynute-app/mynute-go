package model

import (
	"github.com/google/uuid"
)

type TenantUser struct {
	User
	TenantID uuid.UUID `gorm:"type:uuid;index;not null" json:"tenant_id"` // Optional tenant ID for multi-tenant setups
}
