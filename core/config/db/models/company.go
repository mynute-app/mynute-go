package models

import (
    "errors"
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
        }
        if val, ok := dest["tax_id"]; ok {
            taxID, _ = val.(string)
        }
    case *Company:
        // When creating or updating with a struct
        name = dest.Name
        taxID = dest.TaxID
    default:
        // Handle other possible types if necessary
    }

    // Validate the Name field if it's being updated
    if tx.Statement.Changed("Name") {
        if name == "" {
            return errors.New("company.name cannot be empty")
        }
    }

    // Validate the TaxID field if it's being updated
    if tx.Statement.Changed("TaxID") {
        if !validateTaxID(taxID) {
            return errors.New("company.tax_id must contain exactly 14 numeric characters")
        }
    }

    return nil
}

// validateTaxID checks if the TaxID is a 14-character string containing only numbers
func validateTaxID(taxID string) bool {
    // Define the regular expression to match exactly 14 digits
    re := regexp.MustCompile(`^\d{14}$`)
    // Check if the TaxID matches the regular expression
    return re.MatchString(taxID)
}
