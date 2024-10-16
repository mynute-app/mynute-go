package models

import (
	"gorm.io/gorm"
)

// Company holds an array of CompanyTypes.
type Company struct {
	gorm.Model
	Name         string        `gorm:"not null;unique" json:"name"`
	TaxID        string        `gorm:"not null;unique" json:"tax_id"`
	CompanyTypes []CompanyType `gorm:"many2many:company_company_types;constraint:OnDelete:CASCADE" json:"company_types"`
	Employees    []Employee    `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"employees"`
	Branches     []Branch      `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"branches"` // Explicit foreign key definition
	Services     []Service     `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"services"` // Explicit foreign key definition
}
