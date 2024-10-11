package models

import "gorm.io/gorm"

// Second step: Choosing the branch.
type Branch struct {
	gorm.Model
	CompanyID uint    `json:"company_id"`                 // Foreign key for Company
	Name      string  `json:"name"`
	Company   Company `gorm:"constraint:OnDelete:CASCADE;"` // No need for foreignKey here, it's inferred from CompanyID
}
