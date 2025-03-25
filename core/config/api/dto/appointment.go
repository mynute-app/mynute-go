package DTO

type CreateAppointment struct {
	ServiceID  uint   `json:"service_id" example:"1"`
	EmployeeID uint   `json:"employee_id" example:"1"`
	ClientID   uint   `json:"client_id" example:"1"`
	BranchID   uint   `json:"branch_id" example:"1"`
	CompanyID  uint   `json:"company_id" example:"1"`
	StartTime  string `json:"start_time" example:"2028-01-01T09:00:00Z"`
}

type UpdateAppointment struct {
	ID        uint   `json:"id" example:"1"`
	StartTime string `json:"start_time" example:"2028-01-01T09:00:00Z"`
}

type Appointment struct {
	ID                uint   `json:"id" example:"1"`
	ServiceID         uint   `json:"service_id" example:"1"`
	EmployeeID        uint   `json:"employee_id" example:"1"`
	ClientID          uint   `json:"client_id" example:"1"`
	BranchID          uint   `json:"branch_id" example:"1"`
	CompanyID         uint   `json:"company_id" example:"1"`
	StartTime         string `json:"start_time" example:"2021-01-01T09:00:00Z"`
	EndTime           string `json:"end_time" example:"2021-01-01T10:00:00Z"`
	Rescheduled       bool   `json:"rescheduled" example:"false"`
	Cancelled         bool   `json:"cancelled" example:"false"`
	RescheduledToID   *uint  `json:"rescheduled_to_id" example:"1"`
	RescheduledFromID *uint  `json:"rescheduled_from_id" example:"1"`
}
