package model

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	CompanyID      uint           `gorm:"not null;index;foreignKey:CompanyID;references:ID;constraint:OnDelete:CASCADE;" json:"company_id"`
	Company        Company        `gorm:"constraint:OnDelete:CASCADE;"`
	UserID         uint           `gorm:"unique;not null;index;foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;" json:"user_id"`
	User           User           `gorm:"constraint:OnDelete:CASCADE;"`
	Branches       []Branch       `gorm:"many2many:employee_branches;"` // Many-to-many relation with Branch
	Services       []Service      `gorm:"many2many:employee_services;"` // Many-to-many relation with Service
	AvailableSlots []TimeRange    `gorm:"type:json" json:"available_slots"`
	Appointments   []Appointment  `gorm:"foreignKey:EmployeeID"`
	Tags           []string       `gorm:"type:json" json:"tag"`
}
