package model

import "gorm.io/gorm"

type Route struct {
	gorm.Model
	Handler     string `gorm:"type:varchar(255)"`
	Description string `gorm:"type:varchar(255)"`
	Method      string `gorm:"type:varchar(10)"`
	Path        string `gorm:"type:varchar(255)"`
	IsPublic    bool   `gorm:"default:false"`
}

// Custom Composite Index
func (Route) TableName() string {
	return "routes"
}

func (Route) Indexes() map[string]string {
	return map[string]string{
		"idx_method_path": "CREATE UNIQUE INDEX idx_method_path ON routes (method, path)",
	}
}
