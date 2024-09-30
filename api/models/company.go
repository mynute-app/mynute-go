package models

import (
	"errors"
	"regexp"
	"gorm.io/gorm"
)

// CompanyType: Represents different types of companies
type CompanyType struct {
	gorm.Model
	Name string `gorm:"not null;unique" json:"name"`
}

// First step: Choosing the company.
type Company struct {
	gorm.Model
	Name  string        `gorm:"not null;unique" json:"name"`
	Types []CompanyType `gorm:"many2many:company_types;"` // Many-to-many relation
	TaxID string        `gorm:"unique" json:"tax_id"`     // TaxID must be unique
}

// BeforeSave is a GORM hook that runs before the record is saved
func (c *Company) BeforeSave(tx *gorm.DB) (err error) {
	if !c.ValidateTaxID() {
		return errors.New("invalid TaxID: it must contain only 15 numeric characters")
	}
	return nil
}

// ValidateTaxID checks if the TaxID is a 15-character string containing only numbers
func (c *Company) ValidateTaxID() bool {
	// Define the regular expression to match exactly 15 digits
	re := regexp.MustCompile(`^\d{14}$`)
	// Check if the TaxID matches the regular expression
	return re.MatchString(c.TaxID)
}
