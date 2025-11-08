package model

import (
	"github.com/google/uuid"
)

type TenantPolicy struct {
	Policy
	TenantID    uuid.UUID       `gorm:"type:uuid;index" json:"tenant_id"` // Optional tenant ID for multi-tenant setups
}