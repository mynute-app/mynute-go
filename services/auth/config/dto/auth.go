package DTO

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Claims represents JWT claims for all user types (admin, tenant, client)
type Claims struct {
	ID        uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Name      string    `json:"name" example:"John"`
	Surname   string    `json:"surname" example:"Doe"`
	Email     string    `json:"email" example:"john.doe@example.com"`
	Phone     string    `json:"phone" example:"+15555555555"`
	Verified  bool      `json:"verified" example:"true"`
	CompanyID uuid.UUID `json:"company_id" example:"00000000-0000-0000-0000-000000000000"`
	Password  string    `json:"password" example:"StrongPswrd1!"`
	Type      string    `json:"type" example:"employee"`
	Roles     []string  `json:"roles" example:"[client]"`
	IsAdmin   bool      `json:"is_admin" example:"false"`
	IsActive  bool      `json:"is_active" example:"true"`
}

type LoginByEmailCode struct {
	Email string `json:"email" example:"john.doe@example.com"`
	Code  string `json:"code" example:"123456"`
}

// AuthRequest represents authorization request for all user types (admin, tenant, client)
type AuthRequest struct {
	Method     string                 `json:"method" validate:"required,oneof=GET POST PUT PATCH DELETE"`
	Path       string                 `json:"path" validate:"required"`
	Resource   map[string]interface{} `json:"resource,omitempty"`
	PathParams map[string]interface{} `json:"path_params,omitempty"`
	Body       map[string]interface{} `json:"body,omitempty"`
	Query      map[string]interface{} `json:"query,omitempty"`
	Headers    map[string]interface{} `json:"headers,omitempty"`
}

// PolicyCreateRequest represents policy creation request for all user types (admin, tenant, client)
type PolicyCreateRequest struct {
	Name        string          `json:"name" validate:"required,min=3,max=100"`
	Description string          `json:"description"`
	Effect      string          `json:"effect" validate:"required,oneof=Allow Deny"`
	EndPointID  string          `json:"end_point_id" validate:"required,uuid"`
	Conditions  json.RawMessage `json:"conditions" validate:"required" swaggertype:"string"`
}

// PolicyUpdateRequest represents policy update request for all user types (admin, tenant, client)
type PolicyUpdateRequest struct {
	Name        *string         `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Description *string         `json:"description,omitempty"`
	Effect      *string         `json:"effect,omitempty" validate:"omitempty,oneof=Allow Deny"`
	EndPointID  *string         `json:"end_point_id,omitempty" validate:"omitempty,uuid"`
	Conditions  json.RawMessage `json:"conditions,omitempty" swaggertype:"string"`
}

// RoleCreateRequest represents role creation request for all user types (admin, tenant, client)
type RoleCreateRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=20"`
	Description string `json:"description" validate:"max=255"`
}

// RoleUpdateRequest represents role update request for all user types (admin, tenant, client)
type RoleUpdateRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=20"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=255"`
}
