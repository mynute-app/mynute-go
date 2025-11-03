package DTO

import "github.com/google/uuid"

// AdminClaims represents JWT claims for admin authentication
type AdminClaims struct {
	ID       uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Name     string    `json:"name" example:"Admin User"`
	Email    string    `json:"email" example:"admin@example.com"`
	IsAdmin  bool      `json:"is_admin" example:"true"`
	IsActive bool      `json:"is_active" example:"true"`
	Type     string    `json:"type" example:"admin"`
	Roles    []string  `json:"roles" example:"[superadmin]"`
}

// AdminLoginRequest represents the request body for admin login
type AdminLoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"admin@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"StrongPassword123!"`
}

// AdminDetail represents detailed admin information (without password)
type Admin struct {
	ID       uuid.UUID   `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Name     string      `json:"name" example:"John"`
	Surname  string      `json:"surname" example:"Doe"`
	Email    string      `json:"email" example:"admin@example.com"`
	IsActive bool        `json:"is_active" example:"true"`
	Roles    []AdminRole `json:"roles"`
}

type AdminList struct {
	Admins []Admin `json:"admins"`
	Total  int     `json:"total" example:"1"`
}

// AdminCreateRequest represents the request body to create a new admin
type AdminCreateRequest struct {
	Name     string   `json:"name" validate:"required,min=3,max=100" example:"Admin User"`
	Surname  string  `json:"surname,omitempty" validate:"omitempty,min=3,max=100" example:"Doe"`
	Email    string   `json:"email" validate:"required,email" example:"admin@example.com"`
	Password string   `json:"password" validate:"required,min=8" example:"StrongPassword123!"`
	IsActive bool     `json:"is_active" example:"true"`
	Roles    []string `json:"roles" example:"[superadmin]"`
}

// AdminUpdateRequest represents the request body to update an admin
type AdminUpdateRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=3,max=100" example:"Admin User"`
	Surname  *string `json:"surname,omitempty" validate:"omitempty,min=3,max=100" example:"Doe"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email" example:"admin@example.com"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=8" example:"NewPassword123!"`
	IsActive *bool   `json:"is_active,omitempty" example:"true"`
}

// RoleAdminCreateRequest represents the request body to create a new admin role
type RoleAdminCreateRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100" example:"support"`
	Description string `json:"description" example:"Customer support role with limited access"`
}

// RoleAdminUpdateRequest represents the request body to update an admin role
type RoleAdminUpdateRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=100" example:"support"`
	Description *string `json:"description,omitempty" example:"Updated description"`
}

// RoleAdminDetail represents detailed admin role information
type AdminRole struct {
	ID          uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Name        string    `json:"name" example:"superadmin"`
	Description string    `json:"description" example:"Full access to all resources"`
	CreatedAt   string    `json:"created_at" example:"2025-01-01T00:00:00Z"`
	UpdatedAt   string    `json:"updated_at" example:"2025-01-01T00:00:00Z"`
}

