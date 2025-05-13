package DTO

import "github.com/google/uuid"

type CreateBranch struct {
	CompanyID    uuid.UUID `json:"company_id" example:"00000000-0000-0000-0000-000000000000"`
	Name         string    `json:"name" example:"Main Branch"`
	Street       string    `gorm:"not null" json:"street" example:"123 Main St"`
	Number       string    `gorm:"not null" json:"number" example:"456"`
	Complement   string    `json:"complement" example:"Suite 100"`
	Neighborhood string    `gorm:"not null" json:"neighborhood" example:"Downtown"`
	ZipCode      string    `gorm:"not null" json:"zip_code" example:"10001"`
	City         string    `gorm:"not null" json:"city" example:"New York"`
	State        string    `gorm:"not null" json:"state" example:"NY"`
	Country      string    `gorm:"not null" json:"country" example:"USA"`
}

type UpdateBranch struct {
	CompanyID uuid.UUID `json:"company_id" example:"00000000-0000-0000-0000-000000000000"`
	Name      string    `json:"name" example:"Main Branch Updated"`
	Street    string    `gorm:"not null" json:"street" example:"556 Main St"`
}

type Branch struct {
	ID             uuid.UUID            `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Name           string               `json:"name" example:"Main Branch"`
	Street         string               `json:"street" example:"123 Main St"`
	Number         string               `json:"number" example:"456"`
	Complement     string               `json:"complement" example:"Suite 100"`
	Neighborhood   string               `json:"neighborhood" example:"Downtown"`
	ZipCode        string               `json:"zip_code" example:"10001"`
	City           string               `json:"city" example:"New York"`
	State          string               `json:"state" example:"NY"`
	Country        string               `json:"country" example:"USA"`
	Employees      []*EmployeePopulated `json:"employees"`
	Services       []*ServicePopulated  `json:"services"`
	CompanyID      uuid.UUID            `json:"company_id" example:"00000000-0000-0000-0000-000000000000"`
	Appointments   []*Appointment       `json:"appointments"`
	ServiceDensity []ServiceDensity     `json:"service_density"`
	BranchDensity  uint                 `json:"branch_density"`
}

type BranchPopulated struct {
	ID             uuid.UUID        `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Name           string           `json:"name" example:"Main Branch"`
	Street         string           `json:"street" example:"123 Main St"`
	Number         string           `json:"number" example:"456"`
	Complement     string           `json:"complement" example:"Suite 100"`
	Neighborhood   string           `json:"neighborhood" example:"Downtown"`
	ZipCode        string           `json:"zip_code" example:"10001"`
	City           string           `json:"city" example:"New York"`
	State          string           `json:"state" example:"NY"`
	Country        string           `json:"country" example:"USA"`
	ServiceDensity []ServiceDensity `json:"service_density"`
	BranchDensity  uint             `json:"branch_density"`
}

type ServiceDensity struct {
	ServiceID           uuid.UUID `json:"service_id" example:"00000000-0000-0000-0000-000000000000"`
	MaxSchedulesOverlap uint      `json:"max_schedules_overlap" example:"5"`
}
