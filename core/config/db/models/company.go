package models

import (
	"errors"
	"gorm.io/gorm"
	"regexp"
)



// Company holds an array of CompanyTypes.
type Company struct {
	gorm.Model
	Name  string        `gorm:"not null;unique" json:"name"`
	Types []CompanyType `gorm:"many2many:company_company_types" json:"company_types"` // Many-to-many relation with a custom join table
	TaxID string        `gorm:"not null;unique" json:"tax_id"`
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
