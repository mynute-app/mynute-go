package model

import (
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/lib"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (Branch) TableName() string { return "branches" }

func (Branch) SchemaType() string { return "company" }

type ServiceDensity struct {
	ServiceID           uuid.UUID `json:"service_id"`
	MaxSchedulesOverlap uint      `json:"max_schedules_overlap"`
}

func (b *Branch) AddService(tx *gorm.DB, service *Service) error {
	if service.CompanyID != b.CompanyID {
		return lib.Error.Company.NotSame
	}
	bID := b.ID.String()
	sID := service.ID.String()
	if err := tx.Exec("INSERT INTO branch_services (branch_id, service_id) VALUES (?, ?)", bID, sID).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	if err := tx.Preload(clause.Associations).First(&b, "id = ?", bID).Error; err != nil {
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
	bID := b.ID.String()
	sID := service.ID.String()
	if err := tx.Exec("DELETE FROM branch_services WHERE branch_id = ? AND service_id = ?", bID, sID).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	if err := tx.Preload(clause.Associations).First(&b, "id = ?", bID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.UpdatedError.WithError(fmt.Errorf("branch not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

func (b *Branch) GetAddress() string {
	number := b.Number

	if number == "" {
		switch strings.ToLower(b.Country) {
		case "brazil", "brasil":
			number = "S/N"
		case "usa", "united states", "united states of america":
			number = "N/A"
		case "spain", "españa", "mexico", "méxico":
			number = "S/N"
		default:
			number = "N/A" // fallback universal
		}
	}

	// Evita espaços e vírgulas sobrando
	parts := []string{b.Street, number}
	if b.Complement != "" {
		parts = append(parts, b.Complement)
	}
	parts = append(parts, b.Neighborhood, b.City, b.State)

	return strings.Join(parts, ", ")
}
