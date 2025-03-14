package DTO

type TimeRange struct {
	Start string `json:"start" example:"09:00"`
	End   string `json:"end" example:"17:00"`
}

type UpdateEmployeeSwagger struct {
	Name    string `json:"name" example:"John"`
	Surname string `json:"surname" example:"Clark"`
}

type CreateEmployee struct {
	CompanyID uint   `json:"company_id"`
	Name      string `json:"name" example:"Joseph"`
	Surname   string `json:"surname" example:"Doe"`
	Role      string `json:"role" example:"user"`
	Email     string `json:"email" example:"joseph.doe@example.com"`
	Phone     string `json:"phone" example:"+15555555551"`
}

type Employee struct {
	ID             uint               `json:"id" example:"1"`
	Name           string             `json:"name" example:"John"`
	Surname        string             `json:"surname" example:"Doe"`
	Role           string             `json:"role" example:"user"`
	Email          string             `json:"email" example:"john.doe@example.com"`
	Phone          string             `json:"phone" example:"+15555555555"`
	Tags           []string           `json:"tags" example:"[\"tag1\", \"tag2\"]"`
	Verified       bool               `json:"verified" example:"true"`
	AvailableSlots []TimeRange        `json:"available_slots" example:"[{\"start\":\"09:00\", \"end\":\"17:00\"}]"`
	Appointments   []Appointment      `json:"appointments"`
	CompanyID      uint               `json:"company_id" example:"1"`
	Company        CompanyPopulated   `json:"company"`
	Branches       []BranchPopulated  `json:"branches"`
	Services       []ServicePopulated `json:"services"`
}

type EmployeePopulated struct {
	ID      uint     `json:"id" example:"1"`
	Name    string   `json:"name" example:"John"`
	Surname string   `json:"surname" example:"Doe"`
	Role    string   `json:"role" example:"user"`
	Email   string   `json:"email" example:"john.doe@example.com"`
	Phone   string   `json:"phone" example:"+15555555555"`
	Tags    []string `json:"tags" example:"[\"tag1\", \"tag2\"]"`
}
