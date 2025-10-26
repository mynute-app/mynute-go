package DTO

import "github.com/google/uuid"

type Subdomain struct {
	ID uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"` // Primary key
	Name       string `json:"name" example:"agenda-yourcompany"` // Subdomain name
	CompanyID  uuid.UUID `json:"company_id" example:"00000000-0000-0000-0000-000000000000"` // Foreign key to Company
}
