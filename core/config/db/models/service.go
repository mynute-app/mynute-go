package models

import "gorm.io/gorm"

// Third step: Choosing the service.
type Service struct {
	gorm.Model
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       int32      `json:"price"`
	Duration    int        `json:"duration"`                   // Duration in minutes (e.g., 30, 60, etc.)
	Branches    []Branch   `gorm:"many2many:branch_services;"` // Many-to-many relation
	CompanyID   uint       `json:"company_id"`
	Company     Company    `gorm:"foreignKey:CompanyID"`         // Foreign key relation to Company
	Employees   []Employee `gorm:"many2many:employee_services;"` // Many-to-many relation
}
