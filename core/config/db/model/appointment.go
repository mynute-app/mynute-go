package model

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Fifth step: Scheduling with service and employee availability.
type Appointment struct {
	gorm.Model
	ServiceID uint     `gorm:"not null;index" json:"service_id"`
	Service   *Service `gorm:"foreignKey:ServiceID;references:ID;constraint:OnDelete:CASCADE;"`

	EmployeeID uint      `gorm:"not null;index" json:"employee_id"`
	Employee   *Employee `gorm:"foreignKey:EmployeeID;references:ID;constraint:OnDelete:CASCADE;"`

	UserID uint  `gorm:"not null;index" json:"user_id"`
	User   *User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`

	BranchID uint    `gorm:"not null;index" json:"branch_id"`
	Branch   *Branch `gorm:"foreignKey:BranchID;references:ID;constraint:OnDelete:CASCADE;"`

	CompanyID uint     `gorm:"not null;index" json:"company_id"`
	Company   *Company `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;" json:"company"`

	StartTime time.Time `gorm:"not null" json:"start_time"`
	EndTime   time.Time `gorm:"not null" json:"end_time"`
}

// Custom Composite Index
func (Appointment) TableName() string {
	return "appointments"
}

func (Appointment) Indexes() map[string]string {
	return map[string]string{
		"idx_employee_time": "CREATE INDEX idx_employee_time ON appointments (employee_id, start_time, end_time)",
		"idx_user_time":     "CREATE INDEX idx_user_time ON appointments (user_id, start_time, end_time)",
		"idx_branch_time":   "CREATE INDEX idx_branch_time ON appointments (branch_id, start_time, end_time)",
	}
}

func (a *Appointment) BeforeCreate(tx *gorm.DB) error {
	if a.StartTime.Before(time.Now()) {
		return errors.New("start time must be in the future")
	}
	if err := tx.First(&a.Service, a.ServiceID).Error; err != nil {
		return err
	}
	a.EndTime = a.StartTime.Add(time.Duration(a.Service.Duration) * time.Minute)
	if a.EndTime.Before(a.StartTime) {
		return errors.New("end time must be after start time")
	}
	if err := tx.First(&a.Employee, a.EmployeeID).Error; err != nil {
		return err
	}
	
	var OverlappingEmployeeAppointmentExists bool
	if err := tx.Model(&Appointment{}).
		Select("count(*) > 0").
		Where(`
					employee_id = ? 
					AND (
							(start_time < ? AND end_time > ?) 
							OR (start_time >= ? AND start_time < ?) 
							OR (end_time > ? AND end_time <= ?)
					)`,
			a.EmployeeID, a.EndTime, a.StartTime, a.StartTime, a.EndTime, a.StartTime, a.EndTime,
		).Find(&OverlappingEmployeeAppointmentExists).Error; err != nil {
		return err
	}

	if OverlappingEmployeeAppointmentExists {
		return errors.New("could not create this appointment as it overlaps with another appointment already created for the employee")
	}

	var OverlappingUserAppointmentExists bool
	if err := tx.Model(&Appointment{}).
		Select("count(*) > 0").
		Where(`
					user_id = ?
					AND (
							(start_time < ? AND end_time > ?)
							OR (start_time >= ? AND start_time < ?)
							OR (end_time > ? AND end_time <= ?)
					)`,
			a.UserID, a.EndTime, a.StartTime, a.StartTime, a.EndTime, a.StartTime, a.EndTime,
		).Find(&OverlappingUserAppointmentExists).Error; err != nil {
		return err
	}

	if OverlappingUserAppointmentExists {
		return errors.New("could not create this appointment as it overlaps with another appointment already created for the user")
	}

	var OverlappingBranchAppointmentCount int64
	var MaxSchedulesAtSameTime int64
	
	// Query to count overlapping appointments
	if err := tx.Model(&Appointment{}).
			Select("count(*)").
			Where(`
					branch_id = ?
					AND service_id = ?
					AND (
							(start_time < ? AND end_time > ?)
							OR (start_time >= ? AND start_time < ?)
							OR (end_time > ? AND end_time <= ?)
					)`,
					a.BranchID, a.ServiceID, a.EndTime, a.StartTime, a.StartTime, a.EndTime, a.StartTime, a.EndTime,
			).Count(&OverlappingBranchAppointmentCount).Error; err != nil {
			return err
	}
	
	// Fetch MaxSchedulesAtSameTime directly from JSONB in PostgreSQL
	query := `
			SELECT COALESCE(
					(SELECT (jsonb_array_elements(service_density)->>'max_schedules_at_same_time')::int 
					 FROM branches 
					 WHERE id = ? 
					 AND jsonb_array_elements(service_density)->>'service_id' = ?
					 LIMIT 1),
					0
			)`
	
	if err := tx.Raw(query, a.BranchID, fmt.Sprintf("%d", a.ServiceID)).Scan(&MaxSchedulesAtSameTime).Error; err != nil {
			return err
	}
	
	// Check if the overlapping count exceeds the limit
	BranchServiceHasReachedMaxOverlappingAllowed := OverlappingBranchAppointmentCount >= MaxSchedulesAtSameTime
	if BranchServiceHasReachedMaxOverlappingAllowed {
		err_text := fmt.Sprintf("could not create this appointment at selected branch as the maximum overlapping limite of %d appointments for the same service has been reached", MaxSchedulesAtSameTime)
		return errors.New(err_text)
	}

	var MaxBranchDensity uint
	if err := tx.Model(&a.Branch).Select("branch_density").Where("id = ?", a.BranchID).Scan(&MaxBranchDensity).Error; err != nil {
		return err
	}

	BranchHasReachedMaxOverlappingAllowed := OverlappingBranchAppointmentCount >= int64(MaxBranchDensity)
	if BranchHasReachedMaxOverlappingAllowed {
		err_text := fmt.Sprintf("could not create this appointment at selected branch as the maximum overlapping limite of %d appointments for the same branch has been reached", MaxBranchDensity)
		return errors.New(err_text)
	}

	return nil
}
