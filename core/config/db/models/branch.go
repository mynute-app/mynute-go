package models

import "gorm.io/gorm"

// Second step: Choosing the branch.
type Branch struct {
	gorm.Model
	CompanyID uint    `json:"company_id"`
	Name      string  `json:"name"`
	Company   Company `gorm:"foreignKey:CompanyID"` // Foreign key relation to Company
}