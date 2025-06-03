package utilsT

import (
	"agenda-kaki-go/core/config/db/model"
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/config/namespace"
	handlerT "agenda-kaki-go/core/test/handlers"
	modelT "agenda-kaki-go/core/test/models"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type FoundAppointmentSlot struct {
	StartTimeRFC3339 string
	BranchID         string
	ServiceID        string
}

const (
	slotSearchHorizonDays = 14               // Example: search up to 2 weeks ahead
	slotSearchTimeStep    = 15 * time.Minute // Example: check every 15 minutes
)

// Searches for a valid appointment slot for an employee within the company's work schedule
// Returns a FoundAppointmentSlot if found, or an error if no valid slot is available
func FindValidAppointmentSlot(employee *modelT.Employee, company *modelT.Company, preferredLocation *time.Location) (slot *FoundAppointmentSlot, found bool, err error) {
	if preferredLocation == nil {
		return nil, false, fmt.Errorf("preferredLocation is nil; timezone must be explicitly passed")
	}
	fmt.Printf("---- Starting findValidAppointmentSlot for Employee ID: %s ----\n", employee.Created.ID)
	fmt.Printf("PreferredLocation: %s, SlotSearchHorizon: %d days, TimeStep: %v\n", preferredLocation.String(), slotSearchHorizonDays, slotSearchTimeStep)

	schedule := employee.Created.WorkSchedule
	weekdaySchedules := map[time.Weekday][]mJSON.WorkRange{
		time.Sunday:    schedule.Sunday,
		time.Monday:    schedule.Monday,
		time.Tuesday:   schedule.Tuesday,
		time.Wednesday: schedule.Wednesday,
		time.Thursday:  schedule.Thursday,
		time.Friday:    schedule.Friday,
		time.Saturday:  schedule.Saturday,
	}

	nowInPreferredLocation := time.Now().In(preferredLocation)
	searchStartDate := time.Date(nowInPreferredLocation.Year(), nowInPreferredLocation.Month(), nowInPreferredLocation.Day(), 0, 0, 0, 0, preferredLocation)
	fmt.Printf("Searching from %s (now in preferred TZ is %s)\n", searchStartDate.Format(time.RFC3339), nowInPreferredLocation.Format(time.RFC3339))

	for dayOffset := range slotSearchHorizonDays {
		currentSearchDate := searchStartDate.AddDate(0, 0, dayOffset)
		currentWeekday := currentSearchDate.Weekday()
		fmt.Printf("  Checking Date: %s (Weekday: %s, DayOffset: %d)\n", currentSearchDate.Format("2006-01-02"), currentWeekday, dayOffset)

		dayScheduleRanges, hasScheduleForDay := weekdaySchedules[currentWeekday]
		if !hasScheduleForDay || len(dayScheduleRanges) == 0 {
			fmt.Printf("    No schedule ranges for this day.\n")
			continue
		}
		fmt.Printf("    Found %d schedule range(s) for %s\n", len(dayScheduleRanges), currentWeekday)

		for wrIdx := range dayScheduleRanges {
			wr := &dayScheduleRanges[wrIdx]
			fmt.Printf("      Processing WorkRange #%d: Start='%s', End='%s', BranchID='%s'\n", wrIdx, wr.Start, wr.End, wr.BranchID)

			if wr.Start == "" || wr.End == "" || wr.BranchID == uuid.Nil {
				fmt.Printf("        Skipping WR: Invalid data (Start, End, or BranchID missing).\n")
				continue
			}

			var targetBranch *modelT.Branch
			branchObjFound := false
			for k := range company.Branches {
				if company.Branches[k].Created.ID == wr.BranchID {
					targetBranch = company.Branches[k]
					branchObjFound = true
					break
				}
			}
			if !branchObjFound {
				fmt.Printf("        Skipping WR: Branch object for ID %s not found in company.Branches.\n", wr.BranchID)
				continue
			}
			if !isEmployeeAssignedToBranch(employee, wr.BranchID) {
				fmt.Printf("        Skipping WR: Employee %s not assigned to Branch %s (ID: %s).\n", employee.Created.ID, targetBranch.Created.Name, wr.BranchID)
				continue
			}
			fmt.Printf("        Branch '%s' (ID: %s) is valid and employee is assigned.\n", targetBranch.Created.Name, wr.BranchID)

			employeeServiceIDs := make(map[uuid.UUID]bool)
			if len(employee.Services) == 0 {
				fmt.Printf("        Skipping WR: Employee has no services assigned.\n")
				continue
			}
			for _, service := range employee.Services {
				employeeServiceIDs[service.Created.ID] = true
			}
			fmt.Printf("        Employee services: %v\n", employeeServiceIDs)

			if len(targetBranch.Services) == 0 {
				fmt.Printf("        Skipping WR: Branch has no services assigned.\n")
				continue
			}
			branchServiceIDs := make(map[uuid.UUID]bool)
			for _, service := range targetBranch.Services {
				branchServiceIDs[service.Created.ID] = true
			}
			fmt.Printf("        Branch services: %v\n", branchServiceIDs)

			validServiceIDs := []uuid.UUID{}
			for empSvcID := range employeeServiceIDs {
				if branchServiceIDs[empSvcID] {
					validServiceIDs = append(validServiceIDs, empSvcID)
				}
			}
			if len(validServiceIDs) == 0 {
				fmt.Printf("        Skipping WR: No common services between employee and branch.\n")
				continue
			}
			fmt.Printf("        Found %d common/valid services: %v\n", len(validServiceIDs), validServiceIDs)

			workRangeStartDateTime, err := parseTimeWithLocation(currentSearchDate, wr.Start, preferredLocation)
			if err != nil {
				fmt.Printf("        Skipping WR: Error parsing WorkRange Start Time '%s': %v\n", wr.Start, err)
				continue
			}
			workRangeEndDateTime, err := parseTimeWithLocation(currentSearchDate, wr.End, preferredLocation)
			if err != nil {
				fmt.Printf("        Skipping WR: Error parsing WorkRange End Time '%s': %v\n", wr.End, err)
				continue
			}
			if !workRangeStartDateTime.Before(workRangeEndDateTime) {
				fmt.Printf("        Skipping WR: WorkRange Start (%s) is not before End (%s).\n", workRangeStartDateTime, workRangeEndDateTime)
				continue
			}
			fmt.Printf("        Parsed WorkRange Times: Start=%s, End=%s\n", workRangeStartDateTime.Format(time.RFC3339), workRangeEndDateTime.Format(time.RFC3339))

			for _, selectedServiceID := range validServiceIDs {
				fmt.Printf("          Trying Service ID: %s\n", selectedServiceID)
				var selectedServiceDurationMinutes uint
				serviceFound := false
				for _, s := range company.Services { // Search in company's master list of services
					if s.Created.ID == selectedServiceID {
						selectedServiceDurationMinutes = s.Created.Duration
						serviceFound = true
						break
					}
				}
				if !serviceFound {
					return nil, false, fmt.Errorf("service with ID %s not found in company's services", selectedServiceID)
				}
				if selectedServiceDurationMinutes == 0 {
					fmt.Printf("          Skipping Service: Duration is 0 minutes.\n")
					continue
				}
				fmt.Printf("          Service Duration: %d minutes\n", selectedServiceDurationMinutes)

				serviceDuration := time.Duration(selectedServiceDurationMinutes) * time.Minute

				for potentialStartTime := workRangeStartDateTime; potentialStartTime.Before(workRangeEndDateTime); potentialStartTime = potentialStartTime.Add(slotSearchTimeStep) {
					potentialEndTime := potentialStartTime.Add(serviceDuration)
					fmt.Printf("            Testing Slot: PotentialStart=%s, PotentialEnd=%s\n", potentialStartTime.Format(time.RFC3339), potentialEndTime.Format(time.RFC3339))

					// If the service cannot be completed by the end of the work range
					if potentialEndTime.After(workRangeEndDateTime) {
						fmt.Printf("              Slot ends after work range. Breaking from time-stepping for this service in this WR.\n")
						// No further time steps for THIS service in THIS work range will fit.
						break // Break from the time-stepping loop (for current service)
					}

					// Skip past slots
					if potentialStartTime.Before(nowInPreferredLocation) {
						fmt.Printf("              Skipping slot: Potential start time is in the past.\n")
						continue
					}

					overlap := false
					fmt.Printf("            Checking for overlaps with %d existing appointments...\n", len(employee.Created.Appointments))
					for _, appt := range employee.Created.Appointments {
						var apptDur uint
						apptServiceFound := false
						if appt.Service != nil && appt.Service.ID != uuid.Nil { // Ensure Service object and its ID are valid
							apptDur = appt.Service.Duration
							apptServiceFound = true
						} else { // Fallback to ServiceID
							for _, s := range company.Services {
								if s.Created.ID == appt.ServiceID {
									apptDur = s.Created.Duration
									apptServiceFound = true
									break
								}
							}
						}
						if !apptServiceFound {
							return nil, false, fmt.Errorf("appointment with ID %s has no valid service duration", appt.ID)
						} else if apptDur == 0 {
							return nil, false, fmt.Errorf("appointment with ID %s has a service duration of 0 minutes", appt.ID)
						}

						existingStart := appt.StartTime.In(preferredLocation)
						existingEnd := existingStart.Add(time.Duration(apptDur) * time.Minute)
						fmt.Printf("              Comparing with existing appt [%s to %s]\n", existingStart.Format(time.RFC3339), existingEnd.Format(time.RFC3339))

						if potentialStartTime.Before(existingEnd) && potentialEndTime.After(existingStart) {
							fmt.Printf("                Overlap detected with existing appointment!\n")
							overlap = true
							break
						}
					}

					if !overlap {
						fmt.Printf("            SUCCESS! Found available slot: StartTime=%s, BranchID=%s, ServiceID=%s\n",
							potentialStartTime.Format(time.RFC3339), wr.BranchID.String(), selectedServiceID.String())
						slot = &FoundAppointmentSlot{
							StartTimeRFC3339: potentialStartTime.Format(time.RFC3339),
							BranchID:         wr.BranchID.String(),
							ServiceID:        selectedServiceID.String(),
						}
						return slot, true, nil // Found a slot
					} else {
						fmt.Printf("            Slot has overlap, trying next potential start time.\n")
					}
				} // End time-stepping loop
				fmt.Printf("          Finished time-stepping for service %s in this WR.\n", selectedServiceID)
			} // End service loop
		} // End work-range loop
	} // End day-offset loop

	return slot, false, fmt.Errorf("no valid appointment slot found for employee %s in company %s within the search horizon", employee.Created.ID, company.Created.ID)
}

func FindValidAppointmentSlotV2(employee *modelT.Employee, preferredLocation *time.Location) (slot *FoundAppointmentSlot, found bool, err error) {
	if preferredLocation == nil {
		return nil, false, fmt.Errorf("preferredLocation is nil; timezone must be explicitly passed")
	}
	fmt.Printf("---- Starting findValidAppointmentSlot for Employee ID: %s ----\n", employee.Created.ID.String())
	fmt.Printf("PreferredLocation: %s, SlotSearchHorizon: %d days, TimeStep: %v\n", preferredLocation.String(), slotSearchHorizonDays, slotSearchTimeStep)

	WorkSchedule := employee.Created.WorkSchedule
	weekdaySchedules := map[time.Weekday][]mJSON.WorkRange{
		time.Sunday:    WorkSchedule.Sunday,
		time.Monday:    WorkSchedule.Monday,
		time.Tuesday:   WorkSchedule.Tuesday,
		time.Wednesday: WorkSchedule.Wednesday,
		time.Thursday:  WorkSchedule.Thursday,
		time.Friday:    WorkSchedule.Friday,
		time.Saturday:  WorkSchedule.Saturday,
	}

	nowInPreferredLocation := time.Now().In(preferredLocation)
	searchStartDate := time.Date(nowInPreferredLocation.Year(), nowInPreferredLocation.Month(), nowInPreferredLocation.Day(), 0, 0, 0, 0, preferredLocation)
	fmt.Printf("Searching from %s (now in preferred TZ is %s)\n", searchStartDate.Format(time.RFC3339), nowInPreferredLocation.Format(time.RFC3339))

	isBranchAssignedToEmployee := func(branch_id string) (bool, error) {
		var Branch model.Branch
		if err := handlerT.NewHttpClient().
			Method("GET").
			URL(fmt.Sprintf("/branch/%s", branch_id)).
			ExpectedStatus(200).
			Header(namespace.HeadersKey.Company, employee.Company.Created.ID.String()).
			Header(namespace.HeadersKey.Auth, employee.Company.Owner.X_Auth_Token).
			Send(nil).
			ParseResponse(&Branch).Error; err != nil {
			return false, fmt.Errorf("failed to get branch by id: %w", err)
		}
		isAssigned := false
		for _, b := range employee.Created.Branches {
			if b.ID.String() == branch_id {
				isAssigned = true
				break
			}
		}
		return isAssigned, nil
	}

	isServiceAssignedToEmployeeV2 := func(service_id string) (bool, error) {
		var Service model.Service
		if err := handlerT.NewHttpClient().
			Method("GET").
			URL(fmt.Sprintf("/service/%s", service_id)).
			ExpectedStatus(200).
			Header(namespace.HeadersKey.Company, employee.Company.Created.ID.String()).
			Header(namespace.HeadersKey.Auth, employee.Company.Owner.X_Auth_Token).
			Send(nil).
			ParseResponse(&Service).Error; err != nil {
			return false, fmt.Errorf("failed to get service by id: %w", err)
		}
		isAssigned := false
		for _, s := range employee.Created.Services {
			if s.ID.String() == service_id {
				isAssigned = true
				break
			}
		}
		return isAssigned, nil
	}

	isServiceAssignedToBranchV2 := func(branch_id, service_id string) (bool, error) {
		var Branch model.Branch
		if err := handlerT.NewHttpClient().
			Method("GET").
			URL(fmt.Sprintf("/branch/%s", branch_id)).
			ExpectedStatus(200).
			Header(namespace.HeadersKey.Company, employee.Company.Created.ID.String()).
			Header(namespace.HeadersKey.Auth, employee.Company.Owner.X_Auth_Token).
			Send(nil).
			ParseResponse(&Branch).Error; err != nil {
			return false, fmt.Errorf("failed to get branch by id: %w", err)
		}
		isAssigned := false
		for _, s := range Branch.Services {
			if s.ID.String() == service_id {
				isAssigned = true
				break
			}
		}
		return isAssigned, nil
	}

	isValidWorkRangeDate := func(wr mJSON.WorkRange, currentSearchDate time.Time) (bool, error) {
		workRangeStartDateTime, err := parseTimeWithLocation(currentSearchDate, wr.Start, preferredLocation)
		if err != nil {
			return false, fmt.Errorf("error parsing WorkRange Start Time '%s': %w", wr.Start, err)
		}
		workRangeEndDateTime, err := parseTimeWithLocation(currentSearchDate, wr.End, preferredLocation)
		if err != nil {
			return false, fmt.Errorf("error parsing WorkRange End Time '%s': %w", wr.End, err)
		}
		if !workRangeStartDateTime.Before(workRangeEndDateTime) {
			return false, fmt.Errorf("WorkRange Start (%s) is not before End (%s)", workRangeStartDateTime, workRangeEndDateTime)
		}
		return true, nil
	}

	for dayOffset := range slotSearchHorizonDays {
		currentSearchDate := searchStartDate.AddDate(0, 0, dayOffset)
		currentWeekday := currentSearchDate.Weekday()
		fmt.Printf("  Checking Date: %s (Weekday: %s, DayOffset: %d)\n", currentSearchDate.Format("2006-01-02"), currentWeekday, dayOffset)

		dayScheduleRanges, hasScheduleForDay := weekdaySchedules[currentWeekday]
		if !hasScheduleForDay || len(dayScheduleRanges) == 0 {
			fmt.Printf("    No schedule ranges for this day.\n")
			continue
		}
		fmt.Printf("    Found %d schedule range(s) for %s\n", len(dayScheduleRanges), currentWeekday)
		for iWr, wr := range dayScheduleRanges {
			fmt.Printf("      Processing WorkRange #%d: Start='%s', End='%s', BranchID='%s'\n", iWr, wr.Start, wr.End, wr.BranchID)
			if wr.Start == "" || wr.End == "" || wr.BranchID == uuid.Nil {
				return nil, false, fmt.Errorf("work range #%d has invalid data for employee %s: Start, End, or BranchID is missing. WorkSchedule: %+v", iWr, employee.Created.ID.String(), WorkSchedule)
			}
			if ok, err := isBranchAssignedToEmployee(wr.BranchID.String()); !ok {
				if err != nil {
					return nil, false, fmt.Errorf("error checking branch assignment for employee %s: %w", employee.Created.ID.String(), err)
				}
				return nil, false, fmt.Errorf("work range #%d has invalid branch assignment as employee %s is not assigned to branch %s", iWr, employee.Created.ID.String(), wr.BranchID.String())
			}
			if ok, err := isValidWorkRangeDate(wr, currentSearchDate); !ok {
				return nil, false, fmt.Errorf("work range #%d has invalid work range assignment for employee %s: %w", iWr, employee.Created.ID.String(), err)
			}
			for iSrv, service_id := range wr.Services {
				if uuid.Nil == service_id {
					return nil, false, fmt.Errorf("employee %s work range #%d has invalid service assignment at #%d: service ID is nil", employee.Created.ID.String(), iWr, iSrv)
				}
				if ok, err := isServiceAssignedToEmployeeV2(service_id.String()); !ok {
					if err != nil {
						return nil, false, fmt.Errorf("error checking service assignment for employee %s: %w", employee.Created.ID.String(), err)
					}
					return nil, false, fmt.Errorf("employee %s work range #%d has invalid service assignment at #%d: employee is not assigned to service %s", employee.Created.ID.String(), iWr, iSrv, service_id.String())
				}
				if ok, err := isServiceAssignedToBranchV2(wr.BranchID.String(), service_id.String()); !ok {
					if err != nil {
						return nil, false, fmt.Errorf("error checking service assignment for branch %s: %w", wr.BranchID.String(), err)
					}
					return nil, false, fmt.Errorf("employee %s work range #%d has invalid service assignment at #%d: branch %s is not assigned to service %s", employee.Created.ID.String(), iWr, iSrv, wr.BranchID.String(), service_id.String())
				}
				fmt.Printf("        WorkRange #%d, Service #%d: Valid service %s found for employee %s in branch %s\n", iWr, iSrv, service_id.String(), employee.Created.ID.String(), wr.BranchID.String())
				serviceDuration := time.Duration(employee.Services[iSrv].Created.Duration) * time.Minute
				for potentialStartTime := time.Date(currentSearchDate.Year(), currentSearchDate.Month(), currentSearchDate.Day(), 0, 0, 0, 0, preferredLocation); potentialStartTime.Before(currentSearchDate.AddDate(0, 0, slotSearchHorizonDays)); potentialStartTime = potentialStartTime.Add(slotSearchTimeStep) {
					potentialEndTime := potentialStartTime.Add(serviceDuration)
					fmt.Printf("          Testing Slot: PotentialStart=%s, PotentialEnd=%s\n", potentialStartTime.Format(time.RFC3339), potentialEndTime.Format(time.RFC3339))
					// If the service cannot be completed by the end of the work range
					if potentialEndTime.After(currentSearchDate.AddDate(0, 0, slotSearchHorizonDays)) {
						fmt.Printf("            Slot ends after work range. Breaking from time-stepping for this service in this work range.\n")
						// No further time steps for THIS service in THIS work range will fit.
						break // Break from the time-stepping loop (for current service)
					}
					// Skip past slots
					if potentialStartTime.Before(nowInPreferredLocation) {
						fmt.Printf("            Slot starts before current time. Skipping...\n")
						continue
					}
					overlap := false
					fmt.Printf("          Checking for overlaps with %d existing appointments...\n", len(employee.Created.Appointments))
					for _, appt := range employee.Created.Appointments {
						if appt.StartTime.Before(potentialEndTime) && appt.EndTime.After(potentialStartTime) {
							overlap = true
							fmt.Printf("            Found overlapping appointment: %s - %s\n", appt.StartTime.Format(time.RFC3339), appt.EndTime.Format(time.RFC3339))
						}
					}
					if overlap {
						fmt.Printf("          Slot %s - %s overlaps with existing appointments. Skipping...\n", potentialStartTime.Format(time.RFC3339), potentialEndTime.Format(time.RFC3339))
						continue
					}
					fmt.Printf("          Slot %s - %s is valid and has no overlaps.\n", potentialStartTime.Format(time.RFC3339), potentialEndTime.Format(time.RFC3339))
					slot = &FoundAppointmentSlot{
						StartTimeRFC3339: potentialStartTime.Format(time.RFC3339),
						BranchID:         wr.BranchID.String(),
						ServiceID:        service_id.String(),
					}
					return slot, true, nil // Found a valid slot
				}
			}
		}
	}
	return nil, false, fmt.Errorf("no valid appointment slot found for employee %s in company %s within the search horizon", employee.Created.ID.String(), employee.Company.Created.ID.String())
}

// Helper to parse HH:MM or HH:MM:SS time string into a full time.Time on a specific date/location
func parseTimeWithLocation(targetDate time.Time, timeStr string, loc *time.Location) (time.Time, error) {
	layout := "15:04" // Default HH:MM
	colonCount := 0
	for _, r := range timeStr {
		if r == ':' {
			colonCount++
		}
	}
	if colonCount == 2 { // Detect HH:MM:SS
		layout = "15:04:05"
	}

	parsedTime, err := time.ParseInLocation(layout, timeStr, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time string '%s' with layout '%s': %w", timeStr, layout, err)
	}
	// Combine the date part from targetDate with the time parts from parsedTime
	return time.Date(
		targetDate.Year(), targetDate.Month(), targetDate.Day(),
		parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0, // Nanoseconds set to 0
		loc,
	), nil
}

// // findNextAvailableSlot attempts to find the next available start time for an employee
// // based on their work schedule, starting from the one provided.
// // NOTE: This is a simplified helper for testing and might not cover all edge cases.
// // Renamed and modified to handle full timestamps
// func findNextAvailableSlotRFC3339(employee *modelT.Employee, currentStartTimeRFC3339 string) (string, error) {
// 	layoutRFC3339 := time.RFC3339
// 	start, err := time.Parse(layoutRFC3339, currentStartTimeRFC3339)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to parse current start time RFC3339 '%s': %v", currentStartTimeRFC3339, err)
// 	}

// 	// --- Service Duration calculation (remains similar) ---
// 	if len(employee.Services) == 0 {
// 		return "", fmt.Errorf("employee has no services, cannot determine next slot based on duration.")
// 	}
// 	duration := time.Duration(employee.Services[0].Created.Duration) * time.Minute
// 	nextPossibleStart := start.Add(duration)

// 	// --- Schedule Check (Now uses the date from input timestamp) ---
// 	schedule := employee.Created.WorkSchedule
// 	targetDate := start     // Use the date from the input timestamp
// 	loc := start.Location() // Preserve the timezone/location

// 	// Get the day schedule based on the start time's day of the week
// 	dayOfWeek := targetDate.Weekday()
// 	var daySchedule []mJSON.WorkRange
// 	switch dayOfWeek {
// 	case time.Monday:
// 		daySchedule = schedule.Monday
// 	case time.Tuesday:
// 		daySchedule = schedule.Tuesday
// 	case time.Wednesday:
// 		daySchedule = schedule.Wednesday
// 	case time.Thursday:
// 		daySchedule = schedule.Thursday
// 	case time.Friday:
// 		daySchedule = schedule.Friday
// 	case time.Saturday:
// 		daySchedule = schedule.Saturday
// 	case time.Sunday:
// 		daySchedule = schedule.Sunday
// 	}

// 	if len(daySchedule) == 0 {
// 		return "", fmt.Errorf("No work schedule found for employee %s on %s.", employee.Created.ID, dayOfWeek)
// 	}

// 	timeLayout := "15:04" // Layout for parsing schedule times HH:MM

// 	for _, block := range daySchedule {
// 		blockStartTimeStr := block.Start
// 		blockEndTimeStr := block.End

// 		// Parse the block start/end times *relative to the targetDate's date and location*
// 		blockStartParsed, err := time.ParseInLocation(timeLayout, blockStartTimeStr, loc)
// 		if err != nil {
// 			return "", fmt.Errorf("Warn: bad block start time %s: %v", blockStartTimeStr, err)
// 		}
// 		blockEndParsed, err := time.ParseInLocation(timeLayout, blockEndTimeStr, loc)
// 		if err != nil {
// 			return "", fmt.Errorf("Warn: bad block end time %s: %v", blockEndTimeStr, err)
// 		}

// 		// Construct full datetime objects for the block boundaries on the target date
// 		blockStartDateTime := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(),
// 			blockStartParsed.Hour(), blockStartParsed.Minute(), 0, 0, loc)
// 		blockEndDateTime := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(),
// 			blockEndParsed.Hour(), blockEndParsed.Minute(), 0, 0, loc)

// 		// Check if the calculated nextPossibleStart fits within this block
// 		if !nextPossibleStart.Before(blockStartDateTime) && !nextPossibleStart.Add(duration).After(blockEndDateTime) {
// 			return nextPossibleStart.Format(layoutRFC3339), nil // Return the RFC3339 string
// 		}
// 	}

// 	return "", fmt.Errorf("Could not find next available slot after %s for employee %s.", currentStartTimeRFC3339, employee.Created.ID)
// }

func isEmployeeAssignedToBranch(e *modelT.Employee, branchID uuid.UUID) bool {
	for _, b := range e.Branches {
		if b.Created.ID == branchID {
			return true
		}
	}
	return false
}

func RescheduleAppointmentRandomly(s int, employee *modelT.Employee, company *modelT.Company, appointment_id, token string) error {
	preferredLocation := time.UTC // Choose your timezone (e.g., UTC)
	appointmentSlot, found, err := FindValidAppointmentSlot(employee, company, preferredLocation)
	if err != nil {
		return fmt.Errorf("failed to find valid appointment slot: %w", err)
	}
	if !found {
		return fmt.Errorf("no valid appointment slot found for employee %s in company %s", employee.Created.ID.String(), company.Created.ID.String())
	}
	if err := handlerT.NewHttpClient().
		Method("PATCH").
		URL("/appointment/"+appointment_id).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, token).
		Header(namespace.HeadersKey.Company, company.Created.ID.String()).
		Send(map[string]any{
			"branch_id":  appointmentSlot.BranchID,
			"service_id": appointmentSlot.ServiceID,
			"start_time": appointmentSlot.StartTimeRFC3339,
		}).Error; err != nil {
		return fmt.Errorf("failed to reschedule appointment: %w", err)
	}
	if err := employee.GetById(200, nil, nil); err != nil {
		return err
	}
	if err := company.GetById(200, company.Owner.X_Auth_Token, nil); err != nil {
		return err
	}
	return nil
}

