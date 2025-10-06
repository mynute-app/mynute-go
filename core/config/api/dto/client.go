package DTO

import (
	"github.com/google/uuid"
)

type LoginClient struct {
	Email    string `json:"email" example:"john.doe@example.com"`
	Password string `json:"password" example:"1SecurePswd!"`
}

type CreateClient struct {
	Name     string `json:"name" example:"John"`
	Surname  string `json:"surname" example:"Doe"`
	Email    string `json:"email" example:"john.doe@example.com"`
	Phone    string `json:"phone" example:"+15555555555"`
	Password string `json:"password" example:"1SecurePswd!"`
}

type Client struct {
	ID       uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Name     string    `json:"name" example:"John"`
	Surname  string    `json:"surname" example:"Doe"`
	Email    string    `json:"email" example:"john.doe@example.com"`
	Phone    string    `json:"phone" example:"+15555555555"`
	Verified bool      `json:"verified" example:"false"`
}

type ClientPopulated struct {
	ID       uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Name     string    `json:"name" example:"John"`
	Surname  string    `json:"surname" example:"Doe"`
	Email    string    `json:"email" example:"john.doe@example.com"`
	Phone    string    `json:"phone" example:"+1-555-555-5555"`
	Verified bool      `json:"verified" example:"false"`
}
