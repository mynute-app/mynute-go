// filter/schedule_filter.go
package schedule_filter

import (
	"fmt"

	"gorm.io/gorm"
)

// ScheduleFilter is the main service for handling schedule-related queries.
type ScheduleFilter struct {
	tx *gorm.DB
}

// NewScheduleFilter creates a new instance of the schedule filter.
func NewScheduleFilter(tx *gorm.DB) *ScheduleFilter {
	return &ScheduleFilter{tx: tx}
}

// GetScheduleOptions is the primary entry point. It selects the correct strategy
// based on the query parameters and executes it.
func (sf *ScheduleFilter) GetScheduleOptions(
	get, branchIDStr, employeeIDStr, serviceIDStr, weekday, timeStr string,
) (any, error) {

	params, err := NewScheduleQueryParams(get, branchIDStr, employeeIDStr, serviceIDStr, weekday, timeStr)
	if err != nil {
		return nil, params.toLibError(err)
	}

	var strategy ScheduleStrategy

	switch params.Get {
	case GetServices:
		strategy = &ServicesStrategy{tx: sf.tx}
	case GetBranches:
		strategy = &BranchesStrategy{tx: sf.tx}
	case GetEmployees:
		strategy = &EmployeesStrategy{tx: sf.tx}
	case GetTimeSlots:
		strategy = &TimeSlotsStrategy{tx: sf.tx}
	default:
		// This case is already handled by NewScheduleQueryParams, but as a safeguard:
		return nil, params.toLibError(fmt.Errorf("unsupported 'get' type: %s", params.Get))
	}

	return strategy.Fetch(params)
}