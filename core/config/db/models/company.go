package models

import (
	"errors"
	"log"
	"regexp"

	"gorm.io/gorm"
)

// Company holds an array of CompanyTypes.
type Company struct {
	gorm.Model
	Name         string        `gorm:"not null;unique" json:"name"`
	CompanyTypes []CompanyType `gorm:"many2many:company_company_types;ForeignKey:id;References:id" json:"company_types"` // Many-to-many relation with a custom join table
	TaxID        string        `gorm:"not null;unique" json:"tax_id"`
}

// BeforeSave is a GORM hook that runs before the record is saved
func (company *Company) BeforeSave(tx *gorm.DB) (err error) {
	log.Printf("Company: %+v", company)
	// Only validate the TaxID if it's being updated
	if tx.Statement.Changed("TaxID") {
			if !company.ValidateTaxID() {
					return errors.New("company.tax_id must contain only 15 numeric characters")
			}
	}

	// Validate the Name field if it’s being updated
	if tx.Statement.Changed("Name") && company.Name == "" {
			return errors.New("company.name cannot be empty")
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
