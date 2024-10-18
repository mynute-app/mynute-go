package DTO

type Company struct {
	ID           uint                `json:"id"`
	Name         string              `json:"name"`
	TaxID        string              `json:"tax_id"`
	CompanyTypes []CompanyType       `json:"company_types"`
	Employees    []EmployeePopulated `json:"employees"`
	Branches     []BranchPopulated   `json:"branches"`
	Services     []ServicePopulated  `json:"services"`
}

type CompanyPopulated struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	TaxID string `json:"tax_id"`
}
