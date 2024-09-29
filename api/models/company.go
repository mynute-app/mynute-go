package models

import "gorm.io/gorm"

// CompanyType: Represents different types of companies
type CompanyType struct {
	gorm.Model
	Name string `json:"name"`
}

// First step: Choosing the company.
type Company struct {
	gorm.Model
	Name  string        `json:"name"`
	Types []CompanyType `gorm:"many2many:company_types;"` // Many-to-many relation
	TaxID string        `gorm:"unique" json:"tax_id"` // TaxID must be unique
}
