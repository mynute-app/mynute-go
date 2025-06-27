package model

import (
	"agenda-kaki-go/core/lib"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Weekday uint8

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

type EmployeeWorkSchedule struct {
	WorkRanges []EmployeeWorkRange `json:"employee_work_ranges"`
}

type EmployeeWorkRange struct {
	BaseModel
	Weekday    time.Weekday `gorm:"not null" validate:"required" json:"weekday" ` // 0 = Sunday, 1 = Monday, ..., 6 = Saturday
	StartTime  time.Time    `gorm:"type:time" validate:"required" json:"start_time" `
	EndTime    time.Time    `gorm:"type:time" validate:"required" json:"end_time" `
	TimeZone   string       `gorm:"type:varchar(100)" validate:"required" json:"timezone"` // Time zone in IANA format (e.g., "America/New_York", "America/Sao_Paulo", etc.)
	EmployeeID uuid.UUID    `gorm:"type:uuid;not null;index:idx_employee_id" validate:"required" json:"employee_id"`
	Employee   Employee     `gorm:"foreignKey:EmployeeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"employee"`
	BranchID   uuid.UUID    `gorm:"type:uuid;not null;index:idx_branch_id" validate:"required" json:"branch_id"`
	Branch     Branch       `gorm:"foreignKey:BranchID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"branch"`
	Services   []*Service   `gorm:"many2many:employee_work_range_services;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"services"`
}

const WorksRangeTableName = "employee_work_ranges"

func (EmployeeWorkRange) TableName() string          { return WorksRangeTableName }
func (EmployeeWorkRange) SchemaType() string         { return "tenant" }
func (EmployeeWorkRange) Indexes() map[string]string { return WorksRangeIndexes(WorksRangeTableName) }

func WorksRangeIndexes(table string) map[string]string {
	return map[string]string{
		"idx_employee_weekday":                   fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_employee_weekday ON %s (employee_id, weekday)", table),
		"idx_employee_branch":                    fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_employee_branch_id ON %s (employee_id, branch_id)", table),
		"idx_employee_start_time":                fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_employee_start_time ON %s (employee_id, start_time)", table),
		"idx_employee_branch_start_time":         fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_employee_branch_start_time ON %s (employee_id, branch_id, start_time)", table),
		"idx_employee_branch_weekday":            fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_employee_branch_weekday ON %s (employee_id, branch_id, weekday)", table),
		"idx_employee_branch_weekday_start_time": fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_employee_branch_weekday_start_time ON %s (employee_id, branch_id, weekday, start_time)", table),
		"idx_branch_weekday":                     fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_branch_weekday ON %s (branch_id, weekday)", table),
		"idx_branch_start_time":                  fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_branch_start_time ON %s (branch_id, start_time)", table),
		"idx_branch_weekday_start_time":          fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_branch_weekday_start_time ON %s (branch_id, weekday, start_time)", table),
	}
}

func (ewr *EmployeeWorkRange) AfterFind(tx *gorm.DB) error {
	var err error

	ewr.StartTime, err = lib.Utc2LocalTime(ewr.TimeZone, ewr.StartTime)
	if err != nil {
		return fmt.Errorf("employee work range (%s) failed to convert start time to local time: %w", ewr.ID, err)
	}

	ewr.EndTime, err = lib.Utc2LocalTime(ewr.TimeZone, ewr.EndTime)
	if err != nil {
		return fmt.Errorf("employee work range (%s) failed to convert end time to local time: %w", ewr.ID, err)
	}

	return nil
}

func (ewr *EmployeeWorkRange) BeforeCreate(tx *gorm.DB) error {
	var err error
	ewr.StartTime, err = lib.LocalTime2UTC(ewr.TimeZone, ewr.StartTime)

	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid start time %s: %w", ewr.StartTime, err))
	}

	ewr.EndTime, err = lib.LocalTime2UTC(ewr.TimeZone, ewr.EndTime)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid end time %s: %w", ewr.EndTime, err))
	}

	if ewr.StartTime.Equal(ewr.EndTime) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range start time cannot be equal to end time"))
	}

	if ewr.StartTime.After(ewr.EndTime) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range start time cannot be after end time"))
	}

	if ewr.Weekday < 0 || ewr.Weekday > 6 {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid weekday %d, must be between 0 (Sunday) and 6 (Saturday)", ewr.Weekday))
	}

	var branch *Branch
	if err := tx.First(&branch, "id = ?", ewr.BranchID.String()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch with ID %s not found", ewr.BranchID))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	if !branch.HasEmployee(tx, ewr.EmployeeID) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee %s does not belong to branch %s", ewr.EmployeeID, ewr.BranchID))
	}

	if err := branch.ValidateEmployeeWorkRangeTime(tx, ewr); err != nil {
		return err
	}

	if len(ewr.Services) > 0 {
		for _, service := range ewr.Services {
			if service == nil {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("service cannot be nil"))
			}
			if !branch.HasService(tx, service.ID) {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch %s does not have service %s", ewr.BranchID, service.ID))
			}
		}
	}

	var employee *Employee
	if err := tx.First(&employee, "id = ?", ewr.EmployeeID.String()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee with ID %s not found", ewr.EmployeeID))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	if !employee.HasBranch(tx, ewr.BranchID) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee %s does not belong to branch %s", ewr.EmployeeID, ewr.BranchID))
	}

	for _, service := range ewr.Services {
		if service == nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("service cannot be nil"))
		}
		if !employee.HasService(tx, service.ID) {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee %s does not have service %s", ewr.EmployeeID, service.ID))
		}
	}

	if err := employee.ValidateEmployeeWorkRangeTime(tx, ewr); err != nil {
		return err
	}

	if ewr.StartTime.Second() != 0 {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("parsing of start time generated a time with seconds, which is not allowed: %d", ewr.StartTime.Second()))
	} else if ewr.EndTime.Second() != 0 {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("parsing of end time generated a time with seconds, which is not allowed: %d", ewr.EndTime.Second()))
	}

	return nil
}

