package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string
	Description string
	IsDefault   bool `gorm:"default:false"`
	CompanyID   uint
	Routes      []*Route `gorm:"many2many:role_routes;constraint:OnDelete:CASCADE"`
}