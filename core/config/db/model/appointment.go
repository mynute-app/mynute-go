package model

import (
	"agenda-kaki-go/core/lib"
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

	ClientID uint    `gorm:"not null;index" json:"client_id"`
	Client   *Client `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE;"`

	BranchID uint    `gorm:"not null;index" json:"branch_id"`
	Branch   *Branch `gorm:"foreignKey:BranchID;references:ID;constraint:OnDelete:CASCADE;"`

	CompanyID uint     `gorm:"not null;index" json:"company_id"`
	Company   *Company `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;" json:"company"`

	StartTime time.Time `gorm:"not null" json:"start_time"`
	EndTime   time.Time `gorm:"not null" json:"end_time"`

	Rescheduled bool `gorm:"default:false" json:"rescheduled"`
	Cancelled   bool `gorm:"default:false" json:"cancelled"`

	RescheduledToID   *uint `json:"rescheduled_to_id"`
	RescheduledFromID *uint `json:"rescheduled_from_id"`
}

// Custom Composite Index
func (Appointment) TableName() string {
	return "appointments"
}

func (Appointment) Indexes() map[string]string {
	return map[string]string{
		"idx_employee_time": "CREATE INDEX idx_employee_time ON appointments (employee_id, start_time, end_time)",
		"idx_client_time":   "CREATE INDEX idx_client_time ON appointments (client_id, start_time, end_time)",
		"idx_branch_time":   "CREATE INDEX idx_branch_time ON appointments (branch_id, start_time, end_time)",
	}
}

func (a *Appointment) BeforeCreate(tx *gorm.DB) error {
	// Ensure start time is in the future
	if a.StartTime.Before(time.Now()) {
		return lib.Error.Appointment.StartTimeInThePast
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
	if err := tx.Model(&Client{}).Find(&a.Client, a.ClientID).Error; err != nil {
		return err
	}

	// Calculate EndTime based on Service duration
	a.EndTime = a.StartTime.Add(time.Duration(a.Service.Duration) * time.Minute)
	if a.EndTime.Before(a.StartTime) {
		return lib.Error.Appointment.EndTimeBeforeStart
	}

	// TODO: Check if branch belongs to company
	if a.Branch.CompanyID != a.CompanyID {
		return lib.Error.Company.BranchDoesNotBelong
	}
	// TODO: Check if employee belongs to company
	if a.Employee.CompanyID != a.CompanyID {
		return lib.Error.Company.EmployeeDoesNotBelong
	}
	// TODO: Check if the service belongs to the company
	if a.Service.CompanyID != a.CompanyID {
		return lib.Error.Company.ServiceDoesNotBelong
	}
	// TODO: Check if service exists in the branch
	var serviceExistsInBranch int64
	if err := tx.Model(&Branch{}).
		Where("id = ? AND EXISTS (SELECT 1 FROM branch_services WHERE branch_id = ? AND service_id = ?)", a.BranchID, a.BranchID, a.ServiceID).
		Count(&serviceExistsInBranch).Error; err != nil {
		return err
	}
	if serviceExistsInBranch == 0 {
		return lib.Error.Branch.ServiceDoesNotBelong
	}
	// TODO: Check if employee has the service
	var employeeHasService int64
	if err := tx.Model(&Employee{}).
		Where("id = ? AND EXISTS (SELECT 1 FROM employee_services WHERE employee_id = ? AND service_id = ?)", a.EmployeeID, a.EmployeeID, a.ServiceID).
		Count(&employeeHasService).Error; err != nil {
		return err
	}
	if employeeHasService == 0 {
		return lib.Error.Employee.ServiceDoesNotBelong
	}
	// TODO: Check if employee exists in the branch
	var employeeExistsInBranch int64
	if err := tx.Model(&Employee{}).
		Where("id = ? AND EXISTS (SELECT 1 FROM employee_branches WHERE employee_id = ? AND branch_id = ?)", a.EmployeeID, a.EmployeeID, a.BranchID).
		Count(&employeeExistsInBranch).Error; err != nil {
		return err
	}
	if employeeExistsInBranch == 0 {
		return lib.Error.Employee.BranchDoesNotBelong
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

	if len(WorkRanges) == 0 {
		return lib.Error.Employee.NoWorkScheduleForDay
	}

	available := false
	for _, wr := range WorkRanges {
		StartTimeDate := fmt.Sprintf("%d-%02d-%02d", a.StartTime.Year(), a.StartTime.Month(), a.StartTime.Day()) // Ensure zero-padding
		wrStart, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", StartTimeDate, wr.Start))
		if err != nil {
			return err
		}
		wrEnd, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", StartTimeDate, wr.End))
		if err != nil {
			return err
		}
		isAfterOrEqual := a.StartTime.After(wrStart) || a.StartTime.Equal(wrStart)
		isBeforeOrEqual := a.StartTime.Before(wrEnd) || a.StartTime.Equal(wrEnd)
		if wr.BranchID == a.BranchID && isAfterOrEqual && isBeforeOrEqual {
			available = true
			break
		}
	}

	if !available {
		return lib.Error.Employee.NotAvailableOnDate
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
		return lib.Error.Employee.ScheduleConflict
	}

	// Checks for Client Overlapping
	var ClientOverlappingCount int64
	if err := tx.Model(&Appointment{}).
		Where(`client_id = ? AND (
			(start_time < ? AND end_time > ?) 
			OR (start_time >= ? AND start_time < ?) 
			OR (end_time > ? AND end_time <= ?)
		)`, a.ClientID, a.EndTime, a.StartTime, a.StartTime, a.EndTime, a.StartTime, a.EndTime).
		Count(&ClientOverlappingCount).Error; err != nil {
		return err
	}

	if ClientOverlappingCount > 0 {
		return lib.Error.Client.ScheduleConflict
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
		return lib.Error.Branch.MaxConcurrentAppointmentsForService
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
		return lib.Error.Branch.MaxConcurrentAppointmentsGeneral
	}

	return nil
}