func CreateAppointmentRandomly(s int, company *modelT.Company, client *modelT.Client, employee *modelT.Employee, token, company_id string, a *modelT.Appointment) error {
	preferredLocation := time.UTC // Choose your timezone (e.g., UTC)
	appointmentSlot, found, err := FindValidAppointmentSlot(employee, company, preferredLocation)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("setup failed: could not find a valid appointment slot for initial booking")
	}
	http := handlerT.NewHttpClient()
	if err := http.
		Method("POST").
		URL("/appointment").
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, token).
		Header(namespace.HeadersKey.Company, company_id).
		Send(map[string]any{
			"branch_id":   appointmentSlot.BranchID,
			"service_id":  appointmentSlot.ServiceID,
			"employee_id": employee.Created.ID.String(),
			"company_id":  company.Created.ID.String(),
			"client_id":   client.Created.ID.String(),
			"start_time":  appointmentSlot.StartTimeRFC3339, // Use found start time
		}).Error; err != nil {
		return fmt.Errorf("failed to create appointment: %w", err)
	}
	if a != nil {
		http.ParseResponse(&a.Created)
	}
	if err := company.GetById(200, company.Owner.X_Auth_Token, nil); err != nil {
		return err
	}
	if err := client.GetByEmail(200); err != nil {
		return err
	}
	if err := employee.GetById(200, nil, nil); err != nil {
		return err
	}
	return nil
}

func GetAppointment(s int, appointment_id string, company_id, token string, a *modelT.Appointment) error {
	if err := handlerT.NewHttpClient().
		Method("GET").
		URL("/appointment/"+appointment_id).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, token).
		Header(namespace.HeadersKey.Company, company_id).
		Send(nil).
		ParseResponse(a).Error; err != nil {
		return fmt.Errorf("failed to get appointment by ID: %w", err)
	}
	return nil
}

func CancelAppointment(s int, appointment_id, company_id, token string) error {
	if err := handlerT.NewHttpClient().
		Method("DELETE").
		URL("/appointment/"+appointment_id).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, token).
		Header(namespace.HeadersKey.Company, company_id).
		Send(nil).Error; err != nil {
		return fmt.Errorf("failed to cancel appointment: %w", err)
	}
	return nil
}
