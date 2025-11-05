package model

import (
	"errors"
	"fmt"
	"mynute-go/services/core/api/lib"
	mJSON "mynute-go/services/core/config/db/json"
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
	StartTime             time.Time  `gorm:"type:time;not null" json:"start_time"`
	EndTime               time.Time  `gorm:"type:time;not null" json:"end_time"`
	TimeZone              string     `gorm:"type:varchar(100);not null" json:"time_zone" validate:"required,myTimezoneValidation"` // Time zone in IANA format (e.g., "America/New_York", "America/Sao_Paulo", etc.)
	ActualStartTime       time.Time  `gorm:"type:time;not null" json:"actual_start_time"`
	ActualEndTime         time.Time  `gorm:"type:time;not null" json:"actual_end_time"`
	CancelTime            time.Time  `gorm:"type:time;not null" json:"cancel_time"`
	IsFulfilled           bool       `gorm:"default:false" json:"is_fulfilled"`
	IsCancelled           bool       `gorm:"default:false" json:"is_cancelled"`
	IsCancelledByClient   bool       `gorm:"default:false" json:"is_cancelled_by_client"`
	IsCancelledByEmployee bool       `gorm:"default:false" json:"is_cancelled_by_employee"`
	IsConfirmedByClient   bool       `gorm:"default:false" json:"is_confirmed_by_client"`
}

// This is the foreign key struct for the Appointment model at company schema level.
type AppointmentFK struct {
	Service  *Service  `gorm:"foreignKey:ServiceID;references:ID;constraint:OnDelete:CASCADE;"`      // Using your Service type
	Employee *Employee `gorm:"foreignKey:EmployeeID;references:UserID;constraint:OnDelete:CASCADE;"` // Using your Employee type (UserID is the primary key)
	Branch   *Branch   `gorm:"foreignKey:BranchID;references:ID;constraint:OnDelete:CASCADE;"`       // Using your Branch type
	Payment  *Payment  `gorm:"foreignKey:PaymentID;references:ID;constraint:OnDelete:CASCADE;"`      // Using your Payment type
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

const AppointmentTableName = "appointments"

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
	client.AddAppointment(a, tx)
	if err := tx.Save(&client).Error; err != nil {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("updating client: %w", err))
	}
	return nil
}

func (a *Appointment) BeforeCreate(tx *gorm.DB) error {
	if err := lib.MyCustomStructValidator(a); err != nil {
		return err
	}
	if !a.History.IsEmpty() {
		return lib.Error.Appointment.HistoryManualUpdateForbidden
	}
	if err := a.ValidateRules(tx, true); err != nil {
		return err
	}
	return nil
}

