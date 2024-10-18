package DTO

type Employee struct {
	ID        uint               `json:"id"`
	CompanyID uint               `json:"company_id"`
	Name      string             `json:"name"`
	Surname   string             `json:"surname"`
	Email     string             `json:"email"`
	Phone     string             `json:"phone"`
	Branches  []BranchPopulated  `json:"branches"`
	Services  []ServicePopulated `json:"services"`
}

type EmployeePopulated struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
}
