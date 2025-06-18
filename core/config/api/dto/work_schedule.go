package DTO

type CreateWorkSchedule struct {
	EmployeeID string    `json:"employee_id" example:"00000000-0000-0000-0000-000000000000"`              // Employee ID
	BranchID   string    `json:"branch_id" example:"00000000-0000-0000-0000-000000000000"`                // Branch ID
	Weekday    uint8     `json:"weekday" example:"1"`                                                     // Weekday (0 = Sunday, 1 = Monday, ..., 6 = Saturday)
	StartTime  string    `json:"start_time" example:"09:00" format:"HH:mm"`                               // Start time (date ignored)
	EndTime    string    `json:"end_time" example:"17:00" format:"HH:mm"`                                 // End time (date ignored)// End time
	TimeZone   string    `json:"timezone" example:"America/New_York"`                                     // Timezone in IANA format, e.g., "America/New_York"
	Services   []Service `json:"services" example:"[{\"id\": \"00000000-0000-0000-0000-000000000000\"}]"` // List of services associated with the work schedule
}
