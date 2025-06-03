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
// func FindValidAppointmentSlot(employee *modelT.Employee, company *modelT.Company, preferredLocation *time.Location) (slot *FoundAppointmentSlot, found bool, err error) {
// 	if preferredLocation == nil {
// 		return nil, false, fmt.Errorf("preferredLocation is nil; timezone must be explicitly passed")
// 	}
// 	fmt.Printf("---- Starting findValidAppointmentSlot for Employee ID: %s ----\n", employee.Created.ID)
// 	fmt.Printf("PreferredLocation: %s, SlotSearchHorizon: %d days, TimeStep: %v\n", preferredLocation.String(), slotSearchHorizonDays, slotSearchTimeStep)

// 	schedule := employee.Created.WorkSchedule
// 	weekdaySchedules := map[time.Weekday][]mJSON.WorkRange{
// 		time.Sunday:    schedule.Sunday,
// 		time.Monday:    schedule.Monday,
// 		time.Tuesday:   schedule.Tuesday,
// 		time.Wednesday: schedule.Wednesday,
// 		time.Thursday:  schedule.Thursday,
// 		time.Friday:    schedule.Friday,
// 		time.Saturday:  schedule.Saturday,
// 	}

// 	nowInPreferredLocation := time.Now().In(preferredLocation)
// 	searchStartDate := time.Date(nowInPreferredLocation.Year(), nowInPreferredLocation.Month(), nowInPreferredLocation.Day(), 0, 0, 0, 0, preferredLocation)
// 	fmt.Printf("Searching from %s (now in preferred TZ is %s)\n", searchStartDate.Format(time.RFC3339), nowInPreferredLocation.Format(time.RFC3339))

// 	for dayOffset := range slotSearchHorizonDays {
// 		currentSearchDate := searchStartDate.AddDate(0, 0, dayOffset)
// 		currentWeekday := currentSearchDate.Weekday()
// 		fmt.Printf("  Checking Date: %s (Weekday: %s, DayOffset: %d)\n", currentSearchDate.Format("2006-01-02"), currentWeekday, dayOffset)

// 		dayScheduleRanges, hasScheduleForDay := weekdaySchedules[currentWeekday]
// 		if !hasScheduleForDay || len(dayScheduleRanges) == 0 {
// 			fmt.Printf("    No schedule ranges for this day.\n")
// 			continue
// 		}
// 		fmt.Printf("    Found %d schedule range(s) for %s\n", len(dayScheduleRanges), currentWeekday)

// 		for wrIdx := range dayScheduleRanges {
// 			wr := &dayScheduleRanges[wrIdx]
// 			fmt.Printf("      Processing WorkRange #%d: Start='%s', End='%s', BranchID='%s'\n", wrIdx, wr.Start, wr.End, wr.BranchID)

// 			if wr.Start == "" || wr.End == "" || wr.BranchID == uuid.Nil {
// 				fmt.Printf("        Skipping WR: Invalid data (Start, End, or BranchID missing).\n")
// 				continue
// 			}

// 			var targetBranch *modelT.Branch
// 			branchObjFound := false
// 			for k := range company.Branches {
// 				if company.Branches[k].Created.ID == wr.BranchID {
// 					targetBranch = company.Branches[k]
// 					branchObjFound = true
// 					break
// 				}
// 			}
// 			if !branchObjFound {
// 				fmt.Printf("        Skipping WR: Branch object for ID %s not found in company.Branches.\n", wr.BranchID)
// 				continue
// 			}
// 			if !isEmployeeAssignedToBranch(employee, wr.BranchID) {
// 				fmt.Printf("        Skipping WR: Employee %s not assigned to Branch %s (ID: %s).\n", employee.Created.ID, targetBranch.Created.Name, wr.BranchID)
// 				continue
// 			}
// 			fmt.Printf("        Branch '%s' (ID: %s) is valid and employee is assigned.\n", targetBranch.Created.Name, wr.BranchID)

// 			employeeServiceIDs := make(map[uuid.UUID]bool)
// 			if len(employee.Services) == 0 {
// 				fmt.Printf("        Skipping WR: Employee has no services assigned.\n")
// 				continue
// 			}
// 			for _, service := range employee.Services {
// 				employeeServiceIDs[service.Created.ID] = true
// 			}
// 			fmt.Printf("        Employee services: %v\n", employeeServiceIDs)