func (ewr *EmployeeWorkRange) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("StartTime") || tx.Statement.Changed("EndTime") || tx.Statement.Changed("TimeZone") {
		var err error
		ewr.StartTime, err = lib.LocalTime2UTC(ewr.TimeZone, ewr.StartTime)

		if err != nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid start time %s: %w", ewr.StartTime, err))
		}

		ewr.EndTime, err = lib.LocalTime2UTC(ewr.TimeZone, ewr.EndTime)
		if err != nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid end time %s: %w", ewr.EndTime, err))
		}
		if ewr.StartTime.Equal(ewr.EndTime) {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range start time cannot be equal to end time"))
		} else if ewr.StartTime.After(ewr.EndTime) {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range start time cannot be after end time"))
		}
	}

	if tx.Statement.Changed("EmployeeID") {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee ID cannot be changed after creation"))
	} else if tx.Statement.Changed("BranchID") {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch ID cannot be changed after creation"))
	}

	var employee *Employee
	if err := tx.First(&employee, "id = ?", ewr.EmployeeID.String()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee with ID %s not found", ewr.EmployeeID))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	var branch *Branch
	if err := tx.First(&branch, "id = ?", ewr.BranchID.String()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch with ID %s not found", ewr.BranchID))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	if tx.Statement.Changed("Weekday") || tx.Statement.Changed("StartTime") || tx.Statement.Changed("EndTime") || tx.Statement.Changed("TimeZone") {
		if err := employee.ValidateEmployeeWorkRangeTime(tx, ewr); err != nil {
			return err
		}
		if err := branch.ValidateEmployeeWorkRangeTime(tx, ewr); err != nil {
			return err
		}
	}

	return nil
}

func (ewr *EmployeeWorkRange) HasService(tx *gorm.DB, serviceID uuid.UUID) bool {
	var count int64
	if err := tx.Raw("SELECT COUNT(*) FROM work_schedule_services WHERE work_range_id = ? AND service_id = ?", ewr.ID, serviceID).Scan(&count).Error; err != nil {
		return false // Error occurred, assume service does not exist
	}
	return count > 0
}

func (ewr *EmployeeWorkRange) GetTimeZone() (*time.Location, error) {
	loc, err := time.LoadLocation(ewr.TimeZone)
	if err != nil {
		return nil, fmt.Errorf("employee work range (%s) has invalid timezone %s: %w", ewr.ID, ewr.TimeZone, err)
	}
	return loc, nil
}

func (ewr *EmployeeWorkRange) Overlaps(other *EmployeeWorkRange) (bool, error) {
	if ewr.Weekday != other.Weekday {
		return false, nil
	}

	loc1, err := ewr.GetTimeZone()
	if err != nil {
		return false, err
	}
	loc2, err := other.GetTimeZone()
	if err != nil {
		return false, err
	}

	return lib.TimeRangeOverlaps(ewr.StartTime, ewr.EndTime, loc1, other.StartTime, other.EndTime, loc2), nil
}

func (ewr *EmployeeWorkRange) AddServices(tx *gorm.DB, services ...*Service) error {
	if ewr.Services == nil {
		ewr.Services = make([]*Service, 0)
	}
	var employee *Employee
	if err := tx.Find(&employee, "id = ?", ewr.EmployeeID).Error; err != nil {
		return fmt.Errorf("error finding employee with ID %s: %w", ewr.EmployeeID, err)
	}
	if !employee.HasBranch(tx, ewr.BranchID) {
		return fmt.Errorf("employee %s does not belong to branch %s", employee.ID, ewr.BranchID)
	}
	var branch *Branch
	if err := tx.Find(&branch, "id = ?", ewr.BranchID).Error; err != nil {
		return fmt.Errorf("error finding branch with ID %s: %w", ewr.BranchID, err)
	}
	for _, service := range services {
		if service != nil && !ewr.HasService(tx, service.ID) {
			// Check if the employee has the service
			if !employee.HasService(tx, service.ID) {
				return fmt.Errorf("employee %s does not have service %s", employee.ID, service.ID)
			}
			if !branch.HasService(tx, service.ID) {
				return fmt.Errorf("branch %s does not have service %s", branch.ID, service.ID)
			}
			if err := tx.Model(ewr).Association("Services").Append(service); err != nil {
				return fmt.Errorf("error adding service %s to work range: %w", service.ID, err)
			}
			ewr.Services = append(ewr.Services, service)
		}
	}
	return nil
}
