package model

import (
	"agenda-kaki-go/core/lib"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppointmentArchive struct {
	BaseModel
	AppointmentBase
	AppointmentJson
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
