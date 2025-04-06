package model

import "github.com/google/uuid"

// Branch model
type Branch struct {
	BaseModel
	Name           string           `gorm:"not null" json:"name"`
	Street         string           `gorm:"not null" json:"street"`
	Number         string           `gorm:"not null" json:"number"`
	Complement     string           `json:"complement"`
	Neighborhood   string           `gorm:"not null" json:"neighborhood"`
	ZipCode        string           `gorm:"not null" json:"zip_code"`
	City           string           `gorm:"not null" json:"city"`
	State          string           `gorm:"not null" json:"state"`
	Country        string           `gorm:"not null" json:"country"`
	CompanyID      uuid.UUID        `gorm:"not null;index" json:"company_id"`
	Company        *Company         `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;" json:"company"`
	Employees      []*Employee      `gorm:"many2many:employee_branches;constraint:OnDelete:CASCADE"`              // Many-to-many relation with Employee
	Services       []*Service       `gorm:"many2many:branch_services;constraint:OnDelete:CASCADE"`                // Many-to-many relation with Service
	Appointments   []Appointment    `gorm:"foreignKey:BranchID;constraint:OnDelete:CASCADE;" json:"appointments"` // One-to-many relation with Appointment
	ServiceDensity []ServiceDensity `gorm:"type:jsonb" json:"service_density"`                                    // One-to-many relation with ServiceDensity
	BranchDensity  uint             `gorm:"not null;default:1" json:"branch_density"`
}

type ServiceDensity struct {
	ServiceID           uuid.UUID `json:"service_id"`
	MaxSchedulesOverlap uint      `json:"max_schedules_overlap"`
}
