package model

import "gorm.io/gorm"

type PolicyRule struct {
	gorm.Model
	CompanyID     uint `gorm:"index"`
	UserCreatorID uint `gorm:"index"`
	Name          string 
	Description   string
	SubjectAttr   string
	SubjectValue  string
	ResourceAttr  string
	ResourceValue string
	AttrCondition string
	Method        string
	Path          string
}
