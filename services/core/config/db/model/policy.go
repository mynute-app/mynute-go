package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

// --- Helper function to build condition JSON ---
// /auth service handles validation and execution
func JsonRawMessage(v any) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic("JsonRawMessage failed: " + err.Error())
	}
	return json.RawMessage(data)
}

// --- TenantPolicy Model for Core Service ---
// /auth service handles all authorization logic and validation
type TenantPolicy struct {
	BaseModel
	TenantID    uuid.UUID       `gorm:"type:uuid;uniqueIndex:idx_tenant_policy_name;not null" json:"tenant_id"`
	Name        string          `gorm:"uniqueIndex:idx_tenant_policy_name;not null" json:"name"`
	Description string          `json:"description"`
	Effect      string          `json:"effect"` // "Allow" / "Deny"
	EndPointID  uuid.UUID       `json:"end_point_id"`
	Conditions  json.RawMessage `gorm:"type:jsonb" json:"conditions"`
}

func (TenantPolicy) TableName() string  { return "public.tenant_policies" }
func (TenantPolicy) SchemaType() string { return "public" }

// --- ClientPolicy Model for Core Service ---
// /auth service handles all authorization logic and validation
type ClientPolicy struct {
	BaseModel
	Name        string          `gorm:"uniqueIndex:idx_client_policy_name;not null" json:"name"`
	Description string          `json:"description"`
	Effect      string          `json:"effect"` // "Allow" / "Deny"
	EndPointID  uuid.UUID       `json:"end_point_id"`
	Conditions  json.RawMessage `gorm:"type:jsonb" json:"conditions"`
}

func (ClientPolicy) TableName() string  { return "public.client_policies" }
func (ClientPolicy) SchemaType() string { return "public" }

// --- AdminPolicy Model for Core Service ---
// /auth service handles all authorization logic and validation
type AdminPolicy struct {
	BaseModel
	Name        string          `gorm:"uniqueIndex:idx_admin_policy_name;not null" json:"name"`
	Description string          `json:"description"`
	Effect      string          `json:"effect"` // "Allow" / "Deny"
	EndPointID  uuid.UUID       `json:"end_point_id"`
	Conditions  json.RawMessage `gorm:"type:jsonb" json:"conditions"`
}

func (AdminPolicy) TableName() string  { return "public.admin_policies" }
func (AdminPolicy) SchemaType() string { return "public" }

// --- Condition structures for seed data ---
// These are used only for building seed data, /auth handles evaluation

type ConditionNode struct {
	Description string          `json:"description,omitempty"`
	LogicType   string          `json:"logic_type,omitempty"`
	Children    []ConditionNode `json:"children,omitempty"`
	Leaf        *ConditionLeaf  `json:"leaf,omitempty"`
}

type ConditionLeaf struct {
	Attribute         string          `json:"attribute"`
	Operator          string          `json:"operator"`
	Description       string          `json:"description,omitempty"`
	Value             json.RawMessage `json:"value,omitempty"`
	ResourceAttribute string          `json:"resource_attribute,omitempty"`
}