// 			if len(targetBranch.Services) == 0 {
// 				fmt.Printf("        Skipping WR: Branch has no services assigned.\n")
// 				continue
// 			}
// 			branchServiceIDs := make(map[uuid.UUID]bool)
// 			for _, service := range targetBranch.Services {
// 				branchServiceIDs[service.Created.ID] = true
// 			}
// 			fmt.Printf("        Branch services: %v\n", branchServiceIDs)

// 			validServiceIDs := []uuid.UUID{}
// 			for empSvcID := range employeeServiceIDs {
// 				if branchServiceIDs[empSvcID] {
// 					validServiceIDs = append(validServiceIDs, empSvcID)
// 				}
// 			}
// 			if len(validServiceIDs) == 0 {
// 				fmt.Printf("        Skipping WR: No common services between employee and branch.\n")
// 				continue
// 			}
// 			fmt.Printf("        Found %d common/valid services: %v\n", len(validServiceIDs), validServiceIDs)

// 			workRangeStartDateTime, err := parseTimeWithLocation(currentSearchDate, wr.Start, preferredLocation)
// 			if err != nil {
// 				fmt.Printf("        Skipping WR: Error parsing WorkRange Start Time '%s': %v\n", wr.Start, err)
// 				continue
// 			}
// 			workRangeEndDateTime, err := parseTimeWithLocation(currentSearchDate, wr.End, preferredLocation)
// 			if err != nil {
// 				fmt.Printf("        Skipping WR: Error parsing WorkRange End Time '%s': %v\n", wr.End, err)
// 				continue
// 			}
// 			if !workRangeStartDateTime.Before(workRangeEndDateTime) {
// 				fmt.Printf("        Skipping WR: WorkRange Start (%s) is not before End (%s).\n", workRangeStartDateTime, workRangeEndDateTime)
// 				continue
// 			}
// 			fmt.Printf("        Parsed WorkRange Times: Start=%s, End=%s\n", workRangeStartDateTime.Format(time.RFC3339), workRangeEndDateTime.Format(time.RFC3339))

// 			for _, selectedServiceID := range validServiceIDs {
// 				fmt.Printf("          Trying Service ID: %s\n", selectedServiceID)
// 				var selectedServiceDurationMinutes uint
// 				serviceFound := false
// 				for _, s := range company.Services { // Search in company's master list of services
// 					if s.Created.ID == selectedServiceID {
// 						selectedServiceDurationMinutes = s.Created.Duration
// 						serviceFound = true
// 						break
// 					}
// 				}
// 				if !serviceFound {
// 					return nil, false, fmt.Errorf("service with ID %s not found in company's services", selectedServiceID)
// 				}
// 				if selectedServiceDurationMinutes == 0 {
// 					fmt.Printf("          Skipping Service: Duration is 0 minutes.\n")
// 					continue
// 				}
// 				fmt.Printf("          Service Duration: %d minutes\n", selectedServiceDurationMinutes)

// 				serviceDuration := time.Duration(selectedServiceDurationMinutes) * time.Minute

// 				for potentialStartTime := workRangeStartDateTime; potentialStartTime.Before(workRangeEndDateTime); potentialStartTime = potentialStartTime.Add(slotSearchTimeStep) {
// 					potentialEndTime := potentialStartTime.Add(serviceDuration)
// 					fmt.Printf("            Testing Slot: PotentialStart=%s, PotentialEnd=%s\n", potentialStartTime.Format(time.RFC3339), potentialEndTime.Format(time.RFC3339))

// 					// If the service cannot be completed by the end of the work range
// 					if potentialEndTime.After(workRangeEndDateTime) {
// 						fmt.Printf("              Slot ends after work range. Breaking from time-stepping for this service in this WR.\n")
// 						// No further time steps for THIS service in THIS work range will fit.
// 						break // Break from the time-stepping loop (for current service)
// 					}

// 					// Skip past slots
// 					if potentialStartTime.Before(nowInPreferredLocation) {
// 						fmt.Printf("              Skipping slot: Potential start time is in the past.\n")
// 						continue
// 					}

