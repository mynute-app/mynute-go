package model

import (
	"gorm.io/gorm"
)

type UpdateEmployeeSwagger struct {
	Name     string `json:"name" example:"John"`
	Surname  string `json:"surname" example:"Clark"`
}

type CreateEmployee struct {
	CreateUser
	CompanyID uint `gorm:"not null" json:"company_id"`
}

type Employee struct {
	gorm.Model
	GeneralUserInfo
	CompanyID      uint          `gorm:"not null;index;foreignKey:CompanyID;references:ID;constraint:OnDelete:CASCADE;" json:"company_id"`
	Company        Company       `gorm:"constraint:OnDelete:CASCADE;"`
	UserID         uint          `gorm:"not null;index;foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;" json:"user_id"`
	Branches       []Branch      `gorm:"many2many:employee_branches;"` // Many-to-many relation with Branch
	Services       []Service     `gorm:"many2many:employee_services;"` // Many-to-many relation with Service
	AvailableSlots []TimeRange   `gorm:"type:json" json:"available_slots"`
	Appointments   []Appointment `gorm:"foreignKey:EmployeeID"`
}
