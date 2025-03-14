package model

import "gorm.io/gorm"

// Branch model
type Branch struct {
	gorm.Model
	Name         string     `gorm:"not null" json:"name"`
	CompanyID    uint       `gorm:"not null" json:"company_id"`  // Foreign key to Company
	Employees    []*Employee `gorm:"many2many:branch_employees;constraint:OnDelete:CASCADE"` // Many-to-many relation with Employee
	Services     []*Service  `gorm:"many2many:branch_services;constraint:OnDelete:CASCADE"`  // Many-to-many relation with Service
	Street       string     `gorm:"not null" json:"street"`
	Number       string     `gorm:"not null" json:"number"`
	Complement   string     `json:"complement"`
	Neighborhood string     `gorm:"not null" json:"neighborhood"`
	ZipCode      string     `gorm:"not null" json:"zip_code"`
	City         string     `gorm:"not null" json:"city"`
	State        string     `gorm:"not null" json:"state"`
	Country      string     `gorm:"not null" json:"country"`
}
