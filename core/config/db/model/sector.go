package model

import "gorm.io/gorm"

type Sector struct {
	gorm.Model
	Name        string `gorm:"not null;unique" json:"name"`
	Description string `json:"description"`
}
