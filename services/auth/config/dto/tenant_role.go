package DTO

import (
	"mynute-go/services/auth/config/db/model"

	"github.com/google/uuid"
)

// =====================
// TENANT ROLE DTOs
// =====================

type PaginatedTenantRolesResponse struct {
	Data   []model.TenantRole `json:"data"`
	Limit  int                `json:"limit"`
	Offset int                `json:"offset"`
}

type TenantRoleBase struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type TenantRoleCreateRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=20"`
	Description string `json:"description" validate:"max=255"`
}

type TenantRoleUpdateRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=20"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=255"`
}
