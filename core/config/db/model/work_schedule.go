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

type WorkRange struct {
	BaseModel
	Weekday    time.Weekday `json:"weekday" gorm:"not null"` // 0 = Sunday, 1 = Monday, ..., 6 = Saturday
	StartTime  time.Time    `json:"start_time" gorm:"type:time;not null"`
	EndTime    time.Time    `json:"end_time" gorm:"type:time;not null"`
	TimeZone   string       `json:"timezone" gorm:"type:varchar(100);not null"` // Time zone in IANA format (e.g., "America/New_York", "America/Sao_Paulo", etc.)
	EmployeeID uuid.UUID    `json:"employee_id" gorm:"type:uuid;not null;index:idx_employee_id,unique"`
	Employee   Employee     `json:"employee" gorm:"foreignKey:EmployeeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	BranchID   uuid.UUID    `json:"branch_id" gorm:"type:uuid;not null;index:idx_branch_id,unique"`
	Branch     Branch       `json:"branch" gorm:"foreignKey:BranchID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Services   []*Service   `json:"services" gorm:"many2many:work_schedule_services;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

const WorksRangeTableName = "work_ranges"

func (WorkRange) TableName() string          { return WorksRangeTableName }
func (WorkRange) SchemaType() string         { return "tenant" }
func (WorkRange) Indexes() map[string]string { return WorksRangeIndexes(WorksRangeTableName) }

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

func (wr *WorkRange) BeforeCreate(tx *gorm.DB) error {
	if err := wr.UTC_with_Zero_YMD_Date(); err != nil {
		return err
	}

	if wr.StartTime.Equal(wr.EndTime) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range start time cannot be equal to end time"))
	}

	if wr.StartTime.After(wr.EndTime) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range start time cannot be after end time"))
	}

	if wr.Weekday < 0 || wr.Weekday > 6 {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid weekday %d, must be between 0 (Sunday) and 6 (Saturday)", wr.Weekday))
	}

	var branch *Branch
	if err := tx.First(&branch, "id = ?", wr.BranchID.String()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch with ID %s not found", wr.BranchID))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	if !branch.HasEmployee(tx, wr.EmployeeID) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee %s does not belong to branch %s", wr.EmployeeID, wr.BranchID))
	}

	if err := branch.ValidateWorkRangeTime(wr); err != nil {
		return err
	}

	if len(wr.Services) > 0 {
		for _, service := range wr.Services {
			if service == nil {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("service cannot be nil"))
			}
			if !branch.HasService(tx, service.ID) {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch %s does not have service %s", wr.BranchID, service.ID))
			}
		}
	}

	var employee *Employee
	if err := tx.First(&employee, "id = ?", wr.EmployeeID.String()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee with ID %s not found", wr.EmployeeID))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	if !employee.HasBranch(tx, wr.BranchID) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee %s does not belong to branch %s", wr.EmployeeID, wr.BranchID))
	}

	for _, service := range wr.Services {
		if service == nil {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("service cannot be nil"))
		}
		if !employee.HasService(tx, service.ID) {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee %s does not have service %s", wr.EmployeeID, service.ID))
		}
	}

	if err := employee.HasOverlappingWorkRange(tx, wr); err != nil {
		return err
	}

	return nil
}

func (wr *WorkRange) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("StartTime") || tx.Statement.Changed("EndTime") || tx.Statement.Changed("TimeZone") {
		if err := wr.UTC_with_Zero_YMD_Date(); err != nil {
			return err
		}
		if wr.StartTime.Equal(wr.EndTime) {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range start time cannot be equal to end time"))
		} else if wr.StartTime.After(wr.EndTime) {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range start time cannot be after end time"))
		}
	}

	if tx.Statement.Changed("EmployeeID") {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee ID cannot be changed after creation"))
	} else if tx.Statement.Changed("BranchID") {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch ID cannot be changed after creation"))
	}

	var employee *Employee
	if err := tx.First(&employee, "id = ?", wr.EmployeeID.String()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee with ID %s not found", wr.EmployeeID))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	var branch *Branch
	if err := tx.First(&branch, "id = ?", wr.BranchID.String()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch with ID %s not found", wr.BranchID))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	if tx.Statement.Changed("Weekday") || tx.Statement.Changed("StartTime") || tx.Statement.Changed("EndTime") || tx.Statement.Changed("TimeZone") {
		if err := employee.HasOverlappingWorkRange(tx, wr); err != nil {
			return err
		}
		if err := branch.ValidateWorkRangeTime(wr); err != nil {
			return err
		}
	}

	return nil
}

func (wr *WorkRange) GetTimeZone() (*time.Location, error) {
	loc, err := time.LoadLocation(wr.TimeZone)
	if err != nil {
		return nil, lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid time zone %s: %w", wr.TimeZone, err))
	}
	return loc, nil
}

