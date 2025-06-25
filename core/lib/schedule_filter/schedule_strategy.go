// filter/schedule_strategy.go
package ScheduleFilter

import (
	"fmt"
	"time"

	"agenda-kaki-go/core/config/db/model" // Ensure correct model path
	// "agenda-kaki-go/core/model" // Your provided path seems to be this
	"agenda-kaki-go/core/lib"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScheduleStrategy interface {
	Fetch(params *ScheduleQueryParams) (any, error)
}

// --- Services Strategy --- (Minor adjustment for clarity on ServiceID filter)
type ServicesStrategy struct {
	tx *gorm.DB
}

func (s *ServicesStrategy) Fetch(params *ScheduleQueryParams) (any, error) {
	if params.allExceptGetAreNil() {
		return nil, params.toLibError(fmt.Errorf("at least one filter is required when getting services"))
	}

	empBuilder := newEmployeeQueryBuilder(s.tx).availableOn(params)
	if params.BranchID != nil {
		empBuilder.inBranch(*params.BranchID)
	}
	if params.EmployeeID != nil {
		empBuilder.byID(*params.EmployeeID)
	}
	// Note: if params.ServiceID is present, it will filter employees who CAN provide this service,
	// which is implicitly handled by the subsequent service query.
	// If we want to find employees who provide a specific service, that's what providesService(id) does.
	// The current availableOn + inBranch + byID finds general employees.
	// Then we find services OF THESE employees.

	employeeIDs, err := empBuilder.getEmployeeIDs()
	if err != nil {
		return nil, err
	}
	if len(employeeIDs) == 0 {
		return []*model.Service{}, nil
	}

	serviceQuery := s.tx.Model(&model.Service{}).
		Joins("JOIN employee_services es ON es.service_id = services.id").
		Where("es.employee_id IN ?", employeeIDs)

	if params.BranchID != nil {
		serviceQuery = serviceQuery.Joins("JOIN branch_services bs ON bs.service_id = services.id").
			Where("bs.branch_id = ?", *params.BranchID)
	}

	// If a specific service_id was part of the input, filter the results to only this service.
	if params.ServiceID != nil {
		serviceQuery = serviceQuery.Where("services.id = ?", *params.ServiceID)
	}

	var services []*model.Service
	err = serviceQuery.Distinct().Find(&services).Error
	if err != nil {
		return nil, lib.Error.General.DatabaseError.WithError(err)
	}
	return services, nil
}

// --- Branches Strategy ---
type BranchesStrategy struct {
	tx *gorm.DB
}

func (s *BranchesStrategy) Fetch(params *ScheduleQueryParams) (any, error) {
	if params.allExceptGetAreNil() {
		var branches []*model.Branch
		err := s.tx.Find(&branches).Error
		if err != nil {
			return nil, lib.Error.General.DatabaseError.WithError(err)
		}
		return branches, nil
	}

	query := s.tx.Model(&model.Branch{})

	// If any filters that depend on employee characteristics (their ID, availability, or services they offer) are present,
	// we need to find those employees first, and then find the branches they are associated with.
	if params.EmployeeID != nil || params.OriginalStartTime != nil || params.ServiceID != nil {
		// Create a temporary params copy for employee search.
		// When GET=branches, params.BranchID is the target entity we're trying to find,
		// not a filter for the initial employee discovery (unless it's the *only* filter, handled later).
		// So, we nil out BranchID in empSearchParams to ensure availableOn and other filters
		// find employees globally or based on other non-branch criteria first.
		empSearchParams := *params
		if params.Get == GetBranches { // Ensure BranchID doesn't prematurely filter employees for branch discovery
			empSearchParams.BranchID = nil
		}

		// Use this builder to find employees based on non-branch specific criteria or global availability
		employeeDiscoverBuilder := newEmployeeQueryBuilder(s.tx).availableOn(&empSearchParams)

		if params.EmployeeID != nil {
			employeeDiscoverBuilder.byID(*params.EmployeeID)
		}
		if params.ServiceID != nil {
			// Find employees who provide the service, irrespective of branch at this stage.
			employeeDiscoverBuilder.providesService(*params.ServiceID)
		}

		employeeIDs, err := employeeDiscoverBuilder.getEmployeeIDs()
		if err != nil {
			return nil, err // Propagate DB errors
		}

		// If employee-specific filters were applied and no employees were found,
		// then no branches can be found via this path.
		if len(employeeIDs) == 0 && (params.EmployeeID != nil || params.OriginalStartTime != nil || params.ServiceID != nil) {
			return []*model.Branch{}, nil
		}

		// If we found employees, filter branches to those where these employees work.
		if len(employeeIDs) > 0 {
			query = query.Joins("JOIN employee_branches eb ON eb.branch_id = branches.id").
				Where("eb.employee_id IN ?", employeeIDs)
		}
	}

	// After potentially filtering by employees, apply direct branch filters.

	// If a specific BranchID is requested (e.g., /path?get=branches&branch_id=X),
	// this acts as the primary or an additional AND filter for the branches.
	if params.BranchID != nil {
		query = query.Where("branches.id = ?", *params.BranchID)
	}

	// Additionally, if a serviceID is given, the resulting branches must offer this service.
	// This ensures branches are relevant to the service, either through employee association or direct offering.
	if params.ServiceID != nil {
		// Using Distinct elsewhere, but multiple Joins on same tables can cause issues
		// if not careful. Assuming GORM handles it or aliases are used if complex.
		// Add a check to prevent duplicate join if already joined for service through employee path
		// For simplicity, current GORM might handle this, or needs an alias for 'bs'.
		query = query.Joins("JOIN branch_services bs ON bs.branch_id = branches.id").
			Where("bs.service_id = ?", *params.ServiceID)
	}

	var branches []*model.Branch
	err := query.Distinct().Find(&branches).Error // Distinct is important due to multiple joins
	if err != nil {
		return nil, lib.Error.General.DatabaseError.WithError(err)
	}
	return branches, nil
}

// --- Employees Strategy ---
type EmployeesStrategy struct {
	tx *gorm.DB
}

func (s *EmployeesStrategy) Fetch(params *ScheduleQueryParams) (any, error) {
	if params.allExceptGetAreNil() {
		var employees []*model.Employee
		// Preload standard associations for a general employee list
		err := s.tx.Preload("Branches").Preload("Services").Preload("EmployeeWorkSchedule").Find(&employees).Error
		if err != nil {
			return nil, lib.Error.General.DatabaseError.WithError(err)
		}
		return employees, nil
	}

	builder := newEmployeeQueryBuilder(s.tx).availableOn(params)
	if params.BranchID != nil {
		builder.inBranch(*params.BranchID)
	}
	if params.ServiceID != nil {
		builder.providesService(*params.ServiceID)
	}
	if params.EmployeeID != nil { // Allow direct query for a specific employee's details
		builder.byID(*params.EmployeeID)
	}

	// .getEmployees() now preloads EmployeeWorkSchedule
	employees, err := builder.getEmployees()
	if err != nil {
		return nil, err
	}
	return employees, nil
}

// --- TimeSlots Strategy ---
type TimeSlot struct {
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	BranchID   uuid.UUID `json:"branch_id"`
	EmployeeID uuid.UUID `json:"employee_id"`
}

type TimeSlotsStrategy struct {
	tx *gorm.DB
}

func getTodayUTC() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

func (s *TimeSlotsStrategy) determineOutputLocation(params *ScheduleQueryParams) (*time.Location, error) {
	if params.OriginalStartTime != nil {
		return params.OriginalStartTime.Location(), nil
	}
	if params.BranchID != nil {
		var branch model.Branch
		if err := s.tx.First(&branch, *params.BranchID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, params.toLibError(fmt.Errorf("branch with ID %s not found for determining timezone", *params.BranchID))
			}
			return nil, lib.Error.General.DatabaseError.WithError(err)
		}
		// Your model.Branch has TimeZone as time.Location, not *time.Location
		// Ensure it's properly initialized. If it can be zero, handle that.
		if branch.TimeZone != "" && branch.TimeZone != "UTC" { // A bit of a hack to check if it's "set"
			loc, err := time.LoadLocation(branch.TimeZone)
			if err != nil {
				return nil, params.toLibError(fmt.Errorf("invalid branch timezone %s: %w", branch.TimeZone, err))
			}
			return loc, nil
		}
	}
	return time.UTC, nil // Default to UTC
}

