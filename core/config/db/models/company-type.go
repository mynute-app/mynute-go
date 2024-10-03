package models

import (
	"errors"

	"gorm.io/gorm"
)

type CompanyType struct {
	gorm.Model
	Name string `gorm:"not null;unique" json:"name"`
}

func (companyType *CompanyType) BeforeCreate(tx *gorm.DB) (err error) {
	// if companyType.Name == "" {
	// 	return errors.New("companyType.Name cannot be empty")
	// }

	var name string

	switch dest := tx.Statement.Dest.(type) {
	case map[string]interface{}:
		if val, ok := dest["name"]; ok {
			name, _ = val.(string)
		}
	case *CompanyType:
		name = dest.Name
	default:
		// Handle other possible types if necessary
	}

	if name == "" {
		return errors.New("companyType.name cannot be empty")
	}

	return nil
}
