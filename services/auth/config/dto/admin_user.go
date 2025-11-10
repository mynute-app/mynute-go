package DTO

import (
	"mynute-go/services/auth/config/db/model"

	"github.com/google/uuid"
)

// AdminClaims represents JWT claims for admin authentication
type AdminUserClaims struct {
	ID       uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Name     string    `json:"name" example:"Admin User"`
	Email    string    `json:"email" example:"admin@example.com"`
	IsAdmin  bool      `json:"is_admin" example:"true"`
	IsActive bool      `json:"is_active" example:"true"`
	Type     string    `json:"type" example:"admin"`
	Roles    []string  `json:"roles" example:"[superadmin]"`
}

// AdminLoginRequest represents the request body for admin login
type AdminUserLoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"admin@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"StrongPassword123!"`
}

// AdminDetail represents detailed admin information (without password)
type AdminUser struct {
	ID       uuid.UUID         `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Name     string            `json:"name" example:"John"`
	Surname  string            `json:"surname" example:"Doe"`
	Email    string            `json:"email" example:"admin@example.com"`
	IsActive bool              `json:"is_active" example:"true"`
	Roles    []model.AdminRole `json:"roles"`
}

type AdminUserList struct {
	Admins []AdminUser `json:"admins"`
}

// AdminCreateRequest represents the request body to create a new admin
type AdminUserCreateRequest struct {
	Name     string   `json:"name" validate:"required,min=3,max=100" example:"Admin User"`
	Surname  string   `json:"surname,omitempty" validate:"omitempty,min=3,max=100" example:"Doe"`
	Email    string   `json:"email" validate:"required,email" example:"admin@example.com"`
	Password string   `json:"password" validate:"required,min=8" example:"StrongPassword123!"`
	IsActive bool     `json:"is_active" example:"true"`
	Roles    []string `json:"roles" example:"[superadmin]"`
}

// AdminUpdateRequest represents the request body to update an admin
type AdminUserUpdateRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=3,max=100" example:"Admin User"`
	Surname  *string `json:"surname,omitempty" validate:"omitempty,min=3,max=100" example:"Doe"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email" example:"admin@example.com"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=8" example:"NewPassword123!"`
	IsActive *bool   `json:"is_active,omitempty" example:"true"`
}
