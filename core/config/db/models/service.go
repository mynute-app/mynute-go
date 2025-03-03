package models

import "gorm.io/gorm"

// Third step: Choosing the service.
type Service struct {
	gorm.Model
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       int32    `json:"price"`
	Duration    int      `json:"duration"`                                                                                         // Duration in minutes
	CompanyID   uint     `gorm:"not null;index;foreignKey:CompanyID;references:ID;constraint:OnDelete:CASCADE;" json:"company_id"` // Foreign key to Company
	Users       []User   `gorm:"many2many:employee_services;" json:"users"`                                                                     // Many-to-many relation with Employee
	Branches    []Branch `gorm:"many2many:branch_services;" json:"branches"`                                                                       // Many-to-many relation with Branch
}
