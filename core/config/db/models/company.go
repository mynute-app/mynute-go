package models

import (
	"agenda-kaki-go/core/lib"
	"errors"
	"log"

	"gorm.io/gorm"
)

// Company holds an array of CompanyTypes.
type Company struct {
	gorm.Model
	Name         string        `gorm:"not null;unique" json:"name"`
	CompanyTypes []CompanyType `gorm:"many2many:company_company_types;ForeignKey:id;References:id" json:"company_types"` // Many-to-many relation with a custom join table
	TaxID        string        `gorm:"not null;unique" json:"tax_id"`
}

func (company *Company) BeforeCreate(tx *gorm.DB) (err error) {
    log.Printf("BeforeCreate hook called")
	if !lib.ValidateName(company.Name) {
		return errors.New("company.name must be at least 3 characters long")
	} else if !lib.ValidateTaxID(company.TaxID) {
		return errors.New("company.tax_id must contain exactly 14 numeric characters")
	}
	return nil
}

// BeforeSave is a GORM hook that runs before the record is saved
func (company *Company) BeforeUpdate(tx *gorm.DB) (err error) {
    log.Printf("BeforeUpdate hook called")

	var name string
	var taxID string

	// Using tx.Statement.Dest instead of the company struct to handle the conditions.
	// This approach ensures that the hook works consistently whether we are creating or
	// updating records using a struct or a map. This is necessary since when updating
	// using a map, the company struct will not be loaded by GORM.
	// Therefore, the only trustable source of data is the tx.Statement.Dest field
	// for all conditions (creating, updating with a struct, and updating with a map).
	switch dest := tx.Statement.Dest.(type) {
	case map[string]interface{}:
		// When updating with a map
		if val, ok := dest["name"]; ok {
			name, _ = val.(string)
			if !lib.ValidateName(name) {
				return errors.New("company.name must be at least 3 characters long")
			}
		} else if val, ok := dest["tax_id"]; ok {
			taxID, _ = val.(string)
			if !lib.ValidateTaxID(taxID) {
				return errors.New("company.tax_id must contain exactly 14 numeric characters")
			}
		}
	}

	return nil
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
