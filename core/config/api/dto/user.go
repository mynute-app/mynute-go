package DTO

type LoginUser struct {
	Email    string `json:"email" example:"john.doe@example.com"`
	Password string `json:"password" example:"1VerySecurePassword!"`
}

type CreatedUser struct {
	ID      uint   `json:"id" example:"1"`
	Name    string `json:"name" example:"John"`
	Surname string `json:"surname" example:"Doe"`
	Email   string `json:"email" example:"john.doe@example.com"`
	Phone   string `json:"phone" example:"+15555555555"`
}

type CreateUser struct {
	Name     string `json:"name" example:"John"`
	Surname  string `json:"surname" example:"Doe"`
	Email    string `json:"email" example:"john.doe@example.com"`
	Phone    string `json:"phone" example:"+15555555555"`
	Password string `json:"password" example:"1VerySecurePassword!"`
}

type User struct {
	ID uint `json:"id" example:"1"`
	Name    string `json:"name" example:"John"`
	Surname string `json:"surname" example:"Doe"`
	Email   string `json:"email" example:"john.doe@example.com"`
	Phone   string `json:"phone" example:"+15555555555"`
	VerificationCode string `json:"verification_code" example:"123456"`
	Verified         bool   `json:"verified" example:"false"`
	EmployeeID       uint   `json:"employee_id" example:"1"`
	CompanyID        uint   `json:"company_id" example:"1"`
}

type UserPopulated struct {
	ID      uint   `json:"id" example:"1"`
	Name    string `json:"name" example:"John"`
	Surname string `json:"surname" example:"Doe"`
	Email   string `json:"email" example:"john.doe@example.com"`
	Phone   string `json:"phone" example:"+1-555-555-5555"`
}
