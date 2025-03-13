package model

import (
	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model

	CompanyID      uint          `gorm:"not null;index" json:"company_id"`
	Company        Company       `gorm:"foreignKey:CompanyID;references:ID;constraint:OnDelete:CASCADE;"`

	UserID         uint          `gorm:"unique;not null;index" json:"user_id"`
	User           User          `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`

	Branches       []Branch      `gorm:"many2many:employee_branches;"`

	Services       []Service     `gorm:"many2many:employee_services;"`

	AvailableSlots []TimeRange   `gorm:"type:json" json:"available_slots"`

	Appointments   []Appointment `gorm:"foreignKey:EmployeeID"`
	
	Tags           []string      `gorm:"type:json" json:"tag"`
}