func (a *Appointment) BeforeUpdate(tx *gorm.DB) error {
	// Fetch the original record from the database using the ID from the `a` struct.
	var originalAppointment Appointment
	if err := tx.First(&originalAppointment, "id = ?", a.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lib.Error.Appointment.NotFound.WithError(fmt.Errorf("appointment ID %s", a.ID))
		}
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("loading appointment: %w", err))
	}

	if a.CompanyID != uuid.Nil && a.CompanyID != originalAppointment.CompanyID {
		return lib.Error.Appointment.UpdateFailed.WithError(fmt.Errorf("cannot change company ID"))
	} else if a.BranchID != uuid.Nil && a.BranchID != originalAppointment.BranchID {
		return lib.Error.Appointment.UpdateFailed.WithError(fmt.Errorf("cannot change branch ID"))
	} else if a.EmployeeID != uuid.Nil && a.EmployeeID != originalAppointment.EmployeeID {
		return lib.Error.Appointment.UpdateFailed.WithError(fmt.Errorf("cannot change employee ID"))
	} else if a.ServiceID != uuid.Nil && a.ServiceID != originalAppointment.ServiceID {
		return lib.Error.Appointment.UpdateFailed.WithError(fmt.Errorf("cannot change service ID"))
	}

	var changes []mJSON.FieldChange

	// --- Start: The Refactored Comparison Logic ---
	originalVal := reflect.ValueOf(originalAppointment)
	newVal := reflect.ValueOf(*a) // The incoming struct 'a' has the new data

	// Iterate through all fields of the Appointment struct
	for i := range newVal.NumField() {
		fieldStruct := newVal.Type().Field(i)
		fieldName := fieldStruct.Name

		// Skip fields that we don't want to track or compare
		if fieldName == "ID" || fieldName == "CreatedAt" || fieldName == "UpdatedAt" || fieldName == "History" || fieldName == "Comments" || fieldName == "BaseModel" {
			continue
		}

		// Prevent changing CompanyID
		if fieldName == "CompanyID" && originalAppointment.CompanyID != a.CompanyID && a.CompanyID != uuid.Nil {
			return lib.Error.General.UpdatedError.WithError(errors.New("the CompanyID cannot be changed after creation"))
		}

		originalFieldVal := originalVal.FieldByName(fieldName).Interface()
		newFieldVal := newVal.Field(i).Interface()

		// Check if the new value is a non-zero value for its type.
		// This aligns with GORM's struct-update behavior.
		isNewValueNonZero := !reflect.DeepEqual(newFieldVal, reflect.Zero(fieldStruct.Type).Interface())

		// If the new value is non-zero AND it's different from the original value...
		if isNewValueNonZero && !reflect.DeepEqual(originalFieldVal, newFieldVal) {
			changes = append(changes, mJSON.FieldChange{
				CreatedAt: time.Now(),
				Field:     fieldName,
				OldValue:  fmt.Sprintf("%v", originalFieldVal),
				NewValue:  fmt.Sprintf("%v", newFieldVal),
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

	// 1. Basic Required Fields & Time Checks

	if a.StartTime.IsZero() {
		return lib.Error.Appointment.StartTimeInThePast
	}
	if a.ServiceID == uuid.Nil || a.EmployeeID == uuid.Nil || a.ClientID == uuid.Nil || a.BranchID == uuid.Nil || a.CompanyID == uuid.Nil {
		return lib.Error.Appointment.MissingRequiredIDs
	}
	if isCreate || a.StartTime.IsZero() {
		if a.StartTime.Before(time.Now().Add(-1 * time.Minute)) {
			return lib.Error.Appointment.StartTimeInThePast // Use specific error for past time
		}
	} else if a.CompanyID != uuid.Nil {
		return lib.Error.Appointment.UpdateFailed.WithError(fmt.Errorf("can not change company id")) // Use specific error for company ID change
	}

	// 2. Calculate & Validate EndTime
	var serviceDuration uint
	if err := tx.Model(&Service{}).Where("id = ?", a.ServiceID).Pluck("duration", &serviceDuration).Error; err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("error loading service duration: %w", err))
	}
	if serviceDuration <= 0 { // Use uint duration from your model
		return lib.Error.Appointment.InvalidServiceDuration
	}

	a.EndTime = a.StartTime.Add(time.Duration(serviceDuration) * time.Minute)
	if !a.EndTime.After(a.StartTime) {
		return lib.Error.Appointment.EndTimeBeforeStart
	}

	// --- If being cancelled, validation stops here --- //
	if a.IsCancelled {
		return nil
	}

	// 3. Relationship & Existence Checks (Use loaded structs)

	// Check if Branch belongs to the same Company as the appointment
	var aBranchCompanyID string
	if err := tx.Model(&Branch{}).Where("id = ?", a.BranchID).Pluck("company_id", &aBranchCompanyID).Error; err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("error loading branch company ID: %w", err))
	}
	if aBranchCompanyID != a.CompanyID.String() {
		return lib.Error.Company.BranchDoesNotBelong
	}

	// Check if Employee belongs to the same Company as the Appointment
	var aEmployeeCompanyID string
	if err := tx.Model(&Employee{}).Where("id = ?", a.EmployeeID).Pluck("company_id", &aEmployeeCompanyID).Error; err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("error loading employee company ID: %w", err))
	}
	if aEmployeeCompanyID != a.CompanyID.String() {
		return lib.Error.Company.EmployeeDoesNotBelong
	}

	// Check if Service belongs to the same Company as the Appointment
	var aServiceCompanyID string
	if err := tx.Model(&Service{}).Where("id = ?", a.ServiceID).Pluck("company_id", &aServiceCompanyID).Error; err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("error loading service company ID: %w", err))
	}
	if aServiceCompanyID != a.CompanyID.String() {
		return lib.Error.Company.ServiceDoesNotBelong
	}

	var count int64

	// Check if service belongs to the branch
	count = 0
	tx.Table("branch_services").Where("branch_id = ? AND service_id = ?", a.BranchID, a.ServiceID).Count(&count)
	if count == 0 {
		return lib.Error.Branch.ServiceDoesNotBelong
	}

	// Check if employee offers the service
	count = 0
	tx.Table("employee_services").Where("employee_id = ? AND service_id = ?", a.EmployeeID, a.ServiceID).Count(&count)
	if count == 0 {
		return lib.Error.Employee.ServiceDoesNotBelong
	}

	// Check if employee works at the branch
	count = 0
	tx.Table("employee_branches").Where("employee_id = ? AND branch_id = ?", a.EmployeeID, a.BranchID).Count(&count)
	if count == 0 {
		return lib.Error.Employee.BranchDoesNotBelong
	}

	// 4. Check Employee Availability (Work Schedule)
	weekday := a.StartTime.Weekday()
	aStartTimeHHMMUTC := time.Date(2020, 1, 1, a.StartTime.Hour(), a.StartTime.Minute(), 0, 0, a.StartTime.Location()).UTC()
	aEndTimeHHMMUTC := time.Date(2020, 1, 1, a.EndTime.Hour(), a.EndTime.Minute(), 0, 0, a.EndTime.Location()).UTC()
	if err := tx.Model(&EmployeeWorkRange{}).
		Where("employee_id = ? AND branch_id = ?", a.EmployeeID, a.BranchID).
		Where("start_time <= ? AND end_time >= ? AND weekday = ?", aStartTimeHHMMUTC, aEndTimeHHMMUTC, weekday).
		Count(&count).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("no work schedule was found that could contain the appointment from (%s) to (%s) on day (%s) for employee %s at branch %s", a.StartTime.Format(time.RFC3339), a.EndTime.Format(time.RFC3339), weekday, a.EmployeeID, a.BranchID))
		}
		return lib.Error.General.InternalError.WithError(fmt.Errorf("error querying work schedule: %w", err))
	}

	if count == 0 {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("no work schedule was found that could contain the appointment from (%s) to (%s) on day (%s) for employee %s at branch %s", a.StartTime.Format(time.RFC3339), a.EndTime.Format(time.RFC3339), weekday, a.EmployeeID, a.BranchID))
	}

	ChangeSchema := func(schema string) error {
		switch schema {
		case "public":
			if err := lib.ChangeToPublicSchema(tx); err != nil {
				return lib.Error.General.InternalError.WithError(fmt.Errorf("error changing to public schema: %w", err))
			}
		case "company":
			companySchema := fmt.Sprintf("company_%s", a.CompanyID.String())
			if err := lib.ChangeToCompanySchema(tx, companySchema); err != nil {
				return lib.Error.General.InternalError.WithError(fmt.Errorf("error changing to company schema: %w", err))
			}
		}
		return nil
	}

	// 5. Overlap and Capacity Checks
	aStartTimeUTC := a.StartTime.UTC()
	aEndTimeUTC := a.EndTime.UTC()
	overlapTime := `? > start_time AND end_time > ?`
	notSameID := `id != ?`
	cancelled := `is_cancelled = ?`

	Query := func() *gorm.DB {
		return tx.Model(&Appointment{}).
			Where(cancelled, false).
			Where(notSameID, a.ID).
			Where(overlapTime, aEndTimeUTC, aStartTimeUTC)
	}

	// Employee Overlap and Capacities
	var employeeAppointmentsCount int64
	if err := Query().
		Where("employee_id = ?", a.EmployeeID).
		Count(&employeeAppointmentsCount).Error; err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("db error checking employee overlap: %w", err))
	}
	if employeeAppointmentsCount > 0 {
		var employeeTotalServiceDensity int64 // As here it is a field with default set to -1 by gorm tags we don't need to initialize it to -1
		if err := tx.Model(&Employee{}).Where("id = ?", a.EmployeeID).Pluck("total_service_density", &employeeTotalServiceDensity).Error; err != nil {
			return lib.Error.General.InternalError.WithError(fmt.Errorf("error loading employee density: %w", err))
		}
		if employeeAppointmentsCount >= employeeTotalServiceDensity {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee %s has reached its maximum density of %d appointments", a.EmployeeID, employeeTotalServiceDensity))
		}
		serviceDensityForTheEmployee := int64(-1) // Force default to -1 if not found (gorm defaults to 0). It will avoid triggering the subsequent error involuntarily
		if err := tx.Model(&EmployeeServiceDensity{}).Where("employee_id = ? AND service_id = ?", a.EmployeeID, a.ServiceID).Pluck("density", &serviceDensityForTheEmployee).Error; err != nil {
			return lib.Error.General.InternalError.WithError(fmt.Errorf("error loading employee service density: %w", err))
		}
		if serviceDensityForTheEmployee >= 0 && employeeAppointmentsCount >= serviceDensityForTheEmployee {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee %s has reached its maximum density of %d appointments for service (%s)", a.EmployeeID, serviceDensityForTheEmployee, a.ServiceID))
		}
	}

	// Branch Overlap and Capacities
	var branchAppointmentsCount int64
	if err := Query().
		Where("branch_id = ?", a.BranchID).
		Count(&branchAppointmentsCount).Error; err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("db error checking branch overlap: %w", err))
	}
	if branchAppointmentsCount > 0 {
		var branchTotalServiceDensity int32 // As here it is a field with default set to -1 by gorm tags we don't need to initialize it to -1
		if err := tx.Model(&Branch{}).Where("id = ?", a.BranchID).Pluck("total_service_density", &branchTotalServiceDensity).Error; err != nil {
			return lib.Error.General.InternalError.WithError(fmt.Errorf("error loading branch density: %w", err))
		}
		if branchTotalServiceDensity >= 0 && branchAppointmentsCount >= int64(branchTotalServiceDensity) {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch %s has reached its maximum density of %d appointments", a.BranchID, branchTotalServiceDensity))
		}
		serviceDensityForTheBranch := int64(-1) // Force default to -1 if not found (gorm defaults to 0). It will avoid triggering the subsequent error involuntarily
		if err := tx.Model(&BranchServiceDensity{}).Where("branch_id = ? AND service_id = ?", a.BranchID, a.ServiceID).Pluck("density", &serviceDensityForTheBranch).Error; err != nil {
			return lib.Error.General.InternalError.WithError(fmt.Errorf("error loading branch service density: %w", err))
		}
		if serviceDensityForTheBranch >= 0 && branchAppointmentsCount >= serviceDensityForTheBranch {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("branch %s has reached its maximum density of %d appointments for service (%s)", a.BranchID, serviceDensityForTheBranch, a.ServiceID))
		}
	}

	// Client Overlap (Under Company Schema Search)
	var clientAppointmentsCount int64
	if err := Query().
		Where("client_id = ?", a.ClientID).
		Count(&clientAppointmentsCount).Error; err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("db error checking client overlap: %w", err))
	}
	if clientAppointmentsCount > 0 {
		return lib.Error.Client.ScheduleConflict
	}

	// Client Overlap (Under Public Schema Search)
	if err := ChangeSchema("public"); err != nil {
		return err
	}
	clientAppointmentsCount = 0
	if err := tx.Model(&ClientAppointment{}).
		Where("client_id = ?", a.ClientID).
		Where("appointment_id != ?", a.ID).
		Where("company_id != ?", a.CompanyID).
		Where(cancelled, false).
		Where(overlapTime, aEndTimeUTC, aStartTimeUTC).
		Count(&clientAppointmentsCount).Error; err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("db error checking client overlap: %w", err))
	}
	if clientAppointmentsCount > 0 {
		return lib.Error.Client.ScheduleConflict
	}
	if err := ChangeSchema("company"); err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("error changing to company schema: %w", err))
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
	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("error changing to public schema: %w", err))
	}
	var clientAppointment ClientAppointment
	if err := tx.Model(&ClientAppointment{}).
		Where("appointment_id = ?", a.ID).
		First(&clientAppointment).Error; err != nil {
		// TODO: Here is a critical point to discuss: if the appointment is
		// cancelled but the ClientAppointment record is missing
		// we must send an alert to the admin team to investigate why it happened.
		// For now, we won't do anything.
		// This situation should never happen in a properly functioning system.
		if err == gorm.ErrRecordNotFound {
			return nil // If not found, nothing more to do
		}
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("error loading client appointment: %w", err))
	}
	clientAppointment.IsCancelled = true
	err = tx.Save(&clientAppointment).Error
	if err != nil {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("error cancelling client appointment: %w", err))
	}
	if err := lib.ChangeToCompanySchema(tx, a.CompanyID.String()); err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("error changing to company schema: %w", err))
	}
	return nil
}
