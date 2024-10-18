package DTO

type Branch struct {
	ID        uint             `json:"id"`
	Name      string           `json:"name"`
	CompanyID uint             `json:"company_id"`
	Company   CompanyPopulated `json:"company"`
	Employees []Employee       `json:"employees"`
	Services  []Service        `json:"services"`
}

type BranchPopulated struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
