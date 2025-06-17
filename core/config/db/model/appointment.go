package model

import (
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/lib"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AppointmentBase struct {
	ServiceID             uuid.UUID  `gorm:"type:uuid;not null" json:"service_id"`
	EmployeeID            uuid.UUID  `gorm:"type:uuid;not null" json:"employee_id"`
	ClientID              uuid.UUID  `gorm:"type:uuid;not null;index" json:"client_id"`
	BranchID              uuid.UUID  `gorm:"type:uuid;not null" json:"branch_id"`
	PaymentID             *uuid.UUID `gorm:"type:uuid;uniqueIndex" json:"payment_id"`
	CompanyID             uuid.UUID  `gorm:"type:uuid;not null;index" json:"company_id"`
	CancelledEmployeeID   *uuid.UUID `gorm:"type:uuid" json:"cancelled_employee_id"`
	StartTime             time.Time  `gorm:"not null" json:"start_time"`
	EndTime               time.Time  `gorm:"not null" json:"end_time"`
	ActualStartTime       time.Time  `json:"actual_start_time"`
	ActualEndTime         time.Time  `json:"actual_end_time"`
	CancelTime            time.Time  `json:"cancel_time"`
	IsFulfilled           bool       `gorm:"default:false" json:"is_fulfilled"`
	IsCancelled           bool       `gorm:"default:false" json:"is_cancelled"`
	IsCancelledByClient   bool       `gorm:"default:false" json:"is_cancelled_by_client"`
	IsCancelledByEmployee bool       `gorm:"default:false" json:"is_cancelled_by_employee"`
	IsConfirmedByClient   bool       `gorm:"default:false" json:"is_confirmed_by_client"`
}

// This is the foreign key struct for the Appointment model at company schema level.
type AppointmentFK struct {
	Service  *Service  `gorm:"foreignKey:ServiceID;references:ID;constraint:OnDelete:CASCADE;"`  // Using your Service type
	Employee *Employee `gorm:"foreignKey:EmployeeID;references:ID;constraint:OnDelete:CASCADE;"` // Using your Employee type
	Branch   *Branch   `gorm:"foreignKey:BranchID;references:ID;constraint:OnDelete:CASCADE;"`   // Using your Branch type
	Payment  *Payment  `gorm:"foreignKey:PaymentID;references:ID;constraint:OnDelete:CASCADE;"`  // Using your Payment type
}

type AppointmentJson struct {
	History  mJSON.AppointmentHistory `gorm:"type:jsonb" json:"history"`  // JSONB field for history changes
	Comments mJSON.Comments           `gorm:"type:jsonb" json:"comments"` // JSONB field for comments
}

// --- Main Appointment Model ---
// Using your actual Service, Employee, Client, Branch, Company types now
type Appointment struct {
	BaseModel
	AppointmentBase
	AppointmentFK
	AppointmentJson
}

var AppointmentTableName = "appointments"

func (Appointment) TableName() string { return AppointmentTableName }

func (Appointment) SchemaType() string { return "company" }

func (Appointment) Indexes() map[string]string {
	return AppointmentIndexes(AppointmentTableName)
}

func AppointmentIndexes(table string) map[string]string {
	return map[string]string{
		"idx_employee_time_active": fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_employee_time_active ON %s (employee_id, start_time, end_time, is_cancelled)", table),
		"idx_client_time_active":   fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_client_time_active ON %s (client_id, start_time, end_time, is_cancelled)", table),
		"idx_branch_time_active":   fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_branch_time_active ON %s (branch_id, start_time, end_time, is_cancelled)", table),
		"idx_company_time_active":  fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_company_time_active ON %s (company_id, start_time, end_time, is_cancelled)", table),
		"idx_start_time_active":    fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_start_time_active ON %s (start_time, is_cancelled)", table),
	}
}

// --- Appointment Hooks ---

