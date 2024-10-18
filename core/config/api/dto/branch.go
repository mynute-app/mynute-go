package DTO

type Branch struct {
	ID        uint                `json:"id"`
	CompanyID uint                `json:"company_id"`
	Name      string              `json:"name"`
	Employees []EmployeePopulated `json:"employees"`
	Services  []ServicePopulated  `json:"services"`
}

type BranchPopulated struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
