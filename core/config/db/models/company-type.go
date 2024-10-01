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
	if companyType.Name == "" {
		return errors.New("companyType.Name cannot be empty")
	}
	return nil
}