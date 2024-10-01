package models

import (
	"errors"

	"gorm.io/gorm"
)

type CompanyType struct {
	gorm.Model
	Name string `gorm:"not null;unique;required" json:"name"`
}

func (company *CompanyType) BeforeCreate(tx *gorm.DB) (err error) {
	if company.Name == "" {
		return errors.New("name cannot be empty")
	}
	return nil
}