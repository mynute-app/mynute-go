package model

import (
	"agenda-kaki-go/core/lib" // Adjust import path if necessary
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AppointmentBase struct {
	ServiceID         uuid.UUID           `gorm:"type:uuid;not null;index" json:"service_id"`
	Service           *Service            `gorm:"foreignKey:ServiceID;references:ID;constraint:OnDelete:CASCADE;"` // Using your Service type
	EmployeeID        uuid.UUID           `gorm:"type:uuid;not null;index" json:"employee_id"`
	Employee          *Employee           `gorm:"foreignKey:EmployeeID;references:ID;constraint:OnDelete:CASCADE;"` // Using your Employee type
	ClientID          uuid.UUID           `gorm:"type:uuid;not null;index" json:"client_id"`
	Client            *Client             `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE;"` // Using your Client type
	BranchID          uuid.UUID           `gorm:"type:uuid;not null;index" json:"branch_id"`
	Branch            *Branch             `gorm:"foreignKey:BranchID;references:ID;constraint:OnDelete:CASCADE;"` // Using your Branch type
	PaymentID         uuid.UUID           `gorm:"type:uuid;index" json:"payment_id"`
	Payment           *Payment            `gorm:"foreignKey:PaymentID;references:ID;constraint:OnDelete:CASCADE;"` // Using your Payment type
	CompanyID         uuid.UUID           `gorm:"type:uuid;not null;index" json:"company_id"`
	StartTime         time.Time           `gorm:"not null;index" json:"start_time"`
	EndTime           time.Time           `gorm:"not null;index" json:"end_time"`
	Cancelled         bool                `gorm:"index;default:false" json:"cancelled"`
	ConfirmedByClient bool                `gorm:"index;default:false" json:"confirmed_by_client"`
	Fulfilled         bool                `gorm:"index;default:false" json:"fulfilled"`
	History           AppointmentHistory  `gorm:"type:jsonb" json:"history"`  // JSONB field for history changes
	Comments          AppointmentComments `gorm:"type:jsonb" json:"comments"` // JSONB field for comments
}

// --- Main Appointment Model ---
// Using your actual Service, Employee, Client, Branch, Company types now
type Appointment struct {
	BaseModel
	AppointmentBase
}

type AppointmentComments []Comment

type Comment struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Comment   string         `json:"comment"`
	CreatedBy uuid.UUID      `gorm:"type:uuid;not null;index" json:"created_by"`
	Type      string         `json:"type"` // "internal" or "external"
}

// --- Implement Scanner/Valuer for AppointmentComments ---
func (ac *AppointmentComments) Value() (driver.Value, error) {
	if ac == nil || len(*ac) == 0 {
		// Return empty JSON array `[]` which is valid JSON
		return json.Marshal([]Comment{})
	}
	return json.Marshal(ac)
}

func (ac *AppointmentComments) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		// Handle nil from DB
		if value == nil {
			*ac = []Comment{} // Initialize to empty slice
			return nil
		}
		return errors.New("failed to scan AppointmentComments: expected []byte")
	}
	// Handle empty JSON array or null from DB
	if len(bytes) == 0 || string(bytes) == "null" {
		*ac = []Comment{} // Initialize to empty slice
		return nil
	}
	// Important: Unmarshal into the pointer *ac
	return json.Unmarshal(bytes, ac)
}

// Optional: Add helper methods directly to the type
func (ac *AppointmentComments) Add(c Comment) {
	if ac != nil {
		*ac = append(*ac, c)
	}
}

type AppointmentHistory struct {
	FieldChanges []FieldChange `json:"field_changes"`
}

type FieldChange struct {
	CreatedAt time.Time `json:"created_at"`
	Field     string    `json:"field"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
}

func (ah *AppointmentHistory) Value() (driver.Value, error) {
	return json.Marshal(ah)
}

func (ah *AppointmentHistory) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan WorkSchedule: expected []byte")
	}

	return json.Unmarshal(bytes, ah)
}

func (ah *AppointmentHistory) IsEmpty() bool {
	return ah == nil || len(ah.FieldChanges) == 0
}

func (ah *AppointmentHistory) FilterByField(field string) []FieldChange {
	if ah == nil {
		return nil
	}
	var filteredChanges []FieldChange
	for _, change := range ah.FieldChanges {
		if change.Field == field {
			filteredChanges = append(filteredChanges, change)
		}
	}
	return filteredChanges
}

func (Appointment) TableName() string { return "appointments" }

func (Appointment) Indexes() map[string]string {
	return map[string]string{
		"idx_employee_time_active": "CREATE INDEX IF NOT EXISTS idx_employee_time_active ON appointments (employee_id, start_time, end_time, cancelled)",
		"idx_client_time_active":   "CREATE INDEX IF NOT EXISTS idx_client_time_active ON appointments (client_id, start_time, end_time, cancelled)",
		"idx_branch_time_active":   "CREATE INDEX IF NOT EXISTS idx_branch_time_active ON appointments (branch_id, start_time, end_time, cancelled)",
		"idx_company_active":       "CREATE INDEX IF NOT EXISTS idx_company_active ON appointments (company_id, cancelled)",
		"idx_start_time_active":    "CREATE INDEX IF NOT EXISTS idx_start_time_active ON appointments (start_time, cancelled)",
	}
}

// --- Appointment Hooks ---

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
	if tx.Statement.Changed("History") {
		return lib.Error.Appointment.HistoryManualUpdateForbidden
	}

	var changes []FieldChange

	if a.History.IsEmpty() {
		a.History = AppointmentHistory{FieldChanges: []FieldChange{}}
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
			changes = append(changes, FieldChange{
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
