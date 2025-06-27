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
	BranchWorkSchedule []BranchWorkRange  `gorm:"foreignKey:BranchID;constraint:OnDelete:CASCADE;" json:"work_schedule"`
	TimeZone           string             `json:"timezone" gorm:"type:varchar(100);not null" ` // Time zone in IANA format (e.g., "America/New_York", "America/Sao_Paulo", etc.)
	Name               string             `gorm:"not null" json:"name"`
	Street             string             `gorm:"not null" json:"street"`
	Number             string             `gorm:"not null" json:"number"`
	Complement         string             `json:"complement"`
	Neighborhood       string             `gorm:"not null" json:"neighborhood"`
	ZipCode            string             `gorm:"not null" json:"zip_code"`
	City               string             `gorm:"not null" json:"city"`
	State              string             `gorm:"not null" json:"state"`
	Country            string             `gorm:"not null" json:"country"`
	CompanyID          uuid.UUID          `gorm:"not null;index" json:"company_id"`
	Employees          []*Employee        `gorm:"many2many:employee_branches;constraint:OnDelete:CASCADE"`              // Many-to-many relation with Employee
	Services           []*Service         `gorm:"many2many:branch_services;constraint:OnDelete:CASCADE"`                // Many-to-many relation with Service
	Appointments       []Appointment      `gorm:"foreignKey:BranchID;constraint:OnDelete:CASCADE;" json:"appointments"` // One-to-many relation with Appointment
	ServiceDensity     []ServiceDensity   `gorm:"type:jsonb" json:"service_density"`                                    // One-to-many relation with ServiceDensity
	BranchDensity      uint               `gorm:"not null;default:1" json:"branch_density"`
	Design             mJSON.DesignConfig `gorm:"type:jsonb" json:"design"`
}

func (Branch) TableName() string { return "branches" }

func (Branch) SchemaType() string { return "company" }

type ServiceDensity struct {
	ServiceID           uuid.UUID `json:"service_id"`
	MaxSchedulesOverlap uint      `json:"max_schedules_overlap"`
}

func (b *Branch) BeforeCreate(tx *gorm.DB) error {
	return nil
}

func (b *Branch) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("CompanyID") {
		return lib.Error.General.UpdatedError.WithError(errors.New("the CompanyID cannot be changed after creation"))
	}
	return nil
}

func (b *Branch) GetTimeZone() (*time.Location, error) {
	loc, err := time.LoadLocation(b.TimeZone)
	if err != nil {
		return nil, fmt.Errorf("branch (%s) has invalid timezone %s: %w", b.ID, b.TimeZone, err)
	}
	return loc, nil
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

func (b *Branch) ValidateEmployeeWorkRangeTime(tx *gorm.DB, ewr *EmployeeWorkRange) error {
	if ewr.BranchID != b.ID {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee work range branch ID %s does not match branch ID %s", ewr.BranchID, b.ID))
	}

	var bwrs []BranchWorkRange
	if err := tx.Find(&bwrs, "branch_id = ? AND weekday = ?", b.ID, ewr.Weekday).Error; err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("failed to retrieve branch work schedule: %w", err))
	}

	if len(bwrs) == 0 {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("no branch work schedule found for branch ID %s and weekday %d", b.ID, ewr.Weekday))
	}

	ewrTZ, err := ewr.GetTimeZone()
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee work range (%s) has invalid time zone %s: %w", ewr.ID, ewr.TimeZone, err))
	}

	for _, bws := range bwrs { // bws is the Branch Work Schedule
		if bws.Weekday != ewr.Weekday {
			continue // Skip if the weekday does not match
		}
		bwsTZ, err := lib.GetTimeZone(bws.TimeZone)
		if err != nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch work range (%s) has invalid time zone %s: %w", bws.ID, bws.TimeZone, err))
		}
		if lib.TimeRangeFullyContained(bws.StartTime, bws.EndTime, bwsTZ, ewr.StartTime, ewr.EndTime, ewrTZ) {
			return nil
		}
	}

	localStartTime, err := lib.Utc2LocalTime(ewr.TimeZone, ewr.StartTime)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid start time %s: %w", ewr.StartTime, err))
	}
	localEndTime, err := lib.Utc2LocalTime(ewr.TimeZone, ewr.EndTime)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid end time %s: %w", ewr.EndTime, err))
	}
	return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee work range (from %s to %s) is not within any defined branch operating hours for weekday %d", localStartTime.Format("15:04"), localEndTime.Format("15:04"), ewr.Weekday))
}

func (b *Branch) ValidateBranchWorkRangeTime(tx *gorm.DB, newRange *BranchWorkRange) error {
	var existing []BranchWorkRange

	err := tx.Find(&existing, "branch_id = ? AND weekday = ? AND id != ?", newRange.BranchID, newRange.Weekday, newRange.ID).Error
	if err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("failed to fetch existing work ranges: %w", err))
	}

	for _, bwr := range existing {
		overlap, err := newRange.Overlaps(&bwr)
		if err != nil {
			return err
		}
		if overlap {
			loc, _ := bwr.GetTimeZone()
			start := bwr.StartTime.In(loc).Format("15:04")
			end := bwr.EndTime.In(loc).Format("15:04")
			return lib.Error.General.BadRequest.WithError(fmt.Errorf(
				"branch already has a work range on %s from %s to %s that overlaps the new range %s-%s",
				bwr.Weekday.String(), start, end, newRange.StartTime, newRange.EndTime,
			))
		}
	}

	return nil
}
