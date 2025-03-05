package DTO

type CreateUser struct {
	Name      string `json:"name" example:"John"`
	Surname   string `json:"surname" example:"Doe"`
	Email     string `json:"email" example:"john.doe@example.com"`
	Password  string `json:"password" example:"1VerySecurePassword!"`
	Phone     string `json:"phone" example:"+1-555-555-5555"`
}

type User struct {
	ID uint `json:"id" example:"1"`
	CreateUser
	VerificationCode string `json:"verification_code" example:"123456"`
	Verified         bool   `json:"verified" example:"false"`
}

type UserPopulated struct {
	ID      uint   `json:"id" example:"1"`
	Name    string `json:"name" example:"John"`
	Surname string `json:"surname" example:"Doe"`
	Email   string `json:"email" example:"john.doe@example.com"`
	Phone   string `json:"phone" example:"+1-555-555-5555"`
}