// 					overlap := false
// 					fmt.Printf("            Checking for overlaps with %d existing appointments...\n", len(employee.Created.Appointments))
// 					for _, appt := range employee.Created.Appointments {
// 						var apptDur uint
// 						apptServiceFound := false
// 						if appt.Service != nil && appt.Service.ID != uuid.Nil { // Ensure Service object and its ID are valid
// 							apptDur = appt.Service.Duration
// 							apptServiceFound = true
// 						} else { // Fallback to ServiceID
// 							for _, s := range company.Services {
// 								if s.Created.ID == appt.ServiceID {
// 									apptDur = s.Created.Duration
// 									apptServiceFound = true
// 									break
// 								}
// 							}
// 						}
// 						if !apptServiceFound {
// 							return nil, false, fmt.Errorf("appointment with ID %s has no valid service duration", appt.ID)
// 						} else if apptDur == 0 {
// 							return nil, false, fmt.Errorf("appointment with ID %s has a service duration of 0 minutes", appt.ID)
// 						}

// 						existingStart := appt.StartTime.In(preferredLocation)
// 						existingEnd := existingStart.Add(time.Duration(apptDur) * time.Minute)
// 						fmt.Printf("              Comparing with existing appt [%s to %s]\n", existingStart.Format(time.RFC3339), existingEnd.Format(time.RFC3339))

// 						if potentialStartTime.Before(existingEnd) && potentialEndTime.After(existingStart) {
// 							fmt.Printf("                Overlap detected with existing appointment!\n")
// 							overlap = true
// 							break
// 						}
// 					}

// 					if !overlap {
// 						fmt.Printf("            SUCCESS! Found available slot: StartTime=%s, BranchID=%s, ServiceID=%s\n",
// 							potentialStartTime.Format(time.RFC3339), wr.BranchID.String(), selectedServiceID.String())
// 						slot = &FoundAppointmentSlot{
// 							StartTimeRFC3339: potentialStartTime.Format(time.RFC3339),
// 							BranchID:         wr.BranchID.String(),
// 							ServiceID:        selectedServiceID.String(),
// 						}
// 						return slot, true, nil // Found a slot
// 					} else {
// 						fmt.Printf("            Slot has overlap, trying next potential start time.\n")
// 					}
// 				} // End time-stepping loop
// 				fmt.Printf("          Finished time-stepping for service %s in this WR.\n", selectedServiceID)
// 			} // End service loop
// 		} // End work-range loop
// 	} // End day-offset loop

// 	return slot, false, fmt.Errorf("no valid appointment slot found for employee %s in company %s within the search horizon", employee.Created.ID, company.Created.ID)
// }

