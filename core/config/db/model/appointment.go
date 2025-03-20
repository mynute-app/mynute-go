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
	// Ensure start time is in the future
	if a.StartTime.Before(time.Now()) {
		return errors.New("start time must be in the future")
	}

	// Calculate EndTime based on Service duration
	a.EndTime = a.StartTime.Add(time.Duration(a.Service.Duration) * time.Minute)
	if a.EndTime.Before(a.StartTime) {
		return errors.New("end time must be after start time")
	}

	// TODO: Load all associations into the appointment struct
	if err := tx.Model(&Service{}).Find(&a.Service, a.ServiceID).Error; err != nil {
		return err
	}
	if err := tx.Model(&Employee{}).Find(&a.Employee, a.EmployeeID).Error; err != nil {
		return err
	}
	if err := tx.Model(&Branch{}).Find(&a.Branch, a.BranchID).Error; err != nil {
		return err
	}
	if err := tx.Model(&Company{}).Find(&a.Company, a.CompanyID).Error; err != nil {
		return err
	}
	if err := tx.Model(&User{}).Find(&a.User, a.UserID).Error; err != nil {
		return err
	}

	// TODO: Check if branch belongs to company
	if a.Branch.CompanyID != a.CompanyID {
		return errors.New("branch does not belong to the company")
	}
	// TODO: Check if employee belongs to company
	if a.Employee.CompanyID != a.CompanyID {
		return errors.New("employee does not belong to the company")
	}
	// TODO: Check if the service belongs to the company
	if a.Service.CompanyID != a.CompanyID {
		return errors.New("service does not belong to the company")
	}
	// TODO: Check if service exists in the branch
	var serviceExistsInBranch int64
	if err := tx.Model(&Branch{}).
		Where("id = ? AND EXISTS (SELECT 1 FROM branch_services WHERE branch_id = ? AND service_id = ?)", a.BranchID, a.BranchID, a.ServiceID).
		Count(&serviceExistsInBranch).Error; err != nil {
		return err
	}
	if serviceExistsInBranch == 0 {
		return errors.New("service does not exist in the branch")
	}
	// TODO: Check if employee has the service
	var employeeHasService int64
	if err := tx.Model(&Employee{}).
		Where("id = ? AND EXISTS (SELECT 1 FROM employee_services WHERE employee_id = ? AND service_id = ?)", a.EmployeeID, a.EmployeeID, a.ServiceID).
		Count(&employeeHasService).Error; err != nil {
		return err
	}
	if employeeHasService == 0 {
		return errors.New("employee does not have the service")
	}
	// TODO: Check if employee exists in the branch
	var employeeExistsInBranch int64
	if err := tx.Model(&Employee{}).
		Where("id = ? AND EXISTS (SELECT 1 FROM employee_branches WHERE employee_id = ? AND branch_id = ?)", a.EmployeeID, a.EmployeeID, a.BranchID).
		Count(&employeeExistsInBranch).Error; err != nil {
		return err
	}
	if employeeExistsInBranch == 0 {
		return errors.New("employee does not exist in the branch")
	}

	// Check if the employee is available in the branch at this start time
	weekday := a.StartTime.Weekday()
	var WorkRanges []WorkRange

	switch weekday {
	case time.Monday:
		WorkRanges = a.Employee.WorkSchedule.Monday
	case time.Tuesday:
		WorkRanges = a.Employee.WorkSchedule.Tuesday
	case time.Wednesday:
		WorkRanges = a.Employee.WorkSchedule.Wednesday
	case time.Thursday:
		WorkRanges = a.Employee.WorkSchedule.Thursday
	case time.Friday:
		WorkRanges = a.Employee.WorkSchedule.Friday
	case time.Saturday:
		WorkRanges = a.Employee.WorkSchedule.Saturday
	case time.Sunday:
		WorkRanges = a.Employee.WorkSchedule.Sunday
	}

	if WorkRanges == nil {
		return errors.New("employee has no work schedule for the selected day")
	}

	available := false
	for _, wr := range WorkRanges {
		StartTimeDate := fmt.Sprintf("%d-%d-%d", a.StartTime.Year(), a.StartTime.Month(), a.StartTime.Day())
		wrStart, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", StartTimeDate, wr.Start))
		if err != nil {
			return err
		}
		wrEnd, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", StartTimeDate, wr.End))
		if err != nil {
			return err
		}
		if wr.BranchID == a.BranchID && a.StartTime.After(wrStart) && a.StartTime.Before(wrEnd) {
			available = true
			break
		}
	}

	if !available {
		return errors.New("employee is not available at schedule time in the specified branch")
	}

	// Checks for Employee Overlapping
	var EmployeeOverlappingCount int64
	if err := tx.Model(&Appointment{}).
		Where(`employee_id = ? AND (
			(start_time < ? AND end_time > ?) 
			OR (start_time >= ? AND start_time < ?) 
			OR (end_time > ? AND end_time <= ?)
		)`, a.EmployeeID, a.EndTime, a.StartTime, a.StartTime, a.EndTime, a.StartTime, a.EndTime).
		Count(&EmployeeOverlappingCount).Error; err != nil {
		return err
	}
	if EmployeeOverlappingCount > 0 {
		return errors.New("appointment conflicts with another employee booking")
	}

	// Checks for User Overlapping
	var UserOverlappingCount int64
	if err := tx.Model(&Appointment{}).
		Where(`user_id = ? AND (
			(start_time < ? AND end_time > ?) 
			OR (start_time >= ? AND start_time < ?) 
			OR (end_time > ? AND end_time <= ?)
		)`, a.UserID, a.EndTime, a.StartTime, a.StartTime, a.EndTime, a.StartTime, a.EndTime).
		Count(&UserOverlappingCount).Error; err != nil {
		return err
	}

	if UserOverlappingCount > 0 {
		return errors.New("appointment conflicts with another user booking")
	}

	// Check for overlapping schedules in the service at the branch
	var ServiceScheduleOverlapping int64
	if err := tx.Model(&Appointment{}).
		Where(`branch_id = ? AND service_id = ? AND (
			(start_time < ? AND end_time > ?) 
			OR (start_time >= ? AND start_time < ?) 
			OR (end_time > ? AND end_time <= ?)
		)`, a.BranchID, a.ServiceID, a.EndTime, a.StartTime, a.StartTime, a.EndTime, a.StartTime, a.EndTime).
		Count(&ServiceScheduleOverlapping).Error; err != nil {
		return err
	}

	var ServiceMaxSchedulesOverlap int64

	for _, sd := range a.Branch.ServiceDensity {
		if sd.ServiceID == a.ServiceID {
			ServiceMaxSchedulesOverlap = int64(sd.MaxSchedulesOverlap)
			break
		}
	}

	if ServiceMaxSchedulesOverlap > 0 && ServiceScheduleOverlapping >= ServiceMaxSchedulesOverlap {
		return fmt.Errorf("max limit of %d concurrent appointments reached for this service at branch", ServiceMaxSchedulesOverlap)
	}

	// Check for overlapping schedules in the branch
	var BranchOverlappingSchedules int64
	if err := tx.Model(&Appointment{}).
		Where(`branch_id = ? AND (
			(start_time < ? AND end_time > ?)
			OR (start_time >= ? AND start_time < ?)
			OR (end_time > ? AND end_time <= ?)
		)`, a.BranchID, a.EndTime, a.StartTime, a.StartTime, a.EndTime, a.StartTime, a.EndTime).
		Count(&BranchOverlappingSchedules).Error; err != nil {
		return err
	}

	BranchMaxSchedulesOverlap := int64(a.Branch.BranchDensity)

	if BranchOverlappingSchedules >= BranchMaxSchedulesOverlap {
		return fmt.Errorf("max limit of %d concurrent appointments reached for this branch", BranchMaxSchedulesOverlap)
	}

	return nil
}
