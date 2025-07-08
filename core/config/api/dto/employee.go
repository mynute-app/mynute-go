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
	CompanyID uuid.UUID `json:"company_id"`
	Name      string    `json:"name" example:"Joseph"`
	Surname   string    `json:"surname" example:"Doe"`
	Role      string    `json:"role" example:"client"`
	Email     string    `json:"email" example:"joseph.doe@example.com"`
	Phone     string    `json:"phone" example:"+15555555551"`
	Password  string    `json:"password" example:"1SecurePswd!"`
	TimeZone  string    `json:"time_zone" example:"America/Sao_Paulo"` // Use a valid timezone
}

// type EmployeeWorkSchedule struct {
// 	Monday    []EmployeeWorkRange `json:"monday"`
// 	Tuesday   []EmployeeWorkRange `json:"tuesday"`
// 	Wednesday []EmployeeWorkRange `json:"wednesday"`
// 	Thursday  []EmployeeWorkRange `json:"thursday"`
// 	Friday    []EmployeeWorkRange `json:"friday"`
// 	Saturday  []EmployeeWorkRange `json:"saturday"`
// 	Sunday    []EmployeeWorkRange `json:"sunday"`
// }

// type EmployeeWorkRange struct {
// 	Start    string      `json:"start"` // Store as "15:30:00"
// 	End      string      `json:"end"`   // Store as "18:00:00"
// 	BranchID uuid.UUID   `json:"branch_id"`
// 	Services []uuid.UUID `json:"services"` // Store as UUIDs
// }

// @description	Employee Full DTO
// @name			EmployeeFullDTO
// @tag.name		employee.full.dto
type EmployeeFull struct {
	EmployeeBase
	Verified             bool                `json:"verified" example:"true"`
	EmployeeWorkSchedule []EmployeeWorkRange `json:"work_schedule"`
	Appointments         []Appointment       `json:"appointments"`
	Branches             []BranchBase        `json:"branches"`
	Services             []ServiceBase       `json:"services"`
	Roles                []RolePopulated     `json:"roles"`
}

// @description	Employee Base DTO
// @name			EmployeeBaseDTO
// @tag.name		employee.base.dto
type EmployeeBase struct {
	ID        uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	CompanyID uuid.UUID `json:"company_id" example:"00000000-0000-0000-0000-000000000000"`
	Name      string    `json:"name" example:"John"`
	Surname   string    `json:"surname" example:"Doe"`
	Email     string    `json:"email" example:"john.doe@example.com"`
	Phone     string    `json:"phone" example:"+15555555555"`
}