func FindValidAppointmentSlotV2(employee *modelT.Employee, preferredLocation *time.Location) (*FoundAppointmentSlot, bool, error) {
	if preferredLocation == nil {
		return nil, false, fmt.Errorf("preferredLocation is nil; timezone must be explicitly passed")
	}

	fmt.Printf("---- Starting findValidAppointmentSlot for Employee ID: %s ----\n", employee.Created.ID.String())

	workSchedule := employee.Created.WorkSchedule
	weekdaySchedules := map[time.Weekday][]mJSON.WorkRange{
		time.Sunday:    workSchedule.Sunday,
		time.Monday:    workSchedule.Monday,
		time.Tuesday:   workSchedule.Tuesday,
		time.Wednesday: workSchedule.Wednesday,
		time.Thursday:  workSchedule.Thursday,
		time.Friday:    workSchedule.Friday,
		time.Saturday:  workSchedule.Saturday,
	}

	now := time.Now().In(preferredLocation)
	searchStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, preferredLocation)

	branchCache := make(map[string]*model.Branch)
	serviceCache := make(map[string]*model.Service)

	httpClient := handlerT.NewHttpClient().
		Header(namespace.HeadersKey.Company, employee.Company.Created.ID.String()).
		Header(namespace.HeadersKey.Auth, employee.Company.Owner.X_Auth_Token)

	for dayOffset := range slotSearchHorizonDays {
		currentDate := searchStart.AddDate(0, 0, dayOffset)
		currentWeekday := currentDate.Weekday()
		workRanges := weekdaySchedules[currentWeekday]

		for iWr, wr := range workRanges {
			if wr.Start == "" || wr.End == "" || wr.BranchID == uuid.Nil {
				return nil, false, fmt.Errorf("work range %d has invalid data (Start, End, or BranchID missing)", iWr)
			}

			branchID := wr.BranchID.String()
			branch, ok := branchCache[branchID]
			if !ok {
				var b model.Branch
				if err := httpClient.Method("GET").URL("/branch/" + branchID).
					ExpectedStatus(200).Send(nil).ParseResponse(&b).Error; err != nil {
					return nil, false, fmt.Errorf("failed to get branch %s: %w", branchID, err)
				}
				branchCache[branchID] = &b
				branch = &b
			}

			// Check if employee is assigned to the branch
			assignedToBranch := false
			for _, e := range branch.Employees {
				if e.ID == employee.Created.ID {
					assignedToBranch = true
					break
				}
			}
			if !assignedToBranch {
				return nil, false, fmt.Errorf("employee %s is not assigned to branch %s", employee.Created.ID, branchID)
			}

			startTime, err := parseTimeWithLocation(currentDate, wr.Start, preferredLocation)
			if err != nil {
				return nil, false, fmt.Errorf("failed to parse start time for work range #%d: %w", iWr, err)
			}
			endTime, err := parseTimeWithLocation(currentDate, wr.End, preferredLocation)
			if err != nil || !startTime.Before(endTime) {
				return nil, false, fmt.Errorf("invalid time range for work range #%d: %w", iWr, err)
			}

			for _, serviceID := range wr.Services {
				if serviceID == uuid.Nil {
					return nil, false, fmt.Errorf("work range %d has a nil service ID", iWr)
				}
				sID := serviceID.String()

				service, ok := serviceCache[sID]
				if !ok {
					var s model.Service
					if err := httpClient.Method("GET").URL("/service/" + sID).
						ExpectedStatus(200).Send(nil).ParseResponse(&s).Error; err != nil {
						return nil, false, fmt.Errorf("failed to get service %s: %w", sID, err)
					}
					serviceCache[sID] = &s
					service = &s
				}

				// Check if employee is assigned to the service
				assignedToService := false
				for _, e := range service.Employees {
					if e.ID == employee.Created.ID {
						assignedToService = true
						break
					}
				}
				if !assignedToService {
					return nil, false, fmt.Errorf("employee %s is not assigned to service %s", employee.Created.ID, sID)
				}

				// Check if branch is assigned to the service
				serviceAvailableAtBranch := false
				for _, s := range branch.Services {
					if s.ID == serviceID {
						serviceAvailableAtBranch = true
						break
					}
				}
				if !serviceAvailableAtBranch {
					return nil, false, fmt.Errorf("service %s is not available at branch %s", sID, branchID)
				}

				duration := time.Duration(service.Duration) * time.Minute

				for t := startTime; t.Add(duration).Before(endTime) || t.Add(duration).Equal(endTime); t = t.Add(slotSearchTimeStep) {
					if t.Before(now) {
						continue
					}
					tEnd := t.Add(duration)

					overlap := false
					for _, appt := range employee.Created.Appointments {
						start := appt.StartTime.In(preferredLocation)
						end := appt.EndTime
						if end.IsZero() && appt.Service != nil {
							end = start.Add(time.Duration(appt.Service.Duration) * time.Minute)
						}
						if start.Before(tEnd) && end.After(t) {
							overlap = true
							break
						}
					}
					if overlap {
						continue
					}

					return &FoundAppointmentSlot{
						StartTimeRFC3339: t.Format(time.RFC3339),
						BranchID:         branchID,
						ServiceID:        sID,
					}, true, nil
				}
			}
		}
	}

	return nil, false, fmt.Errorf("no valid appointment slot found for employee %s", employee.Created.ID.String())
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
	appointmentSlot, found, err := FindValidAppointmentSlotV2(employee, preferredLocation)
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
	appointmentSlot, found, err := FindValidAppointmentSlotV2(employee, preferredLocation)
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
