package model

import (
	"fmt"
	"mynute-go/core/src/lib"

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
	WorkRangeBase
	EmployeeID uuid.UUID  `gorm:"type:uuid;not null;index:idx_employee_id" validate:"required" json:"employee_id"`
	Employee   Employee   `gorm:"foreignKey:EmployeeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"employee" validate:"-"`
	Services   []*Service `json:"services" gorm:"many2many:employee_work_range_services;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" validate:"-"`
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

func (ewr *EmployeeWorkRange) BeforeCreate(tx *gorm.DB) error {
	if err := ewr.WorkRangeBase.BeforeCreate(tx); err != nil {
		return err
	}

	branch := &Branch{BaseModel: BaseModel{ID: ewr.BranchID}}
	employee := &Employee{UserID: ewr.EmployeeID}

	if err := branch.HasEmployee(tx, employee); err != nil {
		return err
	}

	if err := employee.HasBranch(tx, ewr.BranchID); err != nil {
		return err
	}

	if err := branch.HasServices(tx, ewr.Services); err != nil {
		return err
	}

	if err := employee.HasServices(tx, ewr.Services); err != nil {
		return err
	}

	if err := branch.ValidateEmployeeWorkRangeTime(tx, ewr); err != nil {
		return err
	}

	if err := employee.ValidateEmployeeWorkRangeTime(tx, ewr); err != nil {
		return err
	}

	return nil
}

func (ewr *EmployeeWorkRange) BeforeUpdate(tx *gorm.DB) error {
	var old EmployeeWorkRange
	if err := tx.Model(&old).Where("id = ?", ewr.ID).First(&old).Error; err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("error fetching old work range: %w", err))
	}

	if err := ewr.WorkRangeBase.BeforeUpdate(tx); err != nil {
		return err
	}

	if ewr.EmployeeID != old.EmployeeID {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee ID cannot be changed after creation"))
	} else if ewr.BranchID != old.BranchID {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch ID cannot be changed after creation"))
	}

	branch := &Branch{BaseModel: BaseModel{ID: ewr.BranchID}}
	employee := &Employee{UserID: ewr.EmployeeID}

	if err := employee.ValidateEmployeeWorkRangeTime(tx, ewr); err != nil {
		return err
	}

	if err := branch.ValidateEmployeeWorkRangeTime(tx, ewr); err != nil {
		return err
	}

	return nil
}

func (ewr *EmployeeWorkRange) HasService(tx *gorm.DB, serviceID uuid.UUID) error {
	var count int64
	if err := tx.Raw("SELECT COUNT(*) FROM work_schedule_services WHERE work_range_id = ? AND service_id = ?", ewr.ID, serviceID).Scan(&count).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	if count == 0 {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee work range %s does not have service %s", ewr.ID, serviceID))
	}
	return nil
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
	employee := &Employee{UserID: ewr.EmployeeID}
	branch := &Branch{BaseModel: BaseModel{ID: ewr.BranchID}}
	if err := employee.HasBranch(tx, ewr.BranchID); err != nil {
		return err
	}
	for _, service := range services {
		if err := ewr.HasService(tx, service.ID); err != nil {
			// Check if the employee has the service
			if err := employee.HasService(tx, service); err != nil {
				return err
			}
			if err := branch.HasService(tx, service); err != nil {
				return err
			}
			if err := tx.Model(ewr).Association("Services").Append(service); err != nil {
				return fmt.Errorf("error adding service %s to work range: %w", service.ID, err)
			}
			ewr.Services = append(ewr.Services, service)
		}
	}
	return nil
}

