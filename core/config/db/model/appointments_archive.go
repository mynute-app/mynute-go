package model

import (
	"agenda-kaki-go/core/lib"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppointmentArchive struct {
	ID         uuid.UUID          `gorm:"type:uuid;primaryKey;<-:create" json:"id"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
	DeletedAt  gorm.DeletedAt     `gorm:"index" json:"deleted_at"`
	ServiceID  uuid.UUID          `gorm:"type:uuid;not null;index" json:"service_id"`
	Service    *Service           `gorm:"foreignKey:ServiceID;references:ID;constraint:OnDelete:CASCADE;"` // Using your Service type
	EmployeeID uuid.UUID          `gorm:"type:uuid;not null;index" json:"employee_id"`
	Employee   *Employee          `gorm:"foreignKey:EmployeeID;references:ID;constraint:OnDelete:CASCADE;"` // Using your Employee type
	ClientID   uuid.UUID          `gorm:"type:uuid;not null;index" json:"client_id"`
	Client     *Client            `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE;"` // Using your Client type
	BranchID   uuid.UUID          `gorm:"type:uuid;not null;index" json:"branch_id"`
	Branch     *Branch            `gorm:"foreignKey:BranchID;references:ID;constraint:OnDelete:CASCADE;"` // Using your Branch type
	CompanyID  uuid.UUID          `gorm:"type:uuid;not null;index" json:"company_id"`
	Company    *Company           `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;" json:"-"` // Company loaded via FK, json:"-" often good practice
	StartTime  time.Time          `gorm:"not null;index" json:"start_time"`
	EndTime    time.Time          `gorm:"not null;index" json:"end_time"`
	Cancelled  bool               `gorm:"index;default:false" json:"cancelled"`
	History    AppointmentHistory `gorm:"type:jsonb" json:"history"` // JSONB field for history changes
}

func (AppointmentArchive) TableName() string { return "appointments_archive" }

func (AppointmentArchive) Indexes() map[string]string {
	return map[string]string{
		"idx_employee_time_active": "CREATE INDEX IF NOT EXISTS idx_employee_time_active ON appointments (employee_id, start_time, end_time, cancelled)",
		"idx_client_time_active":   "CREATE INDEX IF NOT EXISTS idx_client_time_active ON appointments (client_id, start_time, end_time, cancelled)",
		"idx_branch_time_active":   "CREATE INDEX IF NOT EXISTS idx_branch_time_active ON appointments (branch_id, start_time, end_time, cancelled)",
		"idx_company_active":       "CREATE INDEX IF NOT EXISTS idx_company_active ON appointments (company_id, cancelled)",
		"idx_start_time_active":    "CREATE INDEX IF NOT EXISTS idx_start_time_active ON appointments (start_time, cancelled)",
	}
}

// --- Appointments Archive Hooks ---

func (a *AppointmentArchive) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		return lib.Error.AppointmentArchive.IdNotSet
	}
	return nil
}

func (a *AppointmentArchive) BeforeUpdate(tx *gorm.DB) (err error) {
	return lib.Error.AppointmentArchive.UpdateForbidden
}

func (a *AppointmentArchive) BeforeDelete(tx *gorm.DB) (err error) {
	return lib.Error.AppointmentArchive.DeleteForbidden
}
