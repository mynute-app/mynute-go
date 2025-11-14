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

// PolicyInterface implementation for AdminPolicy
func (p *AdminPolicy) GetID() uuid.UUID               { return p.ID }
func (p *AdminPolicy) GetName() string                { return p.Name }
func (p *AdminPolicy) GetDescription() string         { return p.Description }
func (p *AdminPolicy) GetEffect() string              { return p.Effect }
func (p *AdminPolicy) GetEndPointID() uuid.UUID       { return p.EndPointID }
func (p *AdminPolicy) GetConditions() json.RawMessage { return p.Conditions }

func (p *AdminPolicy) GetConditionsNode() (ConditionNode, error) {
	return GetConditionsNode(p.Name, p.ID.String(), p.Conditions)
}
