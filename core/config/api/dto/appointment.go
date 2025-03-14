package DTO

type Appointment struct {
	ID         uint   `json:"id"`
	ServiceID  uint   `json:"service_id"`
	EmployeeID uint   `json:"employee_id"`
	UserID     uint   `json:"user_id"`
	BranchID   uint   `json:"branch_id"`
	CompanyID  uint   `json:"company_id"`
	StartTime  string `json:"start_time" example:"2021-01-01T09:00:00Z"`
	EndTime    string `json:"end_time" example:"2021-01-01T10:00:00Z"`
}