package ScheduleFilter

import (
	"fmt"
	"strings"
	"time"

	"mynute-go/core/lib"

	"github.com/google/uuid"
)

type GetType string

const (
	GetServices  GetType = "services"
	GetBranches  GetType = "branches"
	GetEmployees GetType = "employees"
	GetTimeSlots GetType = "time_slots"
	// GetInvalid GetType = "invalid" // Not explicitly used if validation is thorough
)

// ScheduleQueryParams holds all the validated and parsed input parameters for schedule filtering.
type ScheduleQueryParams struct {
	BranchID   *uuid.UUID
	EmployeeID *uuid.UUID
	ServiceID  *uuid.UUID
	Get        GetType

	// OriginalStartTime is the raw input from the query.
	OriginalStartTime *time.Time

	// Derived fields for easier processing
	IsSpecificDateQuery bool         // True if OriginalStartTime has a non-generic date.
	IsTimeOfDayQuery    bool         // True if OriginalStartTime implies time-of-day only (e.g., year is 1).
	QueryDate           *time.Time   // If IsSpecificDateQuery, this is the date part (time at 00:00:00).
	QueryTimeOfDay      *time.Time   // If OriginalStartTime is present, this is its time component (date part normalized to 0001-01-01).
	QueryWeekday        time.Weekday // If IsSpecificDateQuery, this is the weekday.
}

func NewScheduleQueryParams(get, branchIDStr, employeeIDStr, serviceIDStr, timeStr string) (*ScheduleQueryParams, error) {
	p := &ScheduleQueryParams{}

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
		parsedTime, err := parseTime(timeStr)
		if err != nil {
			return nil, err
		}
		p.OriginalStartTime = &parsedTime
		if parsedTime.Year() == 1 && parsedTime.Month() == 1 && parsedTime.Day() == 1 {
			p.IsTimeOfDayQuery = true
			normalizedTimeOfDayUTC := time.Date(0, time.January, 1, parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), parsedTime.Nanosecond(), time.UTC)
			p.QueryTimeOfDay = &normalizedTimeOfDayUTC
		} else {
			p.IsSpecificDateQuery = true
			p.QueryWeekday = parsedTime.Weekday()
			datePart := time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(), 0, 0, 0, 0, parsedTime.Location())
			p.QueryDate = &datePart

			normalizedTimeOfDayUTC := time.Date(0, time.January, 1, parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), parsedTime.Nanosecond(), time.UTC)
			p.QueryTimeOfDay = &normalizedTimeOfDayUTC
		}
	}

	return p, nil
}

func parseUUID(idStr, fieldName string) (*uuid.UUID, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID format for %s: %w", fieldName, err)
	}
	return &id, nil
}

func parseTime(timeStr string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		// Attempt to parse if only time is provided, assuming a zero date
		// This is tricky because RFC3339 expects a full date-time.
		// The convention "0001-01-01THH:MM:SSZ" handles the "time-of-day" intent.
		return time.Time{}, fmt.Errorf("invalid time format for start_time, expected RFC3339 (e.g., \"2023-10-26T10:00:00Z\" or \"0001-01-01T10:00:00Z\"): %w", err)
	}
	return t, nil
}

func (p *ScheduleQueryParams) allExceptGetAreNil() bool {
	return p.BranchID == nil && p.EmployeeID == nil && p.ServiceID == nil && p.OriginalStartTime == nil
}

func (p *ScheduleQueryParams) toLibError(err error) error {
	return lib.Error.General.BadRequest.WithError(err)
}
