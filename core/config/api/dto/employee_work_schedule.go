package DTO

import "github.com/google/uuid"

//	@description	represents a work schedule for an employee, which is a collection of work ranges.
//	@name			CreateEmployeeWorkSchedule
//	@tag.name		dto.workrange.create_work_schedule
type EmployeeWorkSchedule struct {
	WorkRanges []EmployeeWorkRange `json:"employee_work_ranges"`
}

//	@description	represents the data required to create a work schedule for an employee.
//	@name			CreateEmployeeWorkSchedule
//	@tag.name		dto.workrange.create_work_schedule
type CreateEmployeeWorkSchedule struct {
	WorkRanges []CreateEmployeeWorkRange `json:"employee_work_ranges"` // List of work ranges for the employee
}

//	@description	represents the data required to create a work range for an employee.
//	@name			CreateEmployeeWorkRange
//	@tag.name		dto.workrange.create_work_range
type CreateEmployeeWorkRange struct {
	EmployeeID uuid.UUID   `json:"employee_id" example:"00000000-0000-0000-0000-000000000000"` // Employee ID
	BranchID   uuid.UUID   `json:"branch_id" example:"00000000-0000-0000-0000-000000000000"`   // Branch ID
	Weekday    uint8       `json:"weekday" example:"1"`                                        // Weekday (0 = Sunday, 1 = Monday, ..., 6 = Saturday)
	StartTime  string      `json:"start_time" example:"09:00" format:"HH:mm"`                  // Start time (date ignored)
	EndTime    string      `json:"end_time" example:"17:00" format:"HH:mm"`                    // End time (date ignored)
	TimeZone   string      `json:"time_zone" example:"America/New_York"`                       // Timezone in IANA format, e.g., "America/New_York"
	EmployeeWorkRangeServices
}

type UpdateWorkRange struct {
	Weekday   uint8  `json:"weekday" example:"1"`                       // Weekday (0 = Sunday, 1 = Monday, ..., 6 = Saturday)
	StartTime string `json:"start_time" example:"09:00" format:"HH:mm"` // Start time (date ignored)
	EndTime   string `json:"end_time" example:"17:00" format:"HH:mm"`   // End time (date ignored)
	TimeZone  string `json:"time_zone" example:"America/New_York"`      // Timezone in IANA format, e.g., "America/New_York"
}

//	@description	represents a work range for an employee, including its ID and the data required to create it.
//	@name			EmployeeWorkRange
//	@tag.name		dto.workrange.full
type EmployeeWorkRange struct {
	ID         string        `json:"id" example:"00000000-0000-0000-0000-000000000000"`          // Work range ID
	EmployeeID uuid.UUID     `json:"employee_id" example:"00000000-0000-0000-0000-000000000000"` // Employee ID
	BranchID   uuid.UUID     `json:"branch_id" example:"00000000-0000-0000-0000-000000000000"`   // Branch ID
	Weekday    uint8         `json:"weekday" example:"1"`                                        // Weekday (0 = Sunday, 1 = Monday, ..., 6 = Saturday)
	StartTime  string        `json:"start_time" example:"09:00" format:"HH:mm"`                  // Start time (date ignored)
	EndTime    string        `json:"end_time" example:"17:00" format:"HH:mm"`                    // End time (date ignored)
	TimeZone   string        `json:"time_zone" example:"America/New_York"`                       // Timezone in IANA format, e.g., "America/New_York"
	Services   []ServiceBase `json:"services" swaggertype:"array,object"`                        // List of services associated with the work range
}

type EmployeeWorkRangeServices struct {
	Services []ServiceID `json:"services" swaggertype:"array,object"` // List of services associated with the work range
}
