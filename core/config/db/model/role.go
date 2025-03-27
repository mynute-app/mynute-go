package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string
	Description string
	IsDefault   bool `gorm:"default:false"`
	IsBusiness  bool `gorm:"default:true"`
	CompanyID   *uint
	Routes      []*Route `gorm:"many2many:role_routes;constraint:OnDelete:CASCADE"`
}
