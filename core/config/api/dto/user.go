package DTO

type LoginUser struct {
	Email    string `json:"email" example:"john.doe@example.com"`
	Password string `json:"password" example:"1SecurePswd!"`
}

type CreateUser struct {
	Name     string `json:"name" example:"John"`
	Surname  string `json:"surname" example:"Doe"`
	Email    string `json:"email" example:"john.doe@example.com"`
	Phone    string `json:"phone" example:"+15555555555"`
	Password string `json:"password" example:"1SecurePswd!"`
}

type User struct {
	ID             uint          `json:"id" example:"1"`
	Name           string        `json:"name" example:"John"`
	Surname        string        `json:"surname" example:"Doe"`
	Email          string        `json:"email" example:"john.doe@example.com"`
	Phone          string        `json:"phone" example:"+15555555555"`
	Verified       bool          `json:"verified" example:"false"`
	AvailableSlots []TimeRange   `json:"available_slots"`
	Appointments   []*Appointment `json:"appointments"`
}

type UserPopulated struct {
	ID       uint   `json:"id" example:"1"`
	Name     string `json:"name" example:"John"`
	Surname  string `json:"surname" example:"Doe"`
	Email    string `json:"email" example:"john.doe@example.com"`
	Phone    string `json:"phone" example:"+1-555-555-5555"`
	Verified bool   `json:"verified" example:"false"`
}
