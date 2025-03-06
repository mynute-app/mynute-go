package DTO

type CreateCompany struct {
	Name    string   `json:"name" example:"Your Company Name"`
	TaxID   string   `json:"tax_id" example:"00000000000000"`
	Sectors []Sector `json:"sectors"`
}

// @description	Company Full DTO
// @name			CompanyFullDTO
// @tag.name		company.full.dto
type Company struct {
	ID uint `json:"id"` // Primary key
	CreateCompany
	Employees []UserPopulated    `json:"employees"`
	Branches  []Branch           `json:"branches"`
	Services  []ServicePopulated `json:"services"`
}

// @description	Company DTO Populated
// @name			CompanyPopulatedDTO
// @tag.name		company_populated.dto
type CompanyPopulated struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	TaxID string `json:"tax_id"`
}
