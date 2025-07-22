package DTO

import (
	mJSON "mynute-go/core/config/db/model/json"

	"github.com/google/uuid"
)

type CreateCompany struct {
	LegalName      string `json:"name" example:"Your Company Legal Name"`
	TradeName      string `json:"trading_name" example:"Your Company Trading Name"`
	TaxID          string `json:"tax_id" example:"00000000000000"`
	OwnerName      string `json:"owner_name" example:"John"`
	OwnerSurname   string `json:"owner_surname" example:"Clark"`
	OwnerEmail     string `json:"owner_email" example:"john.clark@gmail.com" validate:"required,email"`
	OwnerPhone     string `json:"owner_phone" example:"+15555555555" validate:"required,e164"`
	OwnerPassword  string `json:"owner_password" example:"1SecurePswd!" validate:"required,myPasswordValidation"`
	OwnerTimeZone  string `json:"owner_time_zone" example:"America/Sao_Paulo" validate:"required"` // Use a valid timezone
	StartSubdomain string `json:"start_subdomain" example:"your-company-subdomain" validate:"required,mySubdomainValidation"`
}

// @description	Company Full DTO
// @name			CompanyFullDTO
// @tag.name		company.full.dto
type CompanyFull struct {
	CompanyBase
	Employees []*EmployeeBase `json:"employees"`
	Branches  []*BranchBase   `json:"branches"`
	Services  []*ServiceBase  `json:"services"`
}

// @description	Company Base DTO
// @name			CompanyBaseDTO
// @tag.name		company.base.dto
type CompanyBase struct {
	ID         uuid.UUID          `json:"id" example:"00000000-0000-0000-0000-000000000000"` // Primary key
	LegalName  string             `json:"legal_name" example:"Your Company Legal Name"`
	TradeName  string             `json:"trading_name" example:"Your Company Trading Name"`
	TaxID      string             `json:"tax_id" example:"00000000000000"`
	Design     mJSON.DesignConfig `json:"design"`
	Sectors    []*Sector          `json:"sectors"`
	Subdomains []*Subdomain       `json:"subdomains"`
}
