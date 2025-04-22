package DTO

import "github.com/google/uuid"

type CreateAppointment struct {
	ServiceID  uuid.UUID `json:"service_id" example:"00000000-0000-0000-0000-000000000000"`
	EmployeeID uuid.UUID `json:"employee_id" example:"00000000-0000-0000-0000-000000000000"`
	ClientID   uuid.UUID `json:"client_id" example:"00000000-0000-0000-0000-000000000000"`
	BranchID   uuid.UUID `json:"branch_id" example:"00000000-0000-0000-0000-000000000000"`
	CompanyID  uuid.UUID `json:"company_id" example:"00000000-0000-0000-0000-000000000000"`
	StartTime  string    `json:"start_time" example:"2028-01-01T09:00:00Z"`
}

type UpdateAppointment struct {
	ID        uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	StartTime string    `json:"start_time" example:"2028-01-01T09:00:00Z"`
}

type Appointment struct {
	ID                uuid.UUID  `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	ServiceID         uuid.UUID  `json:"service_id" example:"00000000-0000-0000-0000-000000000000"`
	EmployeeID        uuid.UUID  `json:"employee_id" example:"00000000-0000-0000-0000-000000000000"`
	ClientID          uuid.UUID  `json:"client_id" example:"00000000-0000-0000-0000-000000000000"`
	BranchID          uuid.UUID  `json:"branch_id" example:"00000000-0000-0000-0000-000000000000"`
	CompanyID         uuid.UUID  `json:"company_id" example:"00000000-0000-0000-0000-000000000000"`
	StartTime         string     `json:"start_time" example:"2021-01-01T09:00:00Z"`
	EndTime           string     `json:"end_time" example:"2021-01-01T10:00:00Z"`
	Rescheduled       bool       `json:"rescheduled" example:"false"`
	Cancelled         bool       `json:"cancelled" example:"false"`
	RescheduledToID   *uuid.UUID `json:"rescheduled_to_id" example:"00000000-0000-0000-0000-000000000000"`
	RescheduledFromID *uuid.UUID `json:"rescheduled_from_id" example:"00000000-0000-0000-0000-000000000000"`
}
