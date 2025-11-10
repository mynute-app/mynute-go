package DTO

import (
	"encoding/json"
	"mynute-go/services/auth/config/db/model"

	"github.com/google/uuid"
)

// =====================
// ADMIN POLICY DTOs
// =====================

type PaginatedAdminPoliciesResponse struct {
	Data   []model.AdminPolicy `json:"data"`
	Limit  int                 `json:"limit"`
	Offset int                 `json:"offset"`
}

type AdminPolicyBase struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Effect      string          `json:"effect"`
	EndPointID  string          `json:"end_point_id"`
	Conditions  json.RawMessage `json:"conditions"`
}

type AdminPolicyCreateRequest struct {
	Name        string          `json:"name" validate:"required,min=3,max=100"`
	Description string          `json:"description"`
	Effect      string          `json:"effect" validate:"required,oneof=Allow Deny"`
	EndPointID  string          `json:"end_point_id" validate:"required,uuid"`
	Conditions  json.RawMessage `json:"conditions" validate:"required"`
}

type AdminPolicyUpdateRequest struct {
	Name        *string         `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Description *string         `json:"description,omitempty"`
	Effect      *string         `json:"effect,omitempty" validate:"omitempty,oneof=Allow Deny"`
	EndPointID  *string         `json:"end_point_id,omitempty" validate:"omitempty,uuid"`
	Conditions  json.RawMessage `json:"conditions,omitempty"`
}
