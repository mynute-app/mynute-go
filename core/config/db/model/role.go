package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name      string
	Description string
}

type UserRole struct {
	ID       uint `gorm:"primaryKey"`
	UserID   uint
	RoleID   uint
	TenantID uint
}

type RolePermission struct {
	ID     uint `gorm:"primaryKey"`
	RoleID uint
	Method string
	Path   string
}