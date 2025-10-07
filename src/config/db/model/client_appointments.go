package model

import (
	"time"

	"github.com/google/uuid"
)

type ClientAppointment struct {
	AppointmentID uuid.UUID `gorm:"type:uuid;not null" json:"appointment_id"`
	ClientID      uuid.UUID `gorm:"type:uuid;not null" json:"client_id"`
	CompanyID     uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	StartTime     time.Time `gorm:"type:time;not null" json:"start_time"`
	EndTime       time.Time `gorm:"type:time;not null" json:"end_time"`
	TimeZone      string    `gorm:"type:varchar(100);not null" json:"time_zone" validate:"required,myTimezoneValidation"` // Time zone in IANA format (e.g., "America/New_York", "America/Sao_Paulo", etc.)
	IsCancelled   bool      `gorm:"default:false" json:"is_cancelled"`
}
