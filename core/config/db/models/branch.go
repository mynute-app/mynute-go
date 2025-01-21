package models

import "gorm.io/gorm"

// Branch model
type Branch struct {
	gorm.Model
	Name      string    `gorm:"not null" json:"name"`
	CompanyID uint      `json:"company_id"`                  // Foreign key to Company
	Employees []User    `gorm:"many2many:branch_employees;"` // Many-to-many relation with Employee
	Services  []Service `gorm:"many2many:branch_services;"`  // Many-to-many relation with Service
}
