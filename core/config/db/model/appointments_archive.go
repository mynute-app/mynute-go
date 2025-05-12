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

var AppointmentArchiveTableName = "appointments_archive"

func (AppointmentArchive) TableName() string { return AppointmentArchiveTableName }

func (AppointmentArchive) Indexes() map[string]string {
	return AppointmentIndexes(AppointmentArchiveTableName)
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
