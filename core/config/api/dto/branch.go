package DTO

import (
	mJSON "mynute-go/core/config/db/model/json"

	"github.com/google/uuid"
)

// @description	Branch Create DTO
// @name			BranchCreateDTO
// @tag.name		branch.create.dto
type CreateBranch struct {
	CompanyID    uuid.UUID `json:"company_id" example:"00000000-0000-0000-0000-000000000000"`
	Name         string    `json:"name" example:"Main Branch"`
	Street       string    `json:"street" example:"123 Main St"`
	Number       string    `json:"number" example:"456"`
	Complement   string    `json:"complement" example:"Suite 100"`
	Neighborhood string    `json:"neighborhood" example:"Downtown"`
	ZipCode      string    `json:"zip_code" example:"10001"`
	City         string    `json:"city" example:"New York"`
	State        string    `json:"state" example:"NY"`
	Country      string    `json:"country" example:"USA"`
	TimeZone     string    `json:"time_zone" example:"America/New_York"` // Time zone in IANA format
}

// @description	Branch Update DTO
// @name			BranchUpdateDTO
// @tag.name		branch.update.dto
type UpdateBranch struct {
	CompanyID uuid.UUID `json:"company_id" example:"00000000-0000-0000-0000-000000000000"`
	Name      string    `json:"name" example:"Main Branch Updated"`
	Street    string    `gorm:"not null" json:"street" example:"556 Main St"`
}

// @description	Branch Full DTO
// @name			BranchFullDTO
// @tag.name		branch.full.dto
type BranchFull struct {
	BranchBase
	Employees      []*EmployeeBase  `json:"employees"`
	Services       []*ServiceBase   `json:"services"`
	Appointments   []*Appointment   `json:"appointments"`
	ServiceDensity []ServiceDensity `json:"service_density"`
}

// @description	Branch Base DTO
// @name			BranchBaseDTO
// @tag.name		branch.base.dto
type BranchBase struct {
	ID                  uuid.UUID          `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	CompanyID           uuid.UUID          `json:"company_id" example:"00000000-0000-0000-0000-000000000000"`
	Name                string             `json:"name" example:"Main Branch"`
	Street              string             `json:"street" example:"123 Main St"`
	Number              string             `json:"number" example:"456"`
	Complement          string             `json:"complement" example:"Suite 100"`
	Neighborhood        string             `json:"neighborhood" example:"Downtown"`
	ZipCode             string             `json:"zip_code" example:"10001"`
	City                string             `json:"city" example:"New York"`
	State               string             `json:"state" example:"NY"`
	Country             string             `json:"country" example:"USA"`
	TimeZone            string             `json:"time_zone" example:"America/New_York"` // Time zone in IANA format
	TotalServiceDensity int32              `json:"total_service_density" example:"100"`
	Design              mJSON.DesignConfig `json:"design"`
}

type ServiceDensity struct {
	ServiceID           uuid.UUID `json:"service_id" example:"00000000-0000-0000-0000-000000000000"`
	MaxSchedulesOverlap uint      `json:"max_schedules_overlap" example:"5"`
}
