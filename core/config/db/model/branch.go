package model

import (
	"errors"
	"fmt"
	mJSON "mynute-go/core/config/db/model/json"
	"mynute-go/core/lib"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Branch model
type Branch struct {
	BaseModel
	Name                string                 `gorm:"type:varchar(100)" validate:"required,min=3,max=100" json:"name"` // Branch name
	Street              string                 `gorm:"type:varchar(100)" validate:"required,min=3,max=100" json:"street"`
	Number              string                 `gorm:"type:varchar(100)" validate:"required,min=1,max=10" json:"number"`
	Complement          string                 `gorm:"type:varchar(100)" validate:"max=100" json:"complement"`
	Neighborhood        string                 `gorm:"type:varchar(100)" validate:"required,min=3,max=100" json:"neighborhood"`
	ZipCode             string                 `gorm:"type:varchar(100)" validate:"required,min=5,max=8" json:"zip_code"`
	City                string                 `gorm:"type:varchar(100)" validate:"required,min=3,max=100" json:"city"`
	State               string                 `gorm:"type:varchar(100)" validate:"required,min=2,max=20" json:"state"` // State code (e.g., "NY" for New York, "SP" for São Paulo)
	Country             string                 `gorm:"type:varchar(100)" validate:"required,min=3,max=100" json:"country"`
	CompanyID           uuid.UUID              `gorm:"not null;index" json:"company_id"`
	Employees           []*Employee            `gorm:"many2many:employee_branches;constraint:OnDelete:CASCADE"`              // Many-to-many relation with Employee
	Services            []*Service             `gorm:"many2many:branch_services;constraint:OnDelete:CASCADE"`                // Many-to-many relation with Service
	Appointments        []Appointment          `gorm:"foreignKey:BranchID;constraint:OnDelete:CASCADE;" json:"appointments"` // One-to-many relation with Appointment
	ServiceDensity      []BranchServiceDensity `gorm:"foreignKey:BranchID;constraint:OnDelete:CASCADE;" json:"service_density"`
	WorkSchedule        []BranchWorkRange      `gorm:"foreignKey:BranchID;constraint:OnDelete:CASCADE;" json:"work_schedule"`
	TimeZone            string                 `gorm:"type:varchar(100)" json:"time_zone" validate:"required,myTimezoneValidation"` // Time zone in IANA format (e.g., "America/New_York", "America/Sao_Paulo", etc.)
	TotalServiceDensity int32                  `gorm:"not null;default:-1" json:"total_service_density"`
	Design              mJSON.DesignConfig     `gorm:"type:jsonb" json:"design"`
}

func (Branch) TableName() string { return "branches" }

func (Branch) SchemaType() string { return "company" }

type ServiceDensity struct {
	ServiceID           uuid.UUID `json:"service_id"`
	MaxSchedulesOverlap uint      `json:"max_schedules_overlap"`
}

func (b *Branch) BeforeCreate(tx *gorm.DB) error {
	if err := lib.MyCustomStructValidator(b); err != nil {
		return err
	}
	return nil
}

func (b *Branch) BeforeUpdate(tx *gorm.DB) error {
	if b.CompanyID != uuid.Nil {
		return lib.Error.General.UpdatedError.WithError(errors.New("the CompanyID cannot be changed after creation"))
	}
	var serviceDensity []BranchServiceDensity
	if err := tx.Find(&serviceDensity, "branch_id = ?", b.ID).Error; err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("error loading branch service densities: %w", err))
	}
	if b.TotalServiceDensity > 0 {
		for _, sd := range serviceDensity {
			if sd.Density > b.TotalServiceDensity {
				return lib.Error.Branch.MaxServiceCapacityReached.WithError(fmt.Errorf("existing service density (%d) from service (%s) exceeding branch density value to update (%d)", sd.Density, sd.ServiceID, b.TotalServiceDensity))
			}
		}
	}
	return nil
}

func (b *Branch) GetTimeZone() (*time.Location, error) {
	loc, err := time.LoadLocation(b.TimeZone)
	if err != nil {
		return nil, fmt.Errorf("branch (%s) has invalid time_zone %s: %w", b.ID, b.TimeZone, err)
	}
	return loc, nil
}

