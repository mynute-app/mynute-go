package model

import (
	"agenda-kaki-go/core/lib"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientAppointment struct {
	BaseModel
	AppointmentBase
}

var ClientAppointmentTableName = "public.client_appointments"

func (ClientAppointment) TableName() string {
	return ClientAppointmentTableName
}

func (ClientAppointment) Indexes() map[string]string {
	return AppointmentIndexes(ClientAppointmentTableName)
}

func (c *ClientAppointment) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		return lib.Error.AppointmentArchive.IdNotSet
	}
	return nil
}

func (c *ClientAppointment) BeforeDelete(tx *gorm.DB) (err error) {
	return lib.Error.AppointmentArchive.DeleteForbidden
}
