package model

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	EmployeeID   uint
	ResourceType string
	EndPointID   *uint // nil means all resources of this type
	Action       string
}
