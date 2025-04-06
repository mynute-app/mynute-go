package DTO

import "github.com/google/uuid"

type CreateAppointment struct {
	ServiceID  uuid.UUID   `json:"service_id" example:"1"`
	EmployeeID uuid.UUID   `json:"employee_id" example:"1"`
	ClientID   uuid.UUID   `json:"client_id" example:"1"`
	BranchID   uuid.UUID   `json:"branch_id" example:"1"`
	CompanyID  uuid.UUID   `json:"company_id" example:"1"`
	StartTime  string `json:"start_time" example:"2028-01-01T09:00:00Z"`
}

type UpdateAppointment struct {
	ID        uuid.UUID   `json:"id" example:"1"`
	StartTime string `json:"start_time" example:"2028-01-01T09:00:00Z"`
}

type Appointment struct {
	ID                uuid.UUID   `json:"id" example:"1"`
	ServiceID         uuid.UUID   `json:"service_id" example:"1"`
	EmployeeID        uuid.UUID   `json:"employee_id" example:"1"`
	ClientID          uuid.UUID   `json:"client_id" example:"1"`
	BranchID          uuid.UUID   `json:"branch_id" example:"1"`
	CompanyID         uuid.UUID   `json:"company_id" example:"1"`
	StartTime         string `json:"start_time" example:"2021-01-01T09:00:00Z"`
	EndTime           string `json:"end_time" example:"2021-01-01T10:00:00Z"`
	Rescheduled       bool   `json:"rescheduled" example:"false"`
	Cancelled         bool   `json:"cancelled" example:"false"`
	RescheduledToID   *uuid.UUID  `json:"rescheduled_to_id" example:"1"`
	RescheduledFromID *uuid.UUID  `json:"rescheduled_from_id" example:"1"`
}
