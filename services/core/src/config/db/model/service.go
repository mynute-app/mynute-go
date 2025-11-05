package model

import (
	"errors"
	mJSON "mynute-go/services/core/src/config/db/model/json"
	"mynute-go/services/core/src/lib"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Third step: Choosing the service.
type Service struct {
	BaseModel
	Name        string             `gorm:"type:varchar(100)" validate:"required,min=3,max=100" json:"name"`
	Description string             `gorm:"type:text" validate:"required,min=3,max=1000" json:"description"`
	Price       int64              `gorm:"not null" validate:"required,min=0" json:"price"`
	Currency    string             `gorm:"type:varchar(3);default:'BRL'" json:"currency"` // Default currency is BRL
	Duration    uint16             `gorm:"not null" json:"duration"`                      // In minutes                    // Duration in minutes
	CompanyID   uuid.UUID          `gorm:"not null;index" json:"company_id"`
	Company     *Company           `gorm:"foreignKey:CompanyID;references:ID;constraint:OnDelete:CASCADE;" json:"company"`
	Employees   []*Employee        `gorm:"many2many:employee_services;constraint:OnDelete:CASCADE;" json:"employees"` // Many-to-many relation with Employee
	Branches    []*Branch          `gorm:"many2many:branch_services;constraint:OnDelete:CASCADE;" json:"branches"`    // Many-to-many relation with Branch
	Design      mJSON.DesignConfig `gorm:"type:jsonb" json:"design"`
}

const ServiceTableName = "services"

func (Service) TableName() string  { return ServiceTableName }
func (Service) SchemaType() string { return "company" }

func (s *Service) BeforeUpdate(tx *gorm.DB) (err error) {
	// Check if CompanyID is being changed
	if tx.Statement.Changed("CompanyID") {
		return lib.Error.General.UpdatedError.WithError(errors.New("the CompanyID cannot be changed after creation"))
	}
	return nil
}
