package model

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	Name             string        `gorm:"not null" json:"name"`
	Surname          string        `json:"surname"`
	Role             string        `json:"role" example:"user"`
	Email            string        `gorm:"not null;unique" json:"email"`
	Phone            string        `gorm:"not null;unique" json:"phone"`
	Password         string        `gorm:"not null" json:"password"`
	Tags             []string      `gorm:"type:json" json:"tag"`
	CompanyID        uint          `gorm:"not null;index;foreignKey:CompanyID;references:ID;constraint:OnDelete:CASCADE;" json:"company_id"`
	Company          Company       `gorm:"constraint:OnDelete:CASCADE;"`
	UserID           uint          `gorm:"not null;index;foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;" json:"user_id"`
	Branches         []Branch      `gorm:"many2many:employee_branches;"` // Many-to-many relation with Branch
	Services         []Service     `gorm:"many2many:employee_services;"` // Many-to-many relation with Service
	AvailableSlots   []TimeRange   `gorm:"type:json" json:"available_slots"`
	Appointments     []Appointment `gorm:"foreignKey:EmployeeID"`
	VerificationCode string        `json:"verification_code"`
	Verified         bool          `gorm:"not null" json:"verified"`
}
