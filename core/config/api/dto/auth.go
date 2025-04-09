package DTO

import "github.com/google/uuid"

type Claims struct {
	ID        uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Name      string    `json:"name" example:"John"`
	Surname   string    `json:"surname" example:"Doe"`
	Email     string    `json:"email" example:"john.doe@example.com"`
	Phone     string    `json:"phone" example:"+15555555555"`
	Tags      []string  `json:"tags" example:"[\"tag1\", \"tag2\"]"`
	Verified  bool      `json:"verified" example:"true"`
	CompanyID uuid.UUID `json:"company_id" example:"00000000-0000-0000-0000-000000000000"`
}
