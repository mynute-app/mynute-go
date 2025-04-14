package DTO

import "github.com/google/uuid"

type Role struct {
	ID           string    `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Name         string    `json:"name" example:"admin"`
	Description  string    `json:"description" example:"Administrator role"`
	CompanyID    uuid.UUID `json:"company_id" example:"00000000-0000-0000-0000-000000000000"`
	Company      CompanyPopulated
	IsSystemRole bool `json:"is_system_role"`
}

type RolePopulated struct {
	ID           string    `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Name         string    `json:"name" example:"admin"`
	Description  string    `json:"description" example:"Administrator role"`
	CompanyID    uuid.UUID `json:"company_id" example:"00000000-0000-0000-0000-000000000000"`
	IsSystemRole bool      `json:"is_system_role"`
}
