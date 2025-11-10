package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

type AdminPolicy struct {
	BaseModel
	Name        string          `gorm:"uniqueIndex:idx_admin_policy_name;not null" json:"name"`
	Description string          `json:"description"`
	Effect      string          `json:"effect"` // "Allow" / "Deny"
	EndPointID  uuid.UUID       `json:"end_point_id"`
	EndPoint    EndPoint        `gorm:"foreignKey:EndPointID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"end_point"`
	Conditions  json.RawMessage `gorm:"type:jsonb" json:"conditions"`
}
