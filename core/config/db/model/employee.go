package model

import (
	"gorm.io/gorm"
)

type CreateEmployee struct {
	CreateUser
	CompanyID uint `gorm:"not null" json:"company_id"`
}

type Employee struct {
	gorm.Model
	GeneralUserInfo
	CompanyID      uint          `gorm:"not null;index;foreignKey:CompanyID;references:ID;constraint:OnDelete:CASCADE;" json:"company_id"`
	Company        Company       `gorm:"constraint:OnDelete:CASCADE;"`
	Branches       []Branch      `gorm:"many2many:employee_branches;"` // Many-to-many relation with Branch
	Services       []Service     `gorm:"many2many:employee_services;"` // Many-to-many relation with Service
	AvailableSlots []TimeRange   `gorm:"type:json" json:"available_slots"`
	Appointments   []Appointment `gorm:"foreignKey:EmployeeID"`
}
