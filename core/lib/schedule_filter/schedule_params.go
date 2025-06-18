// filter/schedule_params.go
package schedule_filter

import (
	"fmt"
	"strings"
	"time"

	"agenda-kaki-go/core/lib"

	"github.com/google/uuid"
)

// GetType defines the valid types for the 'get' parameter.
type GetType string

const (
	GetServices   GetType = "services"
	GetBranches   GetType = "branches"
	GetEmployees  GetType = "employees"
	GetTimeSlots  GetType = "time_slots"
	GetInvalid    GetType = "invalid"
)

// ScheduleQueryParams holds all the validated and parsed input parameters for schedule filtering.
type ScheduleQueryParams struct {
	BranchID   *uuid.UUID
	EmployeeID *uuid.UUID
	ServiceID  *uuid.UUID
	Weekday    *string
	StartTime  *time.Time
	Get        GetType
}

// NewScheduleQueryParams creates and validates a new params object from raw string inputs.
func NewScheduleQueryParams(get, branchIDStr, employeeIDStr, serviceIDStr, weekday, timeStr string) (*ScheduleQueryParams, error) {
	p := &ScheduleQueryParams{}

	// Validate and set 'get' parameter
	switch GetType(strings.ToLower(get)) {
	case GetServices:
		p.Get = GetServices
	case GetBranches:
		p.Get = GetBranches
	case GetEmployees:
		p.Get = GetEmployees
	case GetTimeSlots:
		p.Get = GetTimeSlots
	default:
		return nil, fmt.Errorf("invalid 'get' parameter: %s. Must be one of 'services', 'branches', 'employees', 'time_slots'", get)
	}

	var err error
	if branchIDStr != "" {
		if p.BranchID, err = parseUUID(branchIDStr, "branch_id"); err != nil {
			return nil, err
		}
	}
	if employeeIDStr != "" {
		if p.EmployeeID, err = parseUUID(employeeIDStr, "employee_id"); err != nil {
			return nil, err
		}
	}
	if serviceIDStr != "" {
		if p.ServiceID, err = parseUUID(serviceIDStr, "service_id"); err != nil {
			return nil, err
		}
	}
	if timeStr != "" {
		if p.StartTime, err = parseTime(timeStr); err != nil {
			return nil, err
		}
	}
	if weekday != "" {
		// Normalize weekday to lowercase for consistency in JSON queries
		lowerWeekday := strings.ToLower(weekday)
		p.Weekday = &lowerWeekday
	}
	return p, nil
}

// Helper functions for parsing
func parseUUID(idStr, fieldName string) (*uuid.UUID, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID format for %s: %w", fieldName, err)
	}
	return &id, nil
}

func parseTime(timeStr string) (*time.Time, error) {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid time format, expected RFC3339: %w", err)
	}
	return &t, nil
}

// allExceptGetAreNil checks if all filter parameters (other than 'get') are nil.
func (p *ScheduleQueryParams) allExceptGetAreNil() bool {
	return p.BranchID == nil && p.EmployeeID == nil && p.ServiceID == nil && p.Weekday == nil && p.StartTime == nil
}

func (p *ScheduleQueryParams) toLibError(err error) error {
	return lib.Error.General.BadRequest.WithError(err)
}