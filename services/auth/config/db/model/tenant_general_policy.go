package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

type TenantGeneralPolicy struct {
	BaseModel
	Name        string          `gorm:"uniqueIndex:idx_tenant_policy_name;not null" json:"name"`
	Description string          `json:"description"`
	Effect      string          `json:"effect"` // "Allow" / "Deny"
	EndPointID  uuid.UUID       `json:"end_point_id"`
	EndPoint    EndPoint        `gorm:"foreignKey:EndPointID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"end_point"`
	Conditions  json.RawMessage `gorm:"type:jsonb" json:"conditions"`
}

// PolicyInterface implementation for TenantGeneralPolicy
func (p *TenantGeneralPolicy) GetID() uuid.UUID               { return p.ID }
func (p *TenantGeneralPolicy) GetName() string                { return p.Name }
func (p *TenantGeneralPolicy) GetDescription() string         { return p.Description }
func (p *TenantGeneralPolicy) GetEffect() string              { return p.Effect }
func (p *TenantGeneralPolicy) GetEndPointID() uuid.UUID       { return p.EndPointID }
func (p *TenantGeneralPolicy) GetConditions() json.RawMessage { return p.Conditions }

func (p *TenantGeneralPolicy) GetConditionsNode() (ConditionNode, error) {
	return GetConditionsNode(p.Name, p.ID.String(), p.Conditions)
}
