package dJSON

import (
	"time"

	"github.com/google/uuid"
)

type ClientAppointments []ClientAppointment

type ClientAppointment struct {
	ServiceName   string    `json:"service_name" example:"Service name example"`
	CompanyName   string    `json:"company_name" example:"Company name example"`
	CompanyID     uuid.UUID `json:"company_id" example:"00000000-0000-0000-0000-000000000000"`
	StartTime     time.Time `json:"start_time" example:"2021-01-01T09:00:00Z"`
	BranchAddress string    `json:"branch_address" example:"76, Example street, My city, My country, 09090790"`
}
