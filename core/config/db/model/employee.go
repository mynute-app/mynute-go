package model

import (
	"gorm.io/gorm"
)

type CreateEmployee struct {
	CompanyID uint     `gorm:"not null" json:"company_id"`
	Name      string   `gorm:"not null" json:"name" example:"Joseph"`
	Surname   string   `json:"surname" example:"Doe"`
	Role      string   `json:"role" example:"user"`
	Email     string   `gorm:"not null;unique" json:"email" example:"joseph.doe@example.com"`
	Phone     string   `gorm:"not null;unique" json:"phone" example:"+15555555551"`
	Tags      []string `gorm:"type:json" json:"tags"` // Tags for the user
}

type Employee struct {
	gorm.Model

	CompanyID *uint   `gorm:"not null;index" json:"company_id"`
	Company   Company `gorm:"foreignKey:CompanyID;references:ID;constraint:OnDelete:CASCADE;"`

	UserID *uint `gorm:"unique;not null;index" json:"user_id"`
	User   User  `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`

	Branches []Branch `gorm:"many2many:employee_branches;"`

	Services []Service `gorm:"many2many:employee_services;"`

	AvailableSlots []TimeRange `gorm:"type:json" json:"available_slots"`

	Appointments []Appointment `gorm:"foreignKey:EmployeeID"`

	Tags []string `gorm:"type:json" json:"tag"`
}
