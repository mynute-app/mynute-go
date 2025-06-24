package model

import (
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/lib"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Branch model
type Branch struct {
	BaseModel
	StartTime      time.Time          `json:"start_time" gorm:"type:time;not null"`
	EndTime        time.Time          `json:"end_time" gorm:"type:time;not null"`
	TimeZone       string             `json:"timezone" gorm:"type:varchar(100);not null"` // Time zone in IANA format (e.g., "America/New_York", "America/Sao_Paulo", etc.)
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

func (b *Branch) BeforeCreate(tx *gorm.DB) error {
	if err := b.UTC_with_Zero_YMD_Date(); err != nil {
		return err
	}
	return nil
}

func (b *Branch) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("CompanyID") {
		return lib.Error.General.UpdatedError.WithError(errors.New("the CompanyID cannot be changed after creation"))
	}
	if tx.Statement.Changed("StartTime") || tx.Statement.Changed("EndTime") || tx.Statement.Changed("TimeZone") {
		if err := b.UTC_with_Zero_YMD_Date(); err != nil {
			return err
		}
		if b.StartTime.Equal(b.EndTime) {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch start time cannot be equal to end time"))
		} else if b.StartTime.Before(b.EndTime) {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch start time cannot be before end time"))
		}
	}
	return nil
}

func (b *Branch) GetTimeZone() (*time.Location, error) {
	loc, err := time.LoadLocation(b.TimeZone)
	if err != nil {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid time zone %s: %w", b.TimeZone, err))
	}
	return loc, nil
}

func (b *Branch) UTC_with_Zero_YMD_Date() error {
	loc, err := time.LoadLocation(b.TimeZone)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid time zone %s: %w", b.TimeZone, err))
	}
	start := time.Date(1, 1, 1, b.StartTime.Hour(), b.StartTime.Minute(), b.StartTime.Second(), 0, loc)
	end := time.Date(1, 1, 1, b.EndTime.Hour(), b.EndTime.Minute(), b.EndTime.Second(), 0, loc)

	b.StartTime = start.UTC()
	b.EndTime = end.UTC()
	return nil
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

func (b *Branch) HasEmployee(tx *gorm.DB, employeeID uuid.UUID) bool {
	var count int64
	// Check if the employee exists in the branch
	if err := tx.Raw("SELECT COUNT(*) FROM employee_branches WHERE branch_id = ? AND employee_id = ?", b.ID, employeeID).Scan(&count).Error; err != nil {
		return false // Error occurred, assume employee does not exist
	}
	return count > 0
}

func (b *Branch) HasService(tx *gorm.DB, serviceID uuid.UUID) bool {
	var count int64
	// Check if the service exists in the branch
	if err := tx.Raw("SELECT COUNT(*) FROM branch_services WHERE branch_id = ? AND service_id = ?", b.ID, serviceID).Scan(&count).Error; err != nil {
		return false // Error occurred, assume service does not exist
	}
	return count > 0
}

func (b *Branch) ValidateWorkRangeTime(wr *WorkRange) error {
	if wr.StartTime.Before(b.StartTime) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range start time %s cannot be before branch start time %s", wr.StartTime, b.StartTime))
	} else if wr.StartTime.After(b.EndTime) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range start time %s cannot be after branch end time %s", wr.StartTime, b.EndTime))
	} else if wr.EndTime.Before(b.StartTime) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range end time %s cannot be before branch start time %s", wr.EndTime, b.StartTime))
	} else if wr.EndTime.After(b.EndTime) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range end time %s cannot be after branch end time %s", wr.EndTime, b.EndTime))
	}
	return nil
}
