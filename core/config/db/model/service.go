package model

import (
	"github.com/google/uuid"
)

// Third step: Choosing the service.
type Service struct {
	BaseModel
	Name        string      `gorm:"not null" json:"name"`
	Description string      `gorm:"not null" json:"description"`
	Price       int64       `gorm:"not null" json:"price"`
	Currency    string      `gorm:"type:varchar(3);not null;default:'BRL'" json:"currency"` // Default currency is BRL
	Duration    uint        `gorm:"not null" json:"duration"`                               // Duration in minutes
	CompanyID   uuid.UUID   `gorm:"not null;index" json:"company_id"`
	Company     *Company    `gorm:"foreignKey:CompanyID;references:ID;constraint:OnDelete:CASCADE;"`
	Employees   []*Employee `gorm:"many2many:employee_services;constraint:OnDelete:CASCADE;" json:"employees"` // Many-to-many relation with Employee
	Branches    []*Branch   `gorm:"many2many:branch_services;constraint:OnDelete:CASCADE;" json:"branches"`    // Many-to-many relation with Branch
}

func (Service) TableName() string  { return "services" }
func (Service) SchemaType() string { return "company" }
