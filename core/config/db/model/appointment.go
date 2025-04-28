package model

import (
	"agenda-kaki-go/core/lib" // Adjust import path if necessary
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	//"gorm.io/gorm/clause" // Required for locking in service layer
)

// BaseModel should be defined elsewhere in your model package or a common place

// --- Main Appointment Model ---
// Using your actual Service, Employee, Client, Branch, Company types now
type Appointment struct {
	BaseModel

	ServiceID  uuid.UUID `gorm:"type:uuid;not null;index" json:"service_id"`
	Service    *Service  `gorm:"foreignKey:ServiceID;references:ID;constraint:OnDelete:CASCADE;"` // Using your Service type
	EmployeeID uuid.UUID `gorm:"type:uuid;not null;index" json:"employee_id"`
	Employee   *Employee `gorm:"foreignKey:EmployeeID;references:ID;constraint:OnDelete:CASCADE;"` // Using your Employee type
	ClientID   uuid.UUID `gorm:"type:uuid;not null;index" json:"client_id"`
	Client     *Client   `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE;"` // Using your Client type
	BranchID   uuid.UUID `gorm:"type:uuid;not null;index" json:"branch_id"`
	Branch     *Branch   `gorm:"foreignKey:BranchID;references:ID;constraint:OnDelete:CASCADE;"` // Using your Branch type
	CompanyID  uuid.UUID `gorm:"type:uuid;not null;index" json:"company_id"`
	Company    *Company  `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;" json:"-"` // Company loaded via FK, json:"-" often good practice

	StartTime time.Time `gorm:"not null;index" json:"start_time"`
	EndTime   time.Time `gorm:"not null;index" json:"end_time"`

	Cancelled bool `gorm:"index;default:false" json:"cancelled"`
}

func (Appointment) TableName() string { return "appointments" }

// Indexes - ensure these align with your DB schema strategy
func (Appointment) Indexes() map[string]string {
	return map[string]string{
		"idx_employee_time_active": "CREATE INDEX IF NOT EXISTS idx_employee_time_active ON appointments (employee_id, start_time, end_time, cancelled)",
		"idx_client_time_active":   "CREATE INDEX IF NOT EXISTS idx_client_time_active ON appointments (client_id, start_time, end_time, cancelled)",
		"idx_branch_time_active":   "CREATE INDEX IF NOT EXISTS idx_branch_time_active ON appointments (branch_id, start_time, end_time, cancelled)",
		"idx_company_active":       "CREATE INDEX IF NOT EXISTS idx_company_active ON appointments (company_id, cancelled)",
		"idx_start_time_active":    "CREATE INDEX IF NOT EXISTS idx_start_time_active ON appointments (start_time, cancelled)",
	}
}

// --- GORM Hooks ---

func (a *Appointment) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	// Run full validation rules, now returns lib.ErrorStruct on failure
	err := a.validateAppointmentRules(tx, true) // isCreate = true
	if err != nil {
		// Return the lib.ErrorStruct directly
		return err
	}
	return nil
}

func (a *Appointment) BeforeUpdate(tx *gorm.DB) error {
	runValidation := tx.Statement.Changed("StartTime", "EndTime", "EmployeeID", "ServiceID", "BranchID") ||
		(tx.Statement.Changed("Cancelled") && !a.Cancelled) // Run validation if reactivating

	if runValidation {
		// validateAppointmentRules loads associations if needed.
		err := a.validateAppointmentRules(tx, false) // isCreate = false
		if err != nil {
			// Return the lib.ErrorStruct directly
			return err
		}
	}
	return nil
}

func (a *Appointment) BeforeDelete(tx *gorm.DB) error {
	return lib.Error.General.DeletedError.WithError(fmt.Errorf("deleting appointments is totally forbidden in this system"))
}

func (a *Appointment) AfterCreate(tx *gorm.DB) error {
	err := CreateAppointmentHistory(tx, a, ActionCreate, "Appointment created.")
	if err != nil {
		// Return the wrapped lib.ErrorStruct directly
		return lib.Error.Appointment.HistoryLoggingFailed.WithError(err)
	}
	return nil
}

func (a *Appointment) AfterUpdate(tx *gorm.DB) error {
	action := ActionUpdate
	notes := "Appointment updated."
	logHistory := true

	if tx.Statement.Changed("Cancelled") {
		if a.Cancelled {
			action, notes = ActionCancel, "Appointment cancelled."
		} else {
			notes = "Appointment reactivated."
		}
	} else if !tx.Statement.Changed() {
		logHistory = false
	}

	if logHistory {
		err := CreateAppointmentHistory(tx, a, action, notes)
		if err != nil {
			// Return the wrapped lib.ErrorStruct directly
			return lib.Error.Appointment.HistoryLoggingFailed.WithError(err)
		}
	}
	return nil
}

