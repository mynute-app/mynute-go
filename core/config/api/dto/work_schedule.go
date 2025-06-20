package DTO

// @description represents a work schedule for an employee, which is a collection of work ranges.
// @name CreateWorkSchedule
// @tag.name dto.workrange.create_work_schedule
type CreateWorkSchedule []CreateWorkRange

// @description represents the data required to create a work range for an employee.
// @name CreateWorkRange
// @tag.name dto.workrange.create_work_range
type CreateWorkRange struct {
	EmployeeID string    `json:"employee_id" example:"00000000-0000-0000-0000-000000000000"`              // Employee ID
	BranchID   string    `json:"branch_id" example:"00000000-0000-0000-0000-000000000000"`                // Branch ID
	Weekday    uint8     `json:"weekday" example:"1"`                                                     // Weekday (0 = Sunday, 1 = Monday, ..., 6 = Saturday)
	StartTime  string    `json:"start_time" example:"09:00" format:"HH:mm"`                               // Start time (date ignored)
	EndTime    string    `json:"end_time" example:"17:00" format:"HH:mm"`                                 // End time (date ignored)
	TimeZone   string    `json:"timezone" example:"America/New_York"`                                     // Timezone in IANA format, e.g., "America/New_York"
	Services   []Service `json:"services" example:"[{\"id\": \"00000000-0000-0000-0000-000000000000\"}]"` // List of services associated with the work range
}

// @description represents a work range for an employee, including its ID and the data required to create it.
// @name WorkRange
// @tag.name dto.workrange.full
type WorkRange struct {
	ID string `json:"id" example:"00000000-0000-0000-0000-000000000000"` // Work range ID
	CreateWorkRange
}
