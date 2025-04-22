package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppointmentHistoryAction string
type ContextKey string

const UserIDContextKey ContextKey = "currentUserID"

var (
	ActionUpdate AppointmentHistoryAction = "Update"
	ActionCreate AppointmentHistoryAction = "Create"
	ActionView   AppointmentHistoryAction = "View"
	ActionCancel AppointmentHistoryAction = "Cancel"
)

// AppointmentHistory struct remains largely the same
type AppointmentHistory struct {
	BaseModel
	Timestamp            time.Time                `gorm:"not null;index;default:CURRENT_TIMESTAMP"`
	AppointmentID        uuid.UUID                `gorm:"type:uuid;not null;index"`
	Action               AppointmentHistoryAction `gorm:"type:varchar(50);not null"`
	UserID               *uuid.UUID               `gorm:"type:uuid"`
	Notes                string                   `gorm:"type:text"`
	ServiceID            uuid.UUID                `gorm:"type:uuid;not null"`
	EmployeeID           uuid.UUID                `gorm:"type:uuid;not null"`
	ClientID             uuid.UUID                `gorm:"type:uuid;not null"`
	BranchID             uuid.UUID                `gorm:"type:uuid;not null"`
	CompanyID            uuid.UUID                `gorm:"type:uuid;not null"`
	StartTime            time.Time                `gorm:"not null"`
	EndTime              time.Time                `gorm:"not null"`
	Cancelled            bool                     `gorm:"not null"`
	AppointmentCreatedAt time.Time
	AppointmentUpdatedAt time.Time
}

func (AppointmentHistory) TableName() string { return "appointment_history" }

// CreateAppointmentHistory Helper
func CreateAppointmentHistory(tx *gorm.DB, app *Appointment, action AppointmentHistoryAction, notes string) error {
	// --- Pre-checks ---
	if app == nil || app.ID == uuid.Nil || app.ServiceID == uuid.Nil || app.EmployeeID == uuid.Nil || app.ClientID == uuid.Nil || app.BranchID == uuid.Nil || app.CompanyID == uuid.Nil || app.StartTime.IsZero() || app.EndTime.IsZero() {
		err := errors.New("attempted to log history for appointment with incomplete data")
		fmt.Printf("ERROR: %v (ID: %s)\n", err, app.ID)
		// Return a standard error, hook will wrap it into the specific lib.ErrorStruct
		return err
	}

	// --- Get User ID from Context ---
	var userID *uuid.UUID
	ctxValue := tx.Statement.Context.Value(UserIDContextKey)
	if ctxValue != nil { /* ... same user ID retrieval logic ... */
	}

	history := AppointmentHistory{
		AppointmentID: app.ID, Action: action, UserID: userID, Notes: notes, Timestamp: time.Now(),
		ServiceID: app.ServiceID, EmployeeID: app.EmployeeID, ClientID: app.ClientID, BranchID: app.BranchID,
		CompanyID: app.CompanyID, StartTime: app.StartTime, EndTime: app.EndTime, Cancelled: app.Cancelled,
		AppointmentCreatedAt: app.CreatedAt, AppointmentUpdatedAt: app.UpdatedAt,
	}

	// Create the record using the hook's transaction
	if err := tx.Create(&history).Error; err != nil {
		// Return standard DB error, hook will wrap it
		return fmt.Errorf("db error creating appointment history record: %w", err)
	}
	return nil
}
