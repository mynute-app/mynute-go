package DTO

type Company struct {
	ID           uint          `json:"id"`
	Name         string        `json:"name"`
	TaxID        string        `json:"tax_id"`
	CompanyTypes []CompanyType `json:"company_types"`
	Employees    []Employee    `json:"employees"`
	Branches     []Branch      `json:"branches"`
	Services     []Service     `json:"services"`
}