func (a *Appointment) AfterCreate(tx *gorm.DB) error {
	var client Client
	if err := tx.Model(&Client{}).Where("id = ?", a.ClientID).First(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lib.Error.Client.NotFound.WithError(fmt.Errorf("client ID %s", a.ClientID))
		}
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("loading client: %w", err))
	}
	var company Company
	if err := tx.Model(&company).Where("id = ?", a.CompanyID).First(&company).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lib.Error.Company.NotFound.WithError(fmt.Errorf("company ID %s", a.CompanyID))
		}
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("loading company: %w", err))
	}
	client.AddAppointment(a, a.Service, &company, a.Branch, a.Employee)
	if err := tx.Save(&client).Error; err != nil {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("updating client: %w", err))
	}
	return nil
}

func (a *Appointment) BeforeCreate(tx *gorm.DB) error {
	if !a.History.IsEmpty() {
		return lib.Error.Appointment.HistoryManualUpdateForbidden
	}
	err := a.ValidateRules(tx, true) // isCreate = true
	if err != nil {
		// Return the lib.ErrorStruct directly
		return err
	}
	return nil
}

func (a *Appointment) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("CompanyID") {
		return lib.Error.General.UpdatedError.WithError(errors.New("the CompanyID cannot be changed after creation"))
	}

	if tx.Statement.Changed("History") {
		return lib.Error.Appointment.HistoryManualUpdateForbidden
	}

	var changes []mJSON.FieldChange

	if a.History.IsEmpty() {
		a.History = mJSON.AppointmentHistory{FieldChanges: []mJSON.FieldChange{}}
	}

	if tx.Statement.Schema == nil {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("statement schema is nil at database transaction layer"))
	}

	var originalAppointment Appointment
	if err := tx.Model(&Appointment{}).
		Where("id = ?", a.ID).
		First(&originalAppointment).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lib.Error.Appointment.NotFound.WithError(fmt.Errorf("appointment ID %s", a.ID))
		}
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("loading appointment: %w", err))
	}

	for _, field := range tx.Statement.Schema.Fields {
		if (!tx.Statement.Changed(field.Name) && field.Name != "history") || field.Name == "id" {
			continue
		}
		oldValue, _ := field.ValueOf(tx.Statement.Context, reflect.ValueOf(originalAppointment))
		newValue, _ := field.ValueOf(tx.Statement.Context, reflect.ValueOf(a))
		if oldValue != newValue {
			changes = append(changes, mJSON.FieldChange{
				CreatedAt: time.Now(),
				Field:     field.Name,
				OldValue:  fmt.Sprintf("%v", oldValue),
				NewValue:  fmt.Sprintf("%v", newValue),
			})
		}
	}

	if len(changes) > 0 {
		err := a.ValidateRules(tx, false)
		if err != nil {
			return err
		}
		a.History.FieldChanges = append(originalAppointment.History.FieldChanges, changes...)
	}

	return nil
}

func (a *Appointment) BeforeDelete(tx *gorm.DB) error {
	return lib.Error.General.DeletedError.WithError(fmt.Errorf("deleting appointments is totally forbidden in this system"))
}

// --- Validation Helper ---
// This function is called from the hooks and can be reused in other contexts if needed
func (a *Appointment) ValidateRules(tx *gorm.DB, isCreate bool) error {

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
	if a.IsCancelled {
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
		Where("is_cancelled = ?", false).
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

func (a *Appointment) Refresh(tx *gorm.DB) error {
	if err := tx.Model(&Appointment{}).Where("id = ?", a.ID).Preload(clause.Associations).First(a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lib.Error.Appointment.NotFound.WithError(fmt.Errorf("appointment ID %s", a.ID))
		}
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("loading appointment: %w", err))
	}
	return nil
}

func (a *Appointment) Cancel(tx *gorm.DB) error {
	if err := tx.Model(&Appointment{}).Where("id = ?", a.ID).First(a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lib.Error.Appointment.NotFound.WithError(fmt.Errorf("appointment ID %s", a.ID))
		}
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("loading appointment: %w", err))
	}
	if a.IsCancelled {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("appointment is already cancelled"))
	} else if a.IsFulfilled {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("appointment is already fulfilled"))
	} else if time.Now().After(a.StartTime) {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("cannot cancel appointment as it already happened"))
	}
	a.IsCancelled = true
	a.CancelTime = time.Now()
	err := tx.Save(a).Error
	if err != nil {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("error cancelling appointment: %w", err))
	}
	return nil
}
