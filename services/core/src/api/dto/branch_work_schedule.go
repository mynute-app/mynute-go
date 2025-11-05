package DTO

import (
	"github.com/google/uuid"
)

// @description	represents a work schedule for a branch, which is a collection of work ranges.
// @name			CreateBranchWorkSchedule
// @tag.name		dto.workrange.create_work_schedule
type BranchWorkSchedule struct {
	WorkRanges []BranchWorkRange `json:"branch_work_ranges"` // List of work ranges for the branch
}

// @description	represents the data required to create a work schedule for a branch.
// @name			CreateBranchWorkSchedule
// @tag.name		dto.workrange.create_work_schedule
type CreateBranchWorkSchedule struct {
	WorkRanges []CreateBranchWorkRange `json:"branch_work_ranges"`
}

// @description	represents the data required to create a work range for a branch.
// @name			CreateBranchWorkRange
// @tag.name		dto.workrange.create_work_range
type CreateBranchWorkRange struct {
	BranchID  uuid.UUID `json:"branch_id" example:"00000000-0000-0000-0000-000000000000"` // Branch ID
	Weekday   uint8     `json:"weekday" example:"1"`                                      // Weekday (0 = Sunday, 1 = Monday, ..., 6 = Saturday)
	StartTime string    `json:"start_time" example:"09:00" format:"HH:mm"`                // Start time (date ignored)
	EndTime   string    `json:"end_time" example:"17:00" format:"HH:mm"`                  // End time (date ignored)
	TimeZone  string    `json:"time_zone" example:"America/New_York"`                     // Timezone in IANA format, e.g., "America/New_York"
	BranchWorkRangeServices
}

// @description	represents a work range for a branch, including its ID and the data required to create it.
// @name			BranchWorkRange
// @tag.name		dto.workrange.full
type BranchWorkRange struct {
	ID string `json:"id" example:"00000000-0000-0000-0000-000000000000"` // Work range ID
	CreateBranchWorkRange
}

type BranchWorkRangeServices struct {
	Services []ServiceBase `json:"services" swaggertype:"array,object"` // List of services associated with the work range
}