func (b *Branch) AddService(tx *gorm.DB, service *Service) error {
	if service.CompanyID != b.CompanyID {
		return lib.Error.Company.NotSame
	}
	bID := b.ID.String()
	sID := service.ID.String()
	// Check if the service already exists in the branch
	var count int64
	if err := tx.Raw("SELECT COUNT(*) FROM branch_services WHERE branch_id = ? AND service_id = ?", bID, sID).Scan(&count).Error; err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("failed to check if service %s is already in branch %s: %w", service.ID, bID, err))
	}
	if count > 0 {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch %s already has service %s", bID, sID))
	}
	if err := tx.Exec("INSERT INTO branch_services (branch_id, service_id) VALUES (?, ?)", bID, sID).Error; err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("failed to add service %s to branch %s: %w", service.ID, bID, err))
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
	// Check if the service exists in the branch
	var count int64
	if err := tx.Raw("SELECT COUNT(*) FROM branch_services WHERE branch_id = ? AND service_id = ?", bID, sID).Scan(&count).Error; err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("failed to check if service %s is in branch %s: %w", service.ID, bID, err))
	}
	if count == 0 {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch %s does not have service %s", bID, sID))
	}
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

func (b *Branch) HasServices(tx *gorm.DB, services []*Service) error {
	if len(services) > 0 {
		for _, service := range services {
			if service == nil {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("service passed is nil when validating branch (%s) services", b.ID))
			}
			if err := b.HasService(tx, service); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Branch) HasService(tx *gorm.DB, service *Service) error {
	if service == nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("service passed is nil when validating branch (%s) services", b.ID))
	}
	var count int64
	// Check if the service exists in the branch
	if err := tx.Raw("SELECT COUNT(*) FROM branch_services WHERE branch_id = ? AND service_id = ?", b.ID, service.ID).Scan(&count).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	if count == 0 {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch %s does not have service %s", b.ID, service.ID))
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

func (b *Branch) HasEmployee(tx *gorm.DB, employee *Employee) error {
	if employee == nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee passed is nil when validating branch (%s) employees", b.ID))
	}
	var count int64
	// Check if the employee exists in the branch
	if err := tx.Raw("SELECT COUNT(*) FROM employee_branches WHERE branch_id = ? AND employee_id = ?", b.ID, employee.ID).Scan(&count).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	if count == 0 {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee %s not found in branch %s", employee.ID, b.ID))
	}
	return nil
}

// ValidateEmployeeWorkRangeTime checks if the employee work range is within the branch's operating hours for the specified weekday.
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
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee work range has invalid time zone %s: %w", ewrTZ, err))
	}

	for _, bws := range bwrs { // bws is the Branch Work Schedule
		if bws.Weekday != ewr.Weekday {
			continue // Skip if the weekday does not match
		}
		bwsTZ, err := bws.GetTimeZone()
		if err != nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch work range (%s) has invalid time zone: %w", bws.ID, err))
		}
		if lib.TimeRangeFullyContained(bws.StartTime, bws.EndTime, bwsTZ, ewr.StartTime, ewr.EndTime, ewrTZ) {
			return nil
		}
	}
	ewrTimeZoneStr, err := ewr.GetTimeZoneString()
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee work range has invalid time zone %s: %w", ewrTimeZoneStr, err))
	}
	localStartTime, err := lib.Utc2LocalTime(ewrTimeZoneStr, ewr.StartTime)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid start time %s: %w", ewr.StartTime, err))
	}
	localEndTime, err := lib.Utc2LocalTime(ewrTimeZoneStr, ewr.EndTime)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid end time %s: %w", ewr.EndTime, err))
	}
	return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee work range (from %s to %s) is not within any defined branch operating hours for weekday %d", localStartTime.Format("15:04"), localEndTime.Format("15:04"), ewr.Weekday))
}

// ValidateBranchWorkRangeTime checks if the new work range overlaps with existing ones for the branch.
func (b *Branch) ValidateBranchWorkRangeTime(tx *gorm.DB, newRange *BranchWorkRange) error {
	var existing []BranchWorkRange

	err := tx.
		Where("branch_id = ? AND weekday = ? AND id != ?", newRange.BranchID, newRange.Weekday, newRange.ID).
		Where("start_time <= ? AND end_time >= ?", newRange.EndTime, newRange.StartTime).
		Find(&existing).Error
	if err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("failed to fetch existing work ranges: %w", err))
	}

	for _, bwr := range existing {
		overlap, err := newRange.WorkRangeBase.Overlaps(&bwr.WorkRangeBase)
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
