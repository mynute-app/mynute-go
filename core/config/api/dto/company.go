package DTO

import "github.com/google/uuid"

type CreateCompany struct {
	Name          string `json:"name" example:"Your Company Name"`
	TaxID         string `json:"tax_id" example:"00000000000000"`
	OwnerName     string `json:"owner_name" example:"John"`
	OwnerSurname  string `json:"owner_surname" example:"Clark"`
	OwnerEmail    string `json:"owner_email" example:"john.clark@gmail.com"`
	OwnerPhone    string `json:"owner_phone" example:"+15555555555"`
	OwnerPassword string `json:"owner_password" example:"1SecurePswd!"`
}

// @description	Company Full DTO
// @name			CompanyFullDTO
// @tag.name		company.full.dto
type Company struct {
	ID        uuid.UUID            `json:"id" example:"00000000-0000-0000-0000-000000000000"` // Primary key
	Name      string               `json:"name" example:"Your Company Name"`
	TaxID     string               `json:"tax_id" example:"00000000000000"`
	Employees []*EmployeePopulated `json:"employees"`
	Branches  []*BranchPopulated   `json:"branches"`
	Services  []*ServicePopulated  `json:"services"`
	SectorID  *uuid.UUID           `json:"sector_id"`
	Sector    *Sector              `json:"sector"`
}

// @description	Company DTO Populated
// @name			CompanyPopulatedDTO
// @tag.name		company_populated.dto
type CompanyPopulated struct {
	ID    uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Name  string    `json:"name"`
	TaxID string    `json:"tax_id"`
}
