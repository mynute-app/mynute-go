package service

import (
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type schedule_filter struct {
	MyGorm     *handler.Gorm
	Ctx        *fiber.Ctx
	BranchID   *uuid.UUID
	EmployeeID *uuid.UUID
	ServiceID  *uuid.UUID
	Weekday    *time.Weekday
	TimeStr    *string
	Get        string
	Error      error
}

func (s *service) ScheduleFilter() *schedule_filter {
	var sf schedule_filter
	if s.Error != nil {
		sf.Error = s.Error
		return &sf
	}
	var (
		branchID, employeeID, serviceID *uuid.UUID
		weekday                         *time.Weekday
		timeStr                         *string
		get                             string
	)

	if val := s.Context.Query("get"); val != "" {
		get = val
	} else {
		sf.Error = lib.Error.General.BadRequest.WithError(errors.New("missing 'get' parameter"))
		return &sf
	}

	if val := s.Context.Query("branch_id"); val != "" {
		if id, err := uuid.Parse(val); err == nil {
			branchID = &id
		}
	}
	if val := s.Context.Query("employee_id"); val != "" {
		if id, err := uuid.Parse(val); err == nil {
			employeeID = &id
		}
	}
	if val := s.Context.Query("service_id"); val != "" {
		if id, err := uuid.Parse(val); err == nil {
			serviceID = &id
		}
	}
	if val := s.Context.Query("start_time"); val != "" {
		if parsed, err := time.Parse(time.RFC3339, val); err == nil {
			d := parsed.Weekday()
			s := parsed.Format("15:04")
			weekday = &d
			timeStr = &s
		} else {
			sf.Error = lib.Error.General.BadRequest.WithError(err)
			return &sf
		}
	}

	return &schedule_filter{
		MyGorm:     s.MyGorm,
		Ctx:        s.Context,
		BranchID:   branchID,
		EmployeeID: employeeID,
		ServiceID:  serviceID,
		Weekday:    weekday,
		TimeStr:    timeStr,
		Get:        get,
		Error:      s.Error,
	}
}

// GetScheduleOptions returns available scheduling options based on the provided query parameters.
//
// Input parameters are:
// - branch_id: UUID of the branch to filter by.
// - employee_id: UUID of the employee to filter by.
// - service_id: UUID of the service to filter by.
// - weekday: Day of the week to filter by (e.g., "Monday").
// - start_time: Time in RFC3339 format to filter by (e.g., "2023-10-01T10:00:00Z").
// - get: Specifies the type of data to retrieve (e.g., "services", "branches", "employees", "time_slots").
// The response may include services, employees, or available time slots, depending on the filters used.
// Behavior rules:
// - If GET is "services":
//   - All parameters nil: returns error.
//   // Branch //
//   - branchID: returns services available at the specified branch.
//   - branchID + employeeID: returns services offered by the employee at that branch.
//   - branchID + weekday: returns services available at the branch on the given weekday. (based on employees’ work schedules)
//   - branchID + timeStr: returns services available at the branch that can start at the specified time in any day of the week. (based on employees’ work schedules)
//   - branchID + weekday + timeStr: returns services at the branch that can start at the specified time. (based on employees’ work schedules)
//   - branchID + employeeID + weekday: returns services offered by the employee at the branch on the given weekday. (based on employees’ work schedules)
//   - branchID + employeeID + weekday + timeStr: returns services offered by the employee at the branch on the given weekday that can start at the specified time. (based on employees’ work schedules)
//   // Employee //
//   - employeeID: returns services offered by the employee, regardless of branch, weekday and time start.
//   - employeeID + weekday: returns services offered by the employee on the given weekday regardless of branch. (based on employees’ work schedules)
//   - employeeID + timeStr: returns services the employee can start at the specified time, regardless of branch and weekday. (based on employees’ work schedules)
//   - employeeID + weekday + timeStr: returns services the employee can start at the specified time. (based on employees’ work schedules)
//   // Weekday //
//   - weekday: returns services available across the company on the given weekday. (based on employees’ work schedules)
//   - weekday + timeStr: returns services available across the company that can start on the given day and time. (based on employees’ work schedules)
//   // Time Start //
//   - timeStr: returns services available across the company that can start at the specified time, regardless of branch, employee and weekday. (based on employees’ work schedules)
// - If GET is "branches":
//   - All parameters nil: returns all branches of the company.
//   // Employee //
//   - employeeID: returns branches where the employee is assigned.
//   - employeeID + weekday: returns branches where the employee is available on the given weekday. (based on employees’ work schedules)
//   - employeeID + timeStr: returns branches where the employee is available to start at the specified time. (based on employees’ work schedules)
//   - employeeID + weekday + timeStr: returns branches where the employee is available to start at the specified day and time. (based on employees’ work schedules)
//   // Service //
//   - serviceID: returns branches that offer the given service.
//   - serviceID + weekday: returns branches that offer the service and have available employees on the given weekday. (based on employees’ work schedules)
//   - serviceID + timeStr: returns branches that offer the service and have available employees at the specified time. (based on employees’ work schedules)
//   - serviceID + weekday + timeStr: returns branches that offer the service with available employees at the specified day and time. (based on employees’ work schedules)
//   // Weekday //
//   - weekday: returns branches that have employees available on the given weekday. (based on employees’ work schedules)
//   - weekday + timeStr: returns branches that have employees available to start at the specified day and time. (based on employees’ work schedules)
//   // Time Start //
//   - timeStr: returns branches that have employees available to start at the specified time. (based on employees’ work schedules)
// - If GET is "employees":
//   - All parameters nil: returns all employees of the company.
//   // Branch //
//   - branchID: returns employees assigned to the specified branch.
//   - branchID + serviceID: returns employees at the branch who offer the given service.
//   - branchID + weekday: returns employees at the branch available on the given weekday. (based on employees’ work schedules)
//   - branchID + timeStr: returns employees at the branch available to start at the specified time in any day of the week. (based on employees’ work schedules)
//   - branchID + weekday + timeStr: returns employees at the branch available to start at the specified time and day. (based on employees’ work schedules)
//   // Service //
//   - serviceID: returns employees who offer the service, regardless of branch.
//   - serviceID + weekday: returns employees who offer the service and are available on the given weekday. (based on employees’ work schedules)
//   - serviceID + timeStr: returns employees who offer the service and are available to start at the specified time. (based on employees’ work schedules)
//   - serviceID + weekday + timeStr: returns employees who offer the service and are available to start at the specified time and day. (based on employees’ work schedules)
//   // Weekday //
//   - weekday: returns employees available on the given weekday. (based on employees’ work schedules)
//   - weekday + timeStr: returns employees available to start at the specified day and time. (based on employees’ work schedules)
//   // Time Start //
//   - timeStr: returns employees available to start at the specified time, regardless of branch or service. (based on employees’ work schedules)
// - If GET is "time_slots":
//   - All parameters nil: returns error.
//   // Branch //
//   - branchID: returns all time slots available at the branch. (based on employees’ work schedules)
//   - branchID + serviceID: returns time slots at the branch for the specified service. (based on employees’ work schedules)
//   - branchID + employeeID: returns time slots at the branch for the specified employee. (based on employees’ work schedules)
//   - branchID + weekday: returns time slots available at the branch on the given weekday. (based on employees’ work schedules)
//   - branchID + timeStr: returns time slots at the branch that can start at the specified time in any day of the week. (based on employees’ work schedules)
//   - branchID + weekday + timeStr: returns time slots at the branch that can start at the specified time and day. (based on employees’ work schedules)
//   - branchID + serviceID + weekday + timeStr: returns time slots at the branch for the service that can start at the specified time and day. (based on employees’ work schedules)
//   - branchID + employeeID + weekday + timeStr: returns time slots at the branch for the employee that can start at the specified time and day. (based on employees’ work schedules)
//   // Employee //
//   - employeeID: returns all available time slots for the employee. (based on employees’ work schedules)
//   - employeeID + weekday: returns time slots the employee is available on the given weekday. (based on employees’ work schedules)
//   - employeeID + timeStr: returns time slots the employee can start at the specified time in any day of the week. (based on employees’ work schedules)
//   - employeeID + weekday + timeStr: returns time slots the employee can start at the specified time and day. (based on employees’ work schedules)
//   // Service //
//   - serviceID: returns all time slots available for the service. (based on employees’ work schedules)
//   - serviceID + weekday: returns time slots available for the service on the given weekday. (based on employees’ work schedules)
//   - serviceID + timeStr: returns time slots available for the service at the specified time. (based on employees’ work schedules)
//   - serviceID + weekday + timeStr: returns time slots available for the service at the specified time and day. (based on employees’ work schedules)
//   // Weekday //
//   - weekday: returns all available time slots on the given weekday. (based on employees’ work schedules)
//   - weekday + timeStr: returns available time slots at the specified day and time. (based on employees’ work schedules)
//   // Time Start //
//   - timeStr: returns available time slots that can start at the specified time. (based on employees’ work schedules)
// The output dynamically adapts to the query context, providing only relevant scheduling options.
func (sf *schedule_filter) GetAvailableOptions() (any, error) {
	if sf.Error != nil {
		return nil, sf.Error
	}

	return nil, nil
}
