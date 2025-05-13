package mJSON

import (
	"time"

	"github.com/google/uuid"
)

type ClientAppointments []ClientAppointment

type ClientAppointment struct {
	ServiceName string    `json:"service_name"`
	CompanyID   uuid.UUID `json:"company_id"`
	StartTime   time.Time `json:"start_time"`
}