// --- Validation Helper ---
// This function is called from the hooks and can be reused in other contexts if needed
func (a *Appointment) validateAppointmentRules(tx *gorm.DB, isCreate bool) error {

	// 0. Basic Required Fields & Time Checks
	if a.StartTime.IsZero() {
		return lib.Error.Appointment.StartTimeInThePast
	}
	if a.ServiceID == uuid.Nil || a.EmployeeID == uuid.Nil || a.ClientID == uuid.Nil || a.BranchID == uuid.Nil || a.CompanyID == uuid.Nil {
		return lib.Error.Appointment.MissingRequiredIDs
	}
	// Check start time only on creation or if it explicitly changed to the past
	if isCreate || tx.Statement.Changed("StartTime") {
		if a.StartTime.Before(time.Now().Add(-1 * time.Minute)) {
			return lib.Error.Appointment.StartTimeInThePast // Use specific error for past time
		}
	} else if tx.Statement.Changed("CompanyID") {
		return lib.Error.Appointment.UpdateFailed.WithError(fmt.Errorf("can not change company id")) // Use specific error for company ID change
	}

	// 1. Load Required Associations (Defensively)
	var err error

	err = tx.Model(&Service{}).Preload(clause.Associations).First(&a.Service, a.ServiceID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return lib.Error.Appointment.NotFound.WithError(fmt.Errorf("service ID %s", a.ServiceID))
	}
	if err != nil {
		return lib.Error.Appointment.AssociationLoadFailed.WithError(fmt.Errorf("loading service: %w", err))
	}

	err = tx.Model(&Employee{}).Preload(clause.Associations).First(&a.Employee, a.EmployeeID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return lib.Error.Employee.NotFound.WithError(fmt.Errorf("employee ID %s", a.EmployeeID))
	}
	if err != nil {
		return lib.Error.Appointment.AssociationLoadFailed.WithError(fmt.Errorf("loading employee: %w", err))
	}

	err = tx.Model(&Branch{}).Preload(clause.Associations).First(&a.Branch, a.BranchID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return lib.Error.Branch.NotFound.WithError(fmt.Errorf("branch ID %s", a.BranchID))
	} // Use Branch NotFound
	if err != nil {
		return lib.Error.Appointment.AssociationLoadFailed.WithError(fmt.Errorf("loading branch: %w", err))
	}

	// 2. Calculate & Validate EndTime
	if a.Service.Duration <= 0 { // Use uint duration from your model
		return lib.Error.Appointment.InvalidServiceDuration
	}
	// Convert uint duration (minutes) to time.Duration
	a.EndTime = a.StartTime.Add(time.Duration(a.Service.Duration) * time.Minute)
	if !a.EndTime.After(a.StartTime) {
		return lib.Error.Appointment.EndTimeBeforeStart
	}

	// --- If being cancelled, validation stops here ---
	if a.Cancelled {
		return nil
	}

	// --- Remaining checks only apply to active (non-cancelled) appointments ---

	// 3. Relationship & Existence Checks (Use loaded structs)
	if a.Branch.CompanyID != a.CompanyID {
		return lib.Error.Company.BranchDoesNotBelong
	}
	if a.Employee.CompanyID != a.CompanyID {
		return lib.Error.Company.EmployeeDoesNotBelong
	}
	if a.Service.CompanyID != a.CompanyID {
		return lib.Error.Company.ServiceDoesNotBelong
	}

	// Use simpler direct checks or JOINs if performance demands, raw SQL less type-safe.
	var serviceInBranch bool
	if a.Branch.Services != nil {
		for _, s := range a.Branch.Services {
			if s.ID == a.ServiceID {
				serviceInBranch = true
				break
			}
		}
	} else {
		var count int64
		tx.Table("branch_services").Where("branch_id = ? AND service_id = ?", a.BranchID, a.ServiceID).Count(&count)
		serviceInBranch = count > 0
	}
	if !serviceInBranch {
		return lib.Error.Branch.ServiceDoesNotBelong
	}

	// Check if employee offers the service
	var employeeService bool
	if a.Employee.Services != nil {
		for _, s := range a.Employee.Services {
			if s.ID == a.ServiceID {
				employeeService = true
				break
			}
		}
	} else {
		var count int64
		tx.Table("employee_services").Where("employee_id = ? AND service_id = ?", a.EmployeeID, a.ServiceID).Count(&count)
		employeeService = count > 0
	}
	if !employeeService {
		return lib.Error.Employee.LacksService
	} // Use specific error

	// Check if employee works at the branch
	var employeeInBranch bool
	if a.Employee.Branches != nil { // Check preloaded (unlikely)
		for _, b := range a.Employee.Branches {
			if b.ID == a.BranchID {
				employeeInBranch = true
				break
			}
		}
	} else {
		var count int64
		tx.Table("employee_branches").Where("employee_id = ? AND branch_id = ?", a.EmployeeID, a.BranchID).Count(&count)
		employeeInBranch = count > 0
	}
	if !employeeInBranch {
		return lib.Error.Employee.BranchDoesNotBelong
	}

	// 4. Check Employee Availability (Work Schedule)
	if a.Employee.WorkSchedule.IsEmpty() {
		return lib.Error.Employee.NoWorkScheduleForDay
	}
	weekday := a.StartTime.Weekday()
	workRanges := a.Employee.WorkSchedule.GetRangesForDay(weekday) // Use helper if available

	if len(workRanges) == 0 {
		return lib.Error.Employee.NoWorkScheduleForDay
	}

	isAvailableInSchedule := false
	appointmentDateStr := a.StartTime.Format("2006-01-02")
	for _, wr := range workRanges {
		if wr.BranchID != a.BranchID {
			continue
		}

		// Use layout matching your "15:30:00" format
		layout := "2006-01-02 15:04"
		scheduleStart, errStart := time.ParseInLocation(layout, fmt.Sprintf("%s %s", appointmentDateStr, wr.Start), a.StartTime.Location())
		scheduleEnd, errEnd := time.ParseInLocation(layout, fmt.Sprintf("%s %s", appointmentDateStr, wr.End), a.StartTime.Location())

		if errStart != nil || errEnd != nil {
			return lib.Error.Appointment.InvalidWorkScheduleFormat.WithError(fmt.Errorf("invalid work schedule format parse for employee %s with range %s-%s. start time error: %v / end time error: %v", a.EmployeeID, wr.Start, wr.End, errStart, errEnd))
		}
		if !scheduleEnd.After(scheduleStart) {
			continue
		}

		fitsStart := !a.StartTime.Before(scheduleStart)
		fitsEnd := a.EndTime.Before(scheduleEnd) || a.EndTime.Equal(scheduleEnd)

		if fitsStart && fitsEnd {
			isAvailableInSchedule = true
			break
		}
	}

	if !isAvailableInSchedule {
		return lib.Error.Employee.NotAvailableWorkSchedule // Use specific error
	}

	// 5. Overlap Checks (for active, non-self appointments)
	overlapCondition := `start_time < ? AND end_time > ?`
	baseOverlapQuery := tx.Model(&Appointment{}).
		Where("cancelled = ?", false).
		Where("id != ?", a.ID).
		Where(overlapCondition, a.EndTime, a.StartTime)

	var count int64

	// Check Employee Overlap
	err = baseOverlapQuery.Where("employee_id = ?", a.EmployeeID).Count(&count).Error
	if err != nil {
		return lib.Error.Appointment.AssociationLoadFailed.WithError(fmt.Errorf("db error checking employee overlap: %w", err))
	} // Generic load/db error
	if count > 0 {
		return lib.Error.Employee.ScheduleConflict
	}

	// Check Client Overlap
	count = 0
	err = baseOverlapQuery.Where("client_id = ?", a.ClientID).Count(&count).Error
	if err != nil {
		return lib.Error.Appointment.AssociationLoadFailed.WithError(fmt.Errorf("db error checking client overlap: %w", err))
	}
	if count > 0 {
		return lib.Error.Client.ScheduleConflict
	}

	// Check Branch Service Capacity (Using JSONB field on Branch)
	count = 0
	err = baseOverlapQuery.Where("branch_id = ? AND service_id = ?", a.BranchID, a.ServiceID).Count(&count).Error
	if err != nil {
		return lib.Error.Appointment.AssociationLoadFailed.WithError(fmt.Errorf("db error checking service capacity: %w", err))
	}

	serviceMax := uint(0) // Use uint to match your model
	// Branch should be loaded, iterate over its ServiceDensity slice
	for _, sd := range a.Branch.ServiceDensity { // Access JSONB field directly
		if sd.ServiceID == a.ServiceID {
			serviceMax = sd.MaxSchedulesOverlap // Use uint MaxSchedulesOverlap
			break
		}
	}
	// Cast serviceMax to int64 for comparison with count
	if serviceMax > 0 && count >= int64(serviceMax) {
		return lib.Error.Branch.MaxServiceCapacityReached // Use new specific error
	}

	// Check Overall Branch Capacity (Using BranchDensity field on Branch)
	count = 0
	err = baseOverlapQuery.Where("branch_id = ?", a.BranchID).Count(&count).Error
	if err != nil {
		return lib.Error.Appointment.AssociationLoadFailed.WithError(fmt.Errorf("db error checking branch capacity: %w", err))
	}

	branchMax := a.Branch.BranchDensity
	if branchMax > 0 && count >= int64(branchMax) {
		return lib.Error.Branch.MaxCapacityReached // Use new specific error
	}

	return nil // All validations passed
}
