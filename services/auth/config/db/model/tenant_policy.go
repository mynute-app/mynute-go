package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

type TenantPolicy struct {
	BaseModel
	TenantID    uuid.UUID       `gorm:"type:uuid;uniqueIndex:idx_tenant_policy_name;not null" json:"tenant_id"`
	Name        string          `gorm:"uniqueIndex:idx_tenant_policy_name;not null" json:"name"`
	Description string          `json:"description"`
	Effect      string          `json:"effect"` // "Allow" / "Deny"
	EndPointID  uuid.UUID       `json:"end_point_id"`
	EndPoint    EndPoint        `gorm:"foreignKey:EndPointID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"end_point"`
	Conditions  json.RawMessage `gorm:"type:jsonb" json:"conditions" swaggertype:"string"`
}

// PolicyInterface implementation for TenantPolicy
func (p *TenantPolicy) GetID() uuid.UUID               { return p.ID }
func (p *TenantPolicy) GetName() string                { return p.Name }
func (p *TenantPolicy) GetDescription() string         { return p.Description }
func (p *TenantPolicy) GetEffect() string              { return p.Effect }
func (p *TenantPolicy) GetEndPointID() uuid.UUID       { return p.EndPointID }
func (p *TenantPolicy) GetConditions() json.RawMessage { return p.Conditions }

func (p *TenantPolicy) GetConditionsNode() (ConditionNode, error) {
	return GetConditionsNode(p.Name, p.ID.String(), p.Conditions)
}
