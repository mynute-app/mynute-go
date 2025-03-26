package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string
	Description string
	IsDefault   bool `gorm:"default:false"`
	CompanyID   *uint
}

type UserRole struct {
	gorm.Model
	UserID    uint
	RoleID    uint
	CompanyID uint
}

type RolePermission struct {
	ID      uint `gorm:"primaryKey"`
	RoleID  uint
	RouteID uint
	Route   Route `gorm:"foreignKey:RouteID"`
}
