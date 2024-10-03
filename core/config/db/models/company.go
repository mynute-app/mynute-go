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

// BeforeCreate is a GORM hook that runs before a new record is created
func (company *Company) BeforeCreate(tx *gorm.DB) (err error) {
	log.Printf("We are creating a company record")
	log.Printf("Company: %+v", company)
	if !validateCompanyName(company.Name) {
		return errors.New("company.name must be at least 3 characters long")
	} else if !validateTaxID(company.TaxID) {
		return errors.New("company.tax_id must contain exactly 14 numeric characters")
	}
	return nil
}

func (company *Company) BeforeUpdate(tx *gorm.DB) (err error) {
	log.Printf("We are updating a company record")
	log.Printf("Company: %+v", company)
	log.Printf(`tx.Statement.Changed("Name"): %+v`, tx.Statement.Changed("Name"))
	log.Printf(`tx.Statement.Changed("TaxID"): %+v`, tx.Statement.Changed("TaxID"))
	if tx.Statement.Changed("Name") && !validateCompanyName(company.Name) {
		return errors.New("company.name must be at least 3 characters long")
	} else if tx.Statement.Changed("TaxID") && !validateTaxID(company.TaxID) {
		return errors.New("company.tax_id must contain exactly 14 numeric characters")
	}
	return nil
}

// validateTaxID checks if the TaxID is a 14-character string containing only numbers
func validateTaxID(taxID string) bool {
	return regexp.MustCompile(`^\d{14}$`).MatchString(taxID)
}

func validateCompanyName(name string) bool {
	return len(name) >= 3
}

// BeforeUpdate is a GORM hook that runs before the record is updated
// func (company *Company) BeforeUpdate(tx *gorm.DB) (err error) {
// 	log.Printf("We are updating a company record")
// 	var name string
// 	var taxID string

// 	// Using tx.Statement.Dest instead of the company struct to handle the conditions.
// 	// This approach ensures that the hook works consistently whether we are
// 	// updating records using a struct or a map. This is necessary since when updating
// 	// using a map, the company struct will not be loaded by GORM.
// 	// Therefore, the only trustable source of data is the tx.Statement.Dest field
// 	// for all conditions (updating with a struct or updating with a map).
// 	switch dest := tx.Statement.Dest.(type) {
// 	case map[string]interface{}:
// 		// When updating with a map
// 		if val, ok := dest["name"]; ok {
// 			name, _ = val.(string)
// 		}
// 		if val, ok := dest["tax_id"]; ok {
// 			taxID, _ = val.(string)
// 		}
// 	case *Company:
// 		// When creating or updating with a struct
// 		name = dest.Name
// 		taxID = dest.TaxID
// 	default:
// 		// Handle other possible types if necessary
// 	}

// 	// Validate the Name field if it's being updated
// 	if len(name) < 3 {
// 		return errors.New("company.name must be at least 3 characters long")
// 	}

// 	// Validate the TaxID field if it's being updated
// 	if !validateTaxID(taxID) {
// 		return errors.New("company.tax_id must contain exactly 14 numeric characters")
// 	}

// 	return nil
// }
