package DTO

type CreateCompany struct {
	Name  string `json:"name" example:"Your Company Name"`
	TaxID string `json:"tax_id" example:"00000000000000"`
}

// @description	Company Full DTO
// @name			CompanyFullDTO
// @tag.name		company.full.dto
type Company struct {
	ID        uint                `json:"id"` // Primary key
	Name      string              `json:"name" example:"Your Company Name"`
	TaxID     string              `json:"tax_id" example:"00000000000000"`
	Employees []EmployeePopulated `json:"employees"`
	Branches  []BranchPopulated   `json:"branches"`
	Services  []ServicePopulated  `json:"services"`
	Sectors   []Sector            `json:"sectors"`
}

// @description	Company DTO Populated
// @name			CompanyPopulatedDTO
// @tag.name		company_populated.dto
type CompanyPopulated struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	TaxID string `json:"tax_id"`
}
