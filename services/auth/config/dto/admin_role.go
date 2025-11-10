package DTO

import (
	"mynute-go/services/auth/config/db/model"

	"github.com/google/uuid"
)

// =====================
// ADMIN ROLE DTOs
// =====================

type PaginatedAdminRolesResponse struct {
	Data   []model.AdminRole `json:"data"`
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
}

type AdminRoleBase struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type AdminRoleCreateRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=20"`
	Description string `json:"description" validate:"max=255"`
}

type AdminRoleUpdateRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=20"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=255"`
}
