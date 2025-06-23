// filter/schedule_filter.go
package ScheduleFilter

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ScheduleFilter is the main service for handling schedule-related queries.
type schedule_filter struct {
	tx     *gorm.DB
	params *ScheduleQueryParams
}

// New creates a new instance of the schedule filter.
func NewFromContext(tx *gorm.DB, c *fiber.Ctx) (*schedule_filter, error) {
	get := c.Query("get")
	branchIDStr := c.Query("branch_id")
	employeeIDStr := c.Query("employee_id")
	serviceIDStr := c.Query("service_id")
	start_timeStr := c.Query("start_time")

	params, err := NewScheduleQueryParams(get, branchIDStr, employeeIDStr, serviceIDStr, start_timeStr)
	if err != nil {
		return nil, err
	}

	return &schedule_filter{tx: tx, params: params}, nil
}

// ## GetScheduleOptions Requirements (Revised)
//
// GetScheduleOptions returns available scheduling options based on the provided query parameters.
//
// Input parameters are:
//   - branch_id: UUID of the branch to filter by. (Optional)
//   - employee_id: UUID of the employee to filter by. (Optional)
//   - service_id: UUID of the service to filter by. (Optional)
//   - start_time: Time in RFC3339 format (e.g., "2023-10-26T10:00:00Z" or "0001-01-01T10:00:00Z"). (Optional)
//       - If a specific date (e.g., "2023-10-26") is provided within the RFC3339 string,
//         queries will target that specific date and its corresponding weekday.
//         The time component will filter for that specific time on that date.
//       - If a generic or minimal date (e.g., "0001-01-01") is used in the RFC3339 string,
//         it implies that primarily the **time component** (e.g., "10:00:00Z") is of interest.
//         The system will then interpret this as:
//           - For general availability (services, branches, employees): "at this time on any day the employee works".
//           - For time slots: "slots starting at this time of day", potentially for a default date
//             (e.g., today) or a defined upcoming period, unless further restricted by other parameters.
//   - get: Specifies the type of data to retrieve (e.g., "services", "branches", "employees", "time_slots"). (Required)
//
// The response may include services, employees, or available time slots, depending on the filters used.
// All availability checks (unless otherwise stated) are based on employees’ work schedules.
//
// Behavior rules:
//
// 
// --- If GET is "services" ---
// 
//   - All parameters (except 'get') nil: returns error.
//
//   // Branch combinations for "services"
//   - branchID:
//       Returns services available at the specified branch.
//   - branchID + employeeID:
//       Returns services offered by the employee at that branch.
//   - branchID + start_time:
//       - If start_time specifies a date: returns services available at the branch on the
//         weekday of the given date, considering employee availability for that weekday.
//       - If start_time primarily specifies a time-of-day: returns services available at the
//         branch that can start at the specified time-of-day, on any day employees
//         are scheduled to work.
//   - branchID + employeeID + start_time:
//       - If start_time specifies a date: returns services offered by the employee at the branch
//         on the weekday of the given date, considering their schedule for that specific time.
//       - If start_time primarily specifies a time-of-day: returns services offered by the
//         employee at that branch that can start at the specified time-of-day, on any day
//         they are scheduled to work.
//
//   // Employee combinations for "services" (without branch)
//   - employeeID:
//       Returns services offered by the employee, regardless of branch or time.
//   - employeeID + start_time:
//       - If start_time specifies a date: returns services offered by the employee on the
//         weekday of the given date, considering their schedule.
//       - If start_time primarily specifies a time-of-day: returns services offered by the
//         employee that they can perform at the specified time-of-day, on any day
//         they are scheduled to work, regardless of branch.
//
//   // Time Start combination for "services" (without branch or employee)
//   - start_time:
//       - If start_time specifies a date: returns services available across the company on
//         the weekday of the given date, considering employee schedules for that day and time.
//       - If start_time primarily specifies a time-of-day: returns services available across
//         the company that can start at the specified time-of-day, on any day employees
//         are scheduled to work.
//
// 
// --- If GET is "branches" ---
// 
//   - All parameters (except 'get') nil: returns all branches of the company.
//
//   // Employee combinations for "branches"
//   - employeeID:
//       Returns branches where the employee is assigned.
//   - employeeID + start_time:
//       - If start_time specifies a date: returns branches where the employee is available
//         on the weekday of the given date and at the specified time.
//       - If start_time primarily specifies a time-of-day: returns branches where the employee
//         is available to start work at the specified time-of-day, on any day they are scheduled.
//
//   // Service combinations for "branches"
//   - serviceID:
//       Returns branches that offer the given service.
//   - serviceID + start_time:
//       - If start_time specifies a date: returns branches that offer the service and have
//         employees available on the weekday of the given date and at the specified time.
//       - If start_time primarily specifies a time-of-day: returns branches that offer the
//         service and have employees available to start at the specified time-of-day,
//         on any day they are scheduled.
//
//   // Time Start combination for "branches"
//   - start_time:
//       - If start_time specifies a date: returns branches that have employees available on
//         the weekday of the given date and at the specified time.
//       - If start_time primarily specifies a time-of-day: returns branches that have
//         employees available to start at the specified time-of-day, on any day they are scheduled.
//
// 
// --- If GET is "employees" ---
//
//   - All parameters (except 'get') nil: returns all employees of the company.
//
//   // Branch combinations for "employees"
//   - branchID:
//       Returns employees assigned to the specified branch.
//   - branchID + serviceID:
//       Returns employees at the branch who offer the given service.
//   - branchID + start_time:
//       - If start_time specifies a date: returns employees at the branch available on the
//         weekday of the given date and at the specified time.
//       - If start_time primarily specifies a time-of-day: returns employees at the branch
//         available to start at the specified time-of-day, on any day they are scheduled.
//
//   // Service combinations for "employees"
//   - serviceID:
//       Returns employees who offer the service, regardless of branch.
//   - serviceID + start_time:
//       - If start_time specifies a date: returns employees who offer the service and are
//         available on the weekday of the given date and at the specified time.
//       - If start_time primarily specifies a time-of-day: returns employees who offer the
//         service and are available to start at the specified time-of-day, on any day
//         they are scheduled.
//
//   // Time Start combination for "employees"
//   - start_time:
//       - If start_time specifies a date: returns employees available on the weekday of the
//         given date and at the specified time.
//       - If start_time primarily specifies a time-of-day: returns employees available to
//         start at the specified time-of-day, regardless of branch or service,
//         on any day they are scheduled.
//
// 
// --- If GET is "time_slots" ---
// 
//   - All parameters (except 'get') nil: returns error.
//
//   - General Note for Time Slots:
//     When `start_time` is provided with a specific, meaningful date, the date component of
//     `start_time` will determine the **specific day** for which time slots are generated.
//     If `start_time` is used to indicate only a time-of-day is of primary interest
//     (e.g., by using a generic date like "0001-01-01T10:00:00Z"), slots will be
//     generated for that specific time-of-day. These slots will be sought within an
//     upcoming period (e.g., the next 30 days starting from today, by default).
//     This means it will return all time slots available that start at the specified
//     time-of-day within this defined upcoming period.
//
//
//   // Branch combinations for "time_slots"
//   - branchID:
//       Returns all time slots available at the branch for a default upcoming period
//       (e.g., today or next 7 days).
//   - branchID + serviceID:
//       Returns time slots at the branch for the specified service for a default upcoming period.
//   - branchID + employeeID:
//       Returns time slots at the branch for the specified employee for a default upcoming period.
//   - branchID + start_time:
//       Returns time slots available at the branch on the specific date derived from `start_time`.
//       If `start_time` also has a specific time component, results are filtered to slots
//       starting at/around that time.
//   - branchID + serviceID + start_time:
//       Returns time slots at the branch for the service that can start on the date and
//       at/around the time specified by `start_time`.
//   - branchID + employeeID + start_time:
//       Returns time slots at the branch for the employee that can start on the date and
//       at/around the time specified by `start_time`.
//
//   // Employee combinations for "time_slots"
//   - employeeID:
//       Returns all available time slots for the employee for a default upcoming period.
//   - employeeID + start_time:
//       Returns time slots the employee is available for on the date and at/around the time
//       specified by `start_time`.
//
//   // Service combinations for "time_slots"
//   - serviceID:
//       Returns all time slots available for the service for a default upcoming period.
//   - serviceID + start_time:
//       Returns time slots available for the service on the date and at/around the time
//       specified by `start_time`.
//
//   // Time Start combination for "time_slots"
//   - start_time:
//       Returns available time slots across the company that can start on the date and
//       at/around the time specified by `start_time`.
//
// The output dynamically adapts to the query context, providing only relevant scheduling options.
func (sf *schedule_filter) GetScheduleOptions() (any, error) {

	if sf.params == nil {
		return nil, fmt.Errorf("schedule filter parameters are not set")
	}

	var strategy ScheduleStrategy

	switch sf.params.Get {
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
		return nil, sf.params.toLibError(fmt.Errorf("unsupported 'get' type: %s", sf.params.Get))
	}

	return strategy.Fetch(sf.params)
}
