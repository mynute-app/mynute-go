package models

import (
	"gorm.io/gorm"
)

type CompanyType struct {
	gorm.Model
	Name string `gorm:"not null;unique" json:"name"`
}