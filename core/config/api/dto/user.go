package DTO

type User struct {
	ID               uint               `json:"id" example:"1"`
	CompanyID        uint               `json:"company_id" example:"1"` // ID da empresa
	Name             string             `json:"name" example:"John"`
	Surname          string             `json:"surname" example:"Doe"`
	Email            string             `json:"email" example:"john.doe@example.com"`
	Password         string             `json:"password" example:"securepassword"`
	VerificationCode string             `json:"verification_code" example:"123456"`
	Phone            string             `json:"phone" example:"+1-555-555-5555"`
	Branches         []BranchPopulated  `json:"branches"`
	Services         []ServicePopulated `json:"services"`
	Tag              []string           `json:"tag" example:"[\"super-admin\", \"branch-manager\"]"`
}

type UserPopulated struct {
	ID      uint   `json:"id" example:"1"`
	Name    string `json:"name" example:"John"`
	Surname string `json:"surname" example:"Doe"`
	Email   string `json:"email" example:"john.doe@example.com"`
	Phone   string `json:"phone" example:"+1-555-555-5555"`
}
