package DTO

import "github.com/google/uuid"

type TimeRange struct {
	Start string `json:"start" example:"09:00"`
	End   string `json:"end" example:"17:00"`
}

type LoginEmployee struct {
	Email    string `json:"email" example:"john.clark@gmail.com"`
	Password string `json:"password" example:"1SecurePswd!"`
}

type UpdateEmployeeSwagger struct {
	Name    string `json:"name" example:"John"`
	Surname string `json:"surname" example:"Clark"`
}

type CreateEmployee struct {
	CompanyID    uuid.UUID    `json:"company_id"`
	Name         string       `json:"name" example:"Joseph"`
	Surname      string       `json:"surname" example:"Doe"`
	Role         string       `json:"role" example:"client"`
	Email        string       `json:"email" example:"joseph.doe@example.com"`
	Phone        string       `json:"phone" example:"+15555555551"`
	Password     string       `json:"password" example:"1SecurePswd!"`
	WorkSchedule WorkSchedule `json:"work_schedule"`
}

type WorkSchedule struct {
	Monday    []WorkRange `json:"monday"`
	Tuesday   []WorkRange `json:"tuesday"`
	Wednesday []WorkRange `json:"wednesday"`
	Thursday  []WorkRange `json:"thursday"`
	Friday    []WorkRange `json:"friday"`
	Saturday  []WorkRange `json:"saturday"`
	Sunday    []WorkRange `json:"sunday"`
}

type WorkRange struct {
	Start    string    `json:"start"` // Store as "15:30:00"
	End      string    `json:"end"`   // Store as "18:00:00"
	BranchID uuid.UUID `json:"branch_id"`
}

type Employee struct {
	ID           uuid.UUID          `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Name         string             `json:"name" example:"John"`
	Surname      string             `json:"surname" example:"Doe"`
	Email        string             `json:"email" example:"john.doe@example.com"`
	Phone        string             `json:"phone" example:"+15555555555"`
	Tags         []string           `json:"tags" example:"[\"tag1\", \"tag2\"]"`
	Verified     bool               `json:"verified" example:"true"`
	WorkSchedule WorkSchedule       `json:"work_schedule"`
	Appointments []Appointment      `json:"appointments"`
	CompanyID    uuid.UUID          `json:"company_id" example:"00000000-0000-0000-0000-000000000000"`
	Branches     []BranchPopulated  `json:"branches"`
	Services     []ServicePopulated `json:"services"`
	Roles        []RolePopulated    `json:"roles"`
}

type EmployeePopulated struct {
	ID        uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Name      string    `json:"name" example:"John"`
	Surname   string    `json:"surname" example:"Doe"`
	Role      string    `json:"role" example:"client"`
	Email     string    `json:"email" example:"john.doe@example.com"`
	Phone     string    `json:"phone" example:"+15555555555"`
	CompanyID uuid.UUID `json:"company_id" example:"00000000-0000-0000-0000-000000000000"`
	Tags      []string  `json:"tags" example:"[\"tag1\", \"tag2\"]"`
}
