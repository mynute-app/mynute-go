package model

import "gorm.io/gorm"

type Route struct {
	gorm.Model
	Handler     string
	Description string
	Method      string
	Path        string
	IsPublic    bool
}
