package DTO

import (
	"encoding/json"
	"mynute-go/services/auth/config/db/model"

	"github.com/google/uuid"
)

// =====================
// TENANT POLICY DTOs
// =====================

type PaginatedTenantPoliciesResponse struct {
	Data   []model.TenantPolicy `json:"data"`
	Limit  int                  `json:"limit"`
	Offset int                  `json:"offset"`
}

type TenantPolicyBase struct {
	ID          uuid.UUID       `json:"id"`
	TenantID    uuid.UUID       `json:"tenant_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Effect      string          `json:"effect"`
	EndPointID  string          `json:"end_point_id"`
	Conditions  json.RawMessage `json:"conditions"`
}

type TenantPolicyCreateRequest struct {
	Name        string          `json:"name" validate:"required,min=3,max=100"`
	Description string          `json:"description"`
	Effect      string          `json:"effect" validate:"required,oneof=Allow Deny"`
	EndPointID  string          `json:"end_point_id" validate:"required,uuid"`
	Conditions  json.RawMessage `json:"conditions" validate:"required"`
}

type TenantPolicyUpdateRequest struct {
	Name        *string         `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Description *string         `json:"description,omitempty"`
	Effect      *string         `json:"effect,omitempty" validate:"omitempty,oneof=Allow Deny"`
	EndPointID  *string         `json:"end_point_id,omitempty" validate:"omitempty,uuid"`
	Conditions  json.RawMessage `json:"conditions,omitempty"`
}
