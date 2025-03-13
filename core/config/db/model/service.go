package model

import "gorm.io/gorm"

type CreateService struct {
	Name        string `gorm:"not null" json:"name"`
	Description string `gorm:"not null" json:"description"`
	Price       int32  `gorm:"not null" json:"price"`
	Duration    uint   `gorm:"not null" json:"duration"`
	CompanyID   uint   `gorm:"not null" json:"company_id"`
}

// Third step: Choosing the service.
type Service struct {
	gorm.Model
	Name        string `gorm:"not null" json:"name"`
	Description string `gorm:"not null" json:"description"`
	Price       int32  `gorm:"not null" json:"price"`
	Duration    uint   `gorm:"not null" json:"duration"` // Duration in minutes

	CompanyID uint    `gorm:"not null;index" json:"company_id"`
	Company   Company `gorm:"foreignKey:CompanyID;references:ID;constraint:OnDelete:CASCADE;"`

	Employees []Employee `gorm:"many2many:employee_services;constraint:OnDelete:CASCADE;" json:"employees"` // Many-to-many relation with Employee
	Branches  []Branch   `gorm:"many2many:branch_services;constraint:OnDelete:CASCADE;" json:"branches"`    // Many-to-many relation with Branch
}
