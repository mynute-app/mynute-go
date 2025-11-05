package DTO

import (
	"github.com/google/uuid"
)

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

// @description	Employee Base DTO - Auth service only handles basic employee info
// @name			EmployeeBaseDTO
// @tag.name		employee.base.dto
type EmployeeBase struct {
	ID        uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	CompanyID uuid.UUID `json:"company_id" example:"00000000-0000-0000-0000-000000000000"`
	Name      string    `json:"name" example:"John"`
	Surname   string    `json:"surname" example:"Doe"`
	Email     string    `json:"email" example:"john.doe@example.com"`
	Phone     string    `json:"phone" example:"+15555555555"`
	Verified  bool      `json:"verified" example:"true"`
}

