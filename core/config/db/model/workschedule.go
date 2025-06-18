package model

import (
	"agenda-kaki-go/core/lib"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	Weekday    time.Weekday  `json:"weekday" gorm:"not null"` // 0 = Sunday, 1 = Monday, ..., 6 = Saturday
	StartTime  time.Time     `json:"start_time" gorm:"type:time;not null"`
	EndTime    time.Time     `json:"end_time" gorm:"type:time;not null"`
	TimeZone   time.Location `json:"timezone" gorm:"not null"` // Timezone in IANA format, e.g., "America/New_York"
	EmployeeID uuid.UUID     `json:"employee_id" gorm:"type:uuid;not null;index:idx_employee_id,unique"`
	Employee   Employee      `json:"employee" gorm:"foreignKey:EmployeeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	BranchID   uuid.UUID     `json:"branch_id" gorm:"type:uuid;not null;index:idx_branch_id,unique"`
	Branch     Branch        `json:"branch" gorm:"foreignKey:BranchID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Services   []*Service    `json:"services" gorm:"many2many:work_schedule_services;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

var WorksRangeTableName = "work_ranges"

func (WorkRange) TableName() string {
	return WorksRangeTableName
}

func (WorkRange) Indexes() map[string]string {
	return WorksRangeIndexes(WorksRangeTableName)
}

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
	wr.UTC_with_Zero_YMD_Date()

	if wr.StartTime.Equal(wr.EndTime) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range start time cannot be equal to end time"))
	}

	if wr.StartTime.After(wr.EndTime) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range start time cannot be after end time"))
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

	if wr.StartTime.Before(branch.StartTime) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range start time %s cannot be before branch start time %s", wr.StartTime, branch.StartTime))
	} else if wr.StartTime.After(branch.EndTime) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range start time %s cannot be after branch end time %s", wr.StartTime, branch.EndTime))
	} else if wr.EndTime.Before(branch.StartTime) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range end time %s cannot be before branch start time %s", wr.EndTime, branch.StartTime))
	} else if wr.EndTime.After(branch.EndTime) {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range end time %s cannot be after branch end time %s", wr.EndTime, branch.EndTime))
	}

	if len(wr.Services) > 0 {
		for _, service := range wr.Services {
			if service == nil {
				continue
			}
			if !branch.HasService(tx, service.ID) {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch %s does not have service %s", wr.BranchID, service.ID))
			}
		}
	}

	return nil
}

func (wr *WorkRange) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("StartTime") || tx.Statement.Changed("EndTime") || tx.Statement.Changed("TimeZone") {
		wr.UTC_with_Zero_YMD_Date()
		if wr.StartTime.Equal(wr.EndTime) {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range start time cannot be equal to end time"))
		} else if wr.StartTime.After(wr.EndTime) {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range start time cannot be after end time"))
		}
	}
	return nil
}

func (wr *WorkRange) UTC_with_Zero_YMD_Date() {
	loc := &wr.TimeZone
	start := time.Date(0, 1, 1, wr.StartTime.Hour(), wr.StartTime.Minute(), wr.StartTime.Second(), 0, loc)
	end := time.Date(0, 1, 1, wr.EndTime.Hour(), wr.EndTime.Minute(), wr.EndTime.Second(), 0, loc)

	wr.StartTime = start.UTC()
	wr.EndTime = end.UTC()
}

func (wr *WorkRange) HasService(tx *gorm.DB, serviceID uuid.UUID) bool {
	var count int64
	if err := tx.Raw("SELECT COUNT(*) FROM work_schedule_services WHERE work_range_id = ? AND service_id = ?", wr.ID, serviceID).Scan(&count).Error; err != nil {
		return false // Error occurred, assume service does not exist
	}
	return count > 0
}

func (wr *WorkRange) AddServices(tx *gorm.DB, services ...*Service) error {
	if wr.Services == nil {
		wr.Services = make([]*Service, 0)
	}
	var employee *Employee
	if err := tx.Association(clause.Associations).Find(&employee, "id = ?", wr.EmployeeID); err != nil {
		return fmt.Errorf("error finding employee with ID %s: %w", wr.EmployeeID, err)
	}
	if !employee.HasBranch(tx, wr.BranchID) {
		return fmt.Errorf("employee %s does not belong to branch %s", employee.ID, wr.BranchID)
	}
	for _, service := range services {
		if service != nil && !wr.HasService(tx, service.ID) {
			// Check if the employee has the service
			if !employee.HasService(tx, service.ID) {
				return fmt.Errorf("employee %s does not have service %s", employee.ID, service.ID)
			}
			if err := tx.Model(wr).Association("Services").Append(service); err != nil {
				return fmt.Errorf("error adding service %s to work range: %w", service.ID, err)
			}
			wr.Services = append(wr.Services, service)
		}
	}
	return nil
}
