package models

import (
	"gorm.io/gorm"
)

// Company holds an array of CompanyTypes.
type Company struct {
	gorm.Model
	Name         string        `gorm:"not null;unique" json:"name"`
	CompanyTypes []CompanyType `gorm:"many2many:company_company_types;constraint:OnDelete:CASCADE" json:"company_types"` // Many-to-many relation with a custom join table
	TaxID        string        `gorm:"not null;unique" json:"tax_id"`
}