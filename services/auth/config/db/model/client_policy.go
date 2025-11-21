package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

type ClientPolicy struct {
	BaseModel
	Name        string          `gorm:"uniqueIndex:idx_client_policy_name;not null" json:"name"`
	Description string          `json:"description"`
	Effect      string          `json:"effect"` // "Allow" / "Deny"
	EndPointID  uuid.UUID       `json:"end_point_id"`
	EndPoint    EndPoint        `gorm:"foreignKey:EndPointID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"end_point"`
	Conditions  json.RawMessage `gorm:"type:jsonb" json:"conditions" swaggertype:"string"`
}

// PolicyInterface implementation for ClientPolicy
func (p *ClientPolicy) GetID() uuid.UUID               { return p.ID }
func (p *ClientPolicy) GetName() string                { return p.Name }
func (p *ClientPolicy) GetDescription() string         { return p.Description }
func (p *ClientPolicy) GetEffect() string              { return p.Effect }
func (p *ClientPolicy) GetEndPointID() uuid.UUID       { return p.EndPointID }
func (p *ClientPolicy) GetConditions() json.RawMessage { return p.Conditions }

func (p *ClientPolicy) GetConditionsNode() (ConditionNode, error) {
	return GetConditionsNode(p.Name, p.ID.String(), p.Conditions)
}
