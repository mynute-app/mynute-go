package DTO

type Company struct {
	ID           uint               `json:"id"` // Primary key
	Name         string             `json:"name"`
	TaxID        string             `json:"tax_id"`
	CompanyTypes []CompanyType      `json:"company_types"`
	Employees    []UserPopulated    `json:"employees"`
	Branches     []BranchPopulated  `json:"branches"`
	Services     []ServicePopulated `json:"services"`
}

type CompanyPopulated struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	TaxID string `json:"tax_id"`
}