func (wr *WorkRange) UTC_with_Zero_YMD_Date() error {
	loc, err := wr.GetTimeZone()
	if err != nil {
		return err
	}
	start := time.Date(1, 1, 1, wr.StartTime.Hour(), wr.StartTime.Minute(), wr.StartTime.Second(), 0, loc)
	end := time.Date(1, 1, 1, wr.EndTime.Hour(), wr.EndTime.Minute(), wr.EndTime.Second(), 0, loc)

	wr.StartTime = start.UTC()
	wr.EndTime = end.UTC()
	return nil
}

func (wr *WorkRange) HasService(tx *gorm.DB, serviceID uuid.UUID) bool {
	var count int64
	if err := tx.Raw("SELECT COUNT(*) FROM work_schedule_services WHERE work_range_id = ? AND service_id = ?", wr.ID, serviceID).Scan(&count).Error; err != nil {
		return false // Error occurred, assume service does not exist
	}
	return count > 0
}

func (wr *WorkRange) Overlaps(other *WorkRange) (bool, error) {
	if wr.Weekday != other.Weekday {
		return false, nil // Different weekdays or branches, no overlap
	}
	wr_loc, err := wr.GetTimeZone()
	if err != nil {
		return false, err
	}
	other_loc, err := other.GetTimeZone()
	if err != nil {
		return false, err
	}
	// Check if the time ranges overlap considering the time zone
	wrStart := wr.StartTime.In(wr_loc)
	wrEnd := wr.EndTime.In(wr_loc)
	otherStart := other.StartTime.In(other_loc)
	otherEnd := other.EndTime.In(other_loc)

	wrStart_before_or_equal_otherStart := wrStart.Before(otherStart) || wrStart.Equal(otherStart) // |<wrStart> <otherStart> || <wrStart otherStart>|
	wrEnd_after_or_equal_otherEnd := wrEnd.After(otherEnd) || wrEnd.Equal(otherEnd)               // |<wrEnd> <otherEnd> || <wrEnd otherEnd>|
	otherStart_before_or_equal_wrStart := otherStart.Before(wrStart) || otherStart.Equal(wrStart) // |<otherStart> <wrStart> || <otherStart wrStart>|
	otherEnd_after_or_equal_wrEnd := otherEnd.After(wrEnd) || otherEnd.Equal(wrEnd)               // |<otherEnd> <wrEnd> || <otherEnd wrEnd>|

	wr_equals_other := wrStart.Equal(otherStart) && wrEnd.Equal(otherEnd)                          // |<wrStart wrEnd> <otherStart otherEnd>|
	wr_contains_other_fully := wrStart_before_or_equal_otherStart && wrEnd_after_or_equal_otherEnd // |<wrStart> <otherStart> <otherEnd> <wrEnd>|
	other_contains_wr_fully := otherStart_before_or_equal_wrStart && otherEnd_after_or_equal_wrEnd // |<otherStart> <wrStart> <wrEnd> <otherEnd>|
	wr_contains_other_start := wrStart_before_or_equal_otherStart && otherEnd_after_or_equal_wrEnd // |<wrStart> <otherStart> <wrEnd> <otherEnd>|
	other_contains_wr_start := otherStart_before_or_equal_wrStart && wrEnd_after_or_equal_otherEnd // |<otherStart> <wrStart> <otherEnd> <wrEnd>|

	isContained := wr_equals_other || wr_contains_other_fully || other_contains_wr_fully || wr_contains_other_start || other_contains_wr_start

	return isContained, nil
}

func (wr *WorkRange) AddServices(tx *gorm.DB, services ...*Service) error {
	if wr.Services == nil {
		wr.Services = make([]*Service, 0)
	}
	var employee *Employee
	if err := tx.Find(&employee, "id = ?", wr.EmployeeID).Error; err != nil {
		return fmt.Errorf("error finding employee with ID %s: %w", wr.EmployeeID, err)
	}
	if !employee.HasBranch(tx, wr.BranchID) {
		return fmt.Errorf("employee %s does not belong to branch %s", employee.ID, wr.BranchID)
	}
	var branch *Branch
	if err := tx.Find(&branch, "id = ?", wr.BranchID).Error; err != nil {
		return fmt.Errorf("error finding branch with ID %s: %w", wr.BranchID, err)
	}
	for _, service := range services {
		if service != nil && !wr.HasService(tx, service.ID) {
			// Check if the employee has the service
			if !employee.HasService(tx, service.ID) {
				return fmt.Errorf("employee %s does not have service %s", employee.ID, service.ID)
			}
			if !branch.HasService(tx, service.ID) {
				return fmt.Errorf("branch %s does not have service %s", branch.ID, service.ID)
			}
			if err := tx.Model(wr).Association("Services").Append(service); err != nil {
				return fmt.Errorf("error adding service %s to work range: %w", service.ID, err)
			}
			wr.Services = append(wr.Services, service)
		}
	}
	return nil
}
