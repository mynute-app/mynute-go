package DTO

import "github.com/google/uuid"

// AdminClaims represents the JWT claims for admin authentication
type AdminClaims struct {
	ID       uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Name     string    `json:"name" example:"Admin User"`
	Email    string    `json:"email" example:"admin@example.com"`
	Password string    `json:"password" example:"HashedPassword123!"`
	IsAdmin  bool      `json:"is_admin" example:"true"`
	IsActive bool      `json:"is_active" example:"true"`
	Roles    []string  `json:"roles" example:"superadmin,auditor"`
	Type     string    `json:"type" example:"admin"`
}

// AdminLoginRequest represents the request body for admin login
type AdminLoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"admin@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"StrongPassword123!"`
}

// AdminLoginResponse represents the response after successful admin login
type AdminLoginResponse struct {
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	Admin *AdminDetail `json:"admin"`
}

// AdminDetail represents detailed admin information (without password)
type AdminDetail struct {
	ID       uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Name     string    `json:"name" example:"Admin User"`
	Email    string    `json:"email" example:"admin@example.com"`
	IsActive bool      `json:"is_active" example:"true"`
	Roles    []string  `json:"roles" example:"superadmin,auditor"`
}

// AdminCreateRequest represents the request body to create a new admin
type AdminCreateRequest struct {
	Name     string   `json:"name" validate:"required,min=3,max=100" example:"Admin User"`
	Email    string   `json:"email" validate:"required,email" example:"admin@example.com"`
	Password string   `json:"password" validate:"required,min=8" example:"StrongPassword123!"`
	IsActive bool     `json:"is_active" example:"true"`
	Roles    []string `json:"roles" example:"superadmin"`
}

// AdminUpdateRequest represents the request body to update an admin
type AdminUpdateRequest struct {
	Name     *string  `json:"name,omitempty" validate:"omitempty,min=3,max=100" example:"Admin User"`
	Email    *string  `json:"email,omitempty" validate:"omitempty,email" example:"admin@example.com"`
	Password *string  `json:"password,omitempty" validate:"omitempty,min=8" example:"NewPassword123!"`
	IsActive *bool    `json:"is_active,omitempty" example:"true"`
	Roles    []string `json:"roles,omitempty" example:"superadmin,support"`
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
type RoleAdminDetail struct {
	ID          uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Name        string    `json:"name" example:"superadmin"`
	Description string    `json:"description" example:"Full access to all resources"`
	CreatedAt   string    `json:"created_at" example:"2025-01-01T00:00:00Z"`
	UpdatedAt   string    `json:"updated_at" example:"2025-01-01T00:00:00Z"`
}
