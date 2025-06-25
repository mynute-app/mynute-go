package DTO

import (
	"time"

	"github.com/google/uuid"
)

type CreateBranchWorkSchedule struct {
	WorkRanges []CreateBranchWorkRange `json:"branch_work_ranges"`
}

type CreateBranchWorkRange struct {
	BranchID  uuid.UUID    `json:"branch_id" example:"00000000-0000-0000-0000-000000000000"`                // Branch ID
	Weekday   time.Weekday `json:"weekday" example:"1"`                                                     // Weekday (0 = Sunday, 1 = Monday, ..., 6 = Saturday)
	StartTime string       `json:"start_time" example:"09:00" format:"HH:mm"`                               // Start time (date ignored)
	EndTime   string       `json:"end_time" example:"17:00" format:"HH:mm"`                                 // End time (date ignored)
	TimeZone  string       `json:"timezone" example:"America/New_York"`                                     // Timezone in IANA format, e.g., "America/New_York"
	Services  []ServiceID  `json:"services" example:"[{\"id\": \"00000000-0000-0000-0000-000000000000\"}]"` // List of services associated with the work range
}

type BranchWorkRange struct {
	ID        string       `json:"id" example:"00000000-0000-0000-0000-000000000000"` // Work range ID
	CreateBranchWorkRange
}
