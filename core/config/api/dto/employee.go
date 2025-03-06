package DTO

type UpdateEmployeeSwagger struct {
	Name    string `json:"name" example:"John"`
	Surname string `json:"surname" example:"Clark"`
}

type CreateEmployee struct {
	Name      string `json:"name" example:"John"`
	Surname   string `json:"surname" example:"Doe"`
	Email     string `json:"email" example:"john.doe@example.com"`
	Phone     string `json:"phone" example:"+15555555555"`
	Password  string `json:"password" example:"1VerySecurePassword!"`
	CompanyID uint   `json:"company_id" example:"1"`
	UserID    uint   `json:"user_id" example:"1"`
}

type Employee struct {
	ID uint `json:"id"`
	CreateEmployee
	VerificationCode string             `json:"verification_code" example:"123456"`
	Verified         bool               `json:"verified" example:"false"`
	Branches         []BranchPopulated  `json:"branches"`
	Services         []ServicePopulated `json:"services"`
	Tags             []string           `json:"tag" example:"[\"super-admin\", \"branch-manager\"]"`
}
