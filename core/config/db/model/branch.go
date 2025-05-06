package model

import (
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/lib"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Branch model
type Branch struct {
	BaseModel
	Name           string             `gorm:"not null" json:"name"`
	Street         string             `gorm:"not null" json:"street"`
	Number         string             `gorm:"not null" json:"number"`
	Complement     string             `json:"complement"`
	Neighborhood   string             `gorm:"not null" json:"neighborhood"`
	ZipCode        string             `gorm:"not null" json:"zip_code"`
	City           string             `gorm:"not null" json:"city"`
	State          string             `gorm:"not null" json:"state"`
	Country        string             `gorm:"not null" json:"country"`
	CompanyID      uuid.UUID          `gorm:"not null;index" json:"company_id"`
	Employees      []*Employee        `gorm:"many2many:employee_branches;constraint:OnDelete:CASCADE"`              // Many-to-many relation with Employee
	Services       []*Service         `gorm:"many2many:branch_services;constraint:OnDelete:CASCADE"`                // Many-to-many relation with Service
	Appointments   []Appointment      `gorm:"foreignKey:BranchID;constraint:OnDelete:CASCADE;" json:"appointments"` // One-to-many relation with Appointment
	ServiceDensity []ServiceDensity   `gorm:"type:jsonb" json:"service_density"`                                    // One-to-many relation with ServiceDensity
	BranchDensity  uint               `gorm:"not null;default:1" json:"branch_density"`
	Design         mJSON.DesignConfig `gorm:"type:jsonb" json:"design"`
}

type ServiceDensity struct {
	ServiceID           uuid.UUID `json:"service_id"`
	MaxSchedulesOverlap uint      `json:"max_schedules_overlap"`
}

func (b *Branch) AddService(tx *gorm.DB, service *Service) error {
	if service.CompanyID != b.CompanyID {
		return lib.Error.Company.NotSame
	}
	if err := tx.Model(&BranchResource).Association("Services").Append(&service); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	if err := tx.Model(&b).Preload("Services").Where("id = ?", b.ID).First(&b).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.UpdatedError.WithError(fmt.Errorf("branch not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

func (b *Branch) RemoveService(tx *gorm.DB, service *Service) error {
	if service.CompanyID != b.CompanyID {
		return lib.Error.Company.NotSame
	}
	if err := tx.Model(&b).Association("Services").Delete(&service); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	if err := tx.Model(&b).Preload("Services").Where("id = ?", b.ID).First(&b).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.UpdatedError.WithError(fmt.Errorf("branch not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}