func (s *TimeSlotsStrategy) Fetch(params *ScheduleQueryParams) (any, error) {
	if params.allExceptGetAreNil() {
		return nil, params.toLibError(fmt.Errorf("at least one filter (branch_id, employee_id, service_id, or start_time) is required when getting time_slots"))
	}

	outputLocation, err := s.determineOutputLocation(params)
	if err != nil {
		return nil, err // Error already wrapped
	}

	var queryStartDate, queryEndDate time.Time
	var specificTimeOfDayFilter *time.Time // This is already 0000-01-01 HH:MM:SS UTC

	if params.OriginalStartTime == nil {
		queryStartDate = getTodayUTC()
		queryEndDate = queryStartDate.AddDate(0, 0, 6)
	} else if params.IsSpecificDateQuery {
		// params.QueryDate is YYYY-MM-DD 00:00:00 in OriginalStartTime's Location.
		// For range calculations, we convert to UTC.
		queryStartDate = params.QueryDate.In(time.UTC)
		queryEndDate = queryStartDate // Slots for this specific day only
		if params.QueryTimeOfDay != nil {
			specificTimeOfDayFilter = params.QueryTimeOfDay // This is 0000-01-01 HH:MM:SS UTC
		}
	} else if params.IsTimeOfDayQuery {
		queryStartDate = getTodayUTC()
		queryEndDate = queryStartDate.AddDate(0, 0, 29)
		specificTimeOfDayFilter = params.QueryTimeOfDay // This is 0000-01-01 HH:MM:SS UTC
	} else {
		return nil, params.toLibError(fmt.Errorf("invalid start_time state for time_slots"))
	}

	// 2. Find potentially available employees
	builder := newEmployeeQueryBuilder(s.tx)
	if params.BranchID != nil {
		builder.inBranch(*params.BranchID)
	}
	if params.EmployeeID != nil {
		builder.byID(*params.EmployeeID)
	}
	if params.ServiceID != nil {
		builder.providesService(*params.ServiceID)
	}

	employees, err := builder.getEmployees() // getEmployees preloads EmployeeWorkSchedule
	if err != nil {
		return nil, err // Error already wrapped by builder or database error
	}
	if len(employees) == 0 {
		return []TimeSlot{}, nil
	}
	employeeIDs := make([]uuid.UUID, len(employees))
	for i, e := range employees {
		employeeIDs[i] = e.ID
	}

	// 3. Fetch existing appointments
	var existingAppointments []model.Appointment
	// Query appointments by their StartTime.
	// IMPORTANT: Appointment.StartTime is stored with its own Appointment.TimeZone.
	// To compare with our UTC queryStartDate/queryEndDate, we need to be careful.
	// It's often best to query StartTime in UTC if possible, or adjust the query.
	// For now, let's assume we can query against the raw StartTime and the DB handles it,
	// but the `blockedSlots` key generation MUST normalize to UTC.
	err = s.tx.Model(&model.Appointment{}).
		Where("employee_id IN ?", employeeIDs).
		Where("is_cancelled = ?", false).
		// The following time condition needs to reliably compare against appointments
		// which might be stored in their local timezones.
		// Option A: Convert queryStartDate/EndDate to a range that covers all timezones, then filter in Go (less efficient).
		// Option B: If DB supports timezone conversion on the fly: WHERE appt.start_time AT TIME ZONE 'UTC' >= queryStartDateUTC
		// Option C (Simpler, assuming StartTime field can be compared across timezones by DB, might not be fully accurate without AT TIME ZONE):
		Where("start_time >= ?", queryStartDate).               // queryStartDate is UTC midnight
		Where("start_time < ?", queryEndDate.AddDate(0, 0, 1)). // queryEndDate is UTC midnight of next day
		Find(&existingAppointments).Error
	if err != nil {
		return nil, lib.Error.General.DatabaseError.WithError(err)
	}

	blockedSlots := make(map[string]bool) // Key: employeeID_RFC3339_UTC
	for _, appt := range existingAppointments {
		apptStartTimeUTC := appt.StartTime.In(time.UTC)
		blockedKey := fmt.Sprintf("%s_%s", appt.EmployeeID.String(), apptStartTimeUTC.Format(time.RFC3339))
		blockedSlots[blockedKey] = true
	}

	// 4. Generate slots
	var allSlots []TimeSlot
	currentDayUTC := queryStartDate // Iterating days in UTC
	for !currentDayUTC.After(queryEndDate) {
		currentDayWeekday := currentDayUTC.Weekday()

		for _, emp := range employees {
			if emp.EmployeeWorkSchedule == nil {
				continue
			}

			slotDurationMinutes := emp.SlotTimeDiff
			if slotDurationMinutes == 0 {
				slotDurationMinutes = 30
			}
			slotDuration := time.Duration(slotDurationMinutes) * time.Minute

			for _, wr := range emp.EmployeeWorkSchedule {
				if wr.Weekday != currentDayWeekday {
					continue
				}
				if params.BranchID != nil && wr.BranchID != *params.BranchID {
					continue
				}

				// wr.StartTime and wr.EndTime are 0000-01-01 HH:MM:SS UTC (from EmployeeWorkRange model's hook)
				// Construct shift start/end for currentDayUTC in UTC.
				shiftStartTimeUTC := time.Date(currentDayUTC.Year(), currentDayUTC.Month(), currentDayUTC.Day(),
					wr.StartTime.Hour(), wr.StartTime.Minute(), wr.StartTime.Second(), 0, time.UTC)

				shiftEndTimeUTC := time.Date(currentDayUTC.Year(), currentDayUTC.Month(), currentDayUTC.Day(),
					wr.EndTime.Hour(), wr.EndTime.Minute(), wr.EndTime.Second(), 0, time.UTC)

				if !shiftEndTimeUTC.After(shiftStartTimeUTC) {
					continue
				}

				for slotCandidateStartUTC := shiftStartTimeUTC; slotCandidateStartUTC.Add(slotDuration).Compare(shiftEndTimeUTC) <= 0; slotCandidateStartUTC = slotCandidateStartUTC.Add(slotDuration) {
					if specificTimeOfDayFilter != nil {
						if slotCandidateStartUTC.Hour() != specificTimeOfDayFilter.Hour() ||
							slotCandidateStartUTC.Minute() != specificTimeOfDayFilter.Minute() {
							continue
						}
					}

					// Check against blocked slots (keys are UTC)
					blockedKey := fmt.Sprintf("%s_%s", emp.ID.String(), slotCandidateStartUTC.Format(time.RFC3339))
					if _, isBlocked := blockedSlots[blockedKey]; isBlocked {
						continue
					}

					// Convert the generated UTC slot times to the determined outputLocation
					slotStartInOutputLoc := slotCandidateStartUTC.In(outputLocation)
					slotEndInOutputLoc := slotCandidateStartUTC.Add(slotDuration).In(outputLocation)

					allSlots = append(allSlots, TimeSlot{
						StartTime:  slotStartInOutputLoc,
						EndTime:    slotEndInOutputLoc,
						BranchID:   wr.BranchID,
						EmployeeID: emp.ID,
					})
				}
			}
		}
		currentDayUTC = currentDayUTC.AddDate(0, 0, 1)
	}
	return allSlots, nil
}
