package models_test

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/tests/handlers"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Appointment struct {
	Created model.Appointment
}

func (a *Appointment) Create(status int, auth_token string, startTime *string, b *Branch, e *Employee, s *Service, cy *Company, ct *Client) error {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/appointment")
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Company, cy.Created.ID.String())
	http.Header(namespace.HeadersKey.Auth, auth_token)
	if startTime == nil {
		tempStartTime := lib.GenerateDateRFC3339(2027, 10, 29)
		startTime = &tempStartTime
	}
	A := DTO.CreateAppointment{
		BranchID:   b.Created.ID,
		EmployeeID: e.Created.ID,
		ServiceID:  s.Created.ID,
		ClientID:   ct.Created.ID,
		CompanyID:  cy.Created.ID,
		StartTime:  *startTime,
	}
	http.Send(A)
	http.ParseResponse(&a.Created)
	if err := b.GetById(200); err != nil {
		return err
	}
	if err := e.GetById(200); err != nil {
		return err
	}
	if err := s.GetById(200, nil); err != nil {
		return err
	}
	if err := cy.GetById(200); err != nil {
		return err
	}
	if err := ct.GetByEmail(200); err != nil {
		return err
	}
	var ClientAppointment mJSON.ClientAppointment
	aCreatedByte, err := json.Marshal(a.Created)
	if err != nil {
		return fmt.Errorf("failed to marshal appointment: %w", err)
	}
	err = json.Unmarshal(aCreatedByte, &ClientAppointment)
	if err != nil {
		return fmt.Errorf("failed to unmarshal appointment: %w", err)
	}
	ct.Created.Appointments.Add(&ClientAppointment)
	e.Created.Appointments = append(e.Created.Appointments, a.Created)
	b.Created.Appointments = append(b.Created.Appointments, a.Created)
	return nil
}

func (a *Appointment) CreateRandom(s int, company *Company, client *Client, employee *Employee, token, company_id string) error {
	if a.Created.ID != uuid.Nil {
		return fmt.Errorf("appointment already created with ID %s, cannot create again", a.Created.ID)
	}
	preferredLocation := time.UTC // Choose your timezone (e.g., UTC)
	appointmentSlot, found, err := FindValidAppointmentSlot(employee, company, preferredLocation)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("setup failed: could not find a valid appointment slot for initial booking")
	}
	http := handler.NewHttpClient()
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
		}).
		ParseResponse(&a.Created).Error; err != nil {
		return fmt.Errorf("failed to create appointment: %w", err)
	}
	if err := company.GetById(200); err != nil {
		return err
	}
	if err := client.GetByEmail(200); err != nil {
		return err
	}
	if err := employee.GetById(200); err != nil {
		return err
	}
	return nil
}

func (a *Appointment) GetById(s int, appointment_id, company_id, token string) error {
	if err := handler.NewHttpClient().
		Method("GET").
		URL("/appointment/"+appointment_id).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, token).
		Header(namespace.HeadersKey.Company, company_id).
		Send(nil).
		ParseResponse(&a.Created).Error; err != nil {
		return fmt.Errorf("failed to get appointment %s: %w", appointment_id, err)
	}
	return nil
}

func (a *Appointment) Cancel(s int, token string) error {
	if a.Created.ID == uuid.Nil {
		return fmt.Errorf("appointment not created, cannot cancel")
	}
	if err := handler.NewHttpClient().
		Method("DELETE").
		URL("/appointment/"+a.Created.ID.String()).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, token).
		Header(namespace.HeadersKey.Company, a.Created.CompanyID.String()).
		Send(nil).Error; err != nil {
		return fmt.Errorf("failed to cancel appointment %s: %w", a.Created.ID.String(), err)
	}
	return nil
}

type FoundAppointmentSlot struct {
	StartTimeRFC3339 string
	BranchID         string
	ServiceID        string
}

// Searches for the first available and valid appointment slot
// for a given employee within a company, considering branch, service, and existing appointment overlaps.
// It returns the necessary details or found=false if no valid slot is found.
// Constants for findValidAppointmentSlot logic
const (
	slotSearchHorizonDays = 14               // Example: search up to 2 weeks ahead
	slotSearchTimeStep    = 15 * time.Minute // Example: check every 15 minutes
)

func FindValidAppointmentSlot(employee *models_test.Employee, company *models_test.Company, preferredLocation *time.Location) (slot *FoundAppointmentSlot, found bool, err error) {
	if preferredLocation == nil {
		return nil, false, fmt.Errorf("preferredLocation is nil; timezone must be explicitly passed")
	}
	fmt.Printf("---- Starting findValidAppointmentSlot for Employee ID: %s ----", employee.Created.ID)
	fmt.Printf("PreferredLocation: %s, SlotSearchHorizon: %d days, TimeStep: %v", preferredLocation.String(), slotSearchHorizonDays, slotSearchTimeStep)

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
	fmt.Printf("Searching from %s (now in preferred TZ is %s)", searchStartDate.Format(time.RFC3339), nowInPreferredLocation.Format(time.RFC3339))

	for dayOffset := range slotSearchHorizonDays {
		currentSearchDate := searchStartDate.AddDate(0, 0, dayOffset)
		currentWeekday := currentSearchDate.Weekday()
		fmt.Printf("  Checking Date: %s (Weekday: %s, DayOffset: %d)", currentSearchDate.Format("2006-01-02"), currentWeekday, dayOffset)

		dayScheduleRanges, hasScheduleForDay := weekdaySchedules[currentWeekday]
		if !hasScheduleForDay || len(dayScheduleRanges) == 0 {
			fmt.Printf("    No schedule ranges for this day.")
			continue
		}
		fmt.Printf("    Found %d schedule range(s) for %s", len(dayScheduleRanges), currentWeekday)

		for wrIdx := range dayScheduleRanges {
			wr := &dayScheduleRanges[wrIdx]
			fmt.Printf("      Processing WorkRange #%d: Start='%s', End='%s', BranchID='%s'", wrIdx, wr.Start, wr.End, wr.BranchID)

			if wr.Start == "" || wr.End == "" || wr.BranchID == uuid.Nil {
				fmt.Printf("        Skipping WR: Invalid data (Start, End, or BranchID missing).")
				continue
			}

			var targetBranch *models_test.Branch
			branchObjFound := false
			for k := range company.Branches {
				if company.Branches[k].Created.ID == wr.BranchID {
					targetBranch = company.Branches[k]
					branchObjFound = true
					break
				}
			}
			if !branchObjFound {
				fmt.Printf("        Skipping WR: Branch object for ID %s not found in company.Branches.", wr.BranchID)
				continue
			}
			if !employeeIsAssignedToBranch(employee, wr.BranchID) {
				fmt.Printf("        Skipping WR: Employee %s not assigned to Branch %s (ID: %s).", employee.Created.ID, targetBranch.Created.Name, wr.BranchID)
				continue
			}
			fmt.Printf("        Branch '%s' (ID: %s) is valid and employee is assigned.", targetBranch.Created.Name, wr.BranchID)

			employeeServiceIDs := make(map[uuid.UUID]bool)
			if len(employee.Services) == 0 {
				fmt.Printf("        Skipping WR: Employee has no services assigned.")
				continue
			}
			for _, service := range employee.Services {
				employeeServiceIDs[service.Created.ID] = true
			}
			fmt.Printf("        Employee services: %v", employeeServiceIDs)

			if len(targetBranch.Services) == 0 {
				fmt.Printf("        Skipping WR: Branch has no services assigned.")
				continue
			}
			branchServiceIDs := make(map[uuid.UUID]bool)
			for _, service := range targetBranch.Services {
				branchServiceIDs[service.Created.ID] = true
			}
			fmt.Printf("        Branch services: %v", branchServiceIDs)

			validServiceIDs := []uuid.UUID{}
			for empSvcID := range employeeServiceIDs {
				if branchServiceIDs[empSvcID] {
					validServiceIDs = append(validServiceIDs, empSvcID)
				}
			}
			if len(validServiceIDs) == 0 {
				fmt.Printf("        Skipping WR: No common services between employee and branch.")
				continue
			}
			fmt.Printf("        Found %d common/valid services: %v", len(validServiceIDs), validServiceIDs)

			workRangeStartDateTime, err := parseTimeWithLocation(currentSearchDate, wr.Start, preferredLocation)
			if err != nil {
				fmt.Printf("        Skipping WR: Error parsing WorkRange Start Time '%s': %v", wr.Start, err)
				continue
			}
			workRangeEndDateTime, err := parseTimeWithLocation(currentSearchDate, wr.End, preferredLocation)
			if err != nil {
				fmt.Printf("        Skipping WR: Error parsing WorkRange End Time '%s': %v", wr.End, err)
				continue
			}
			if !workRangeStartDateTime.Before(workRangeEndDateTime) {
				fmt.Printf("        Skipping WR: WorkRange Start (%s) is not before End (%s).", workRangeStartDateTime, workRangeEndDateTime)
				continue
			}
			fmt.Printf("        Parsed WorkRange Times: Start=%s, End=%s", workRangeStartDateTime.Format(time.RFC3339), workRangeEndDateTime.Format(time.RFC3339))

			for _, selectedServiceID := range validServiceIDs {
				fmt.Printf("          Trying Service ID: %s", selectedServiceID)
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
					fmt.Printf("          Skipping Service: Definition for ID %s not found in company.Services. This is a test setup error.", selectedServiceID)
					continue // Should probably be t.Fatalf
				}
				if selectedServiceDurationMinutes == 0 {
					fmt.Printf("          Skipping Service: Duration is 0 minutes.")
					continue
				}
				fmt.Printf("          Service Duration: %d minutes", selectedServiceDurationMinutes)

				serviceDuration := time.Duration(selectedServiceDurationMinutes) * time.Minute

				for potentialStartTime := workRangeStartDateTime; potentialStartTime.Before(workRangeEndDateTime); potentialStartTime = potentialStartTime.Add(slotSearchTimeStep) {
					potentialEndTime := potentialStartTime.Add(serviceDuration)
					fmt.Printf("            Testing Slot: PotentialStart=%s, PotentialEnd=%s", potentialStartTime.Format(time.RFC3339), potentialEndTime.Format(time.RFC3339))

					// If the service cannot be completed by the end of the work range
					if potentialEndTime.After(workRangeEndDateTime) {
						fmt.Printf("              Slot ends after work range. Breaking from time-stepping for this service in this WR.")
						// No further time steps for THIS service in THIS work range will fit.
						break // Break from the time-stepping loop (for current service)
					}

					// Skip past slots
					if potentialStartTime.Before(nowInPreferredLocation) {
						fmt.Printf("              Skipping slot: Potential start time is in the past.")
						continue
					}

					overlap := false
					fmt.Printf("            Checking for overlaps with %d existing appointments...", len(employee.Created.Appointments))
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
						if !apptServiceFound || apptDur == 0 {
							fmt.Printf("              Warning: Could not determine duration for existing appt ID %s (ServiceID: %s). Skipping overlap with it.", appt.ID, appt.ServiceID)
							continue // Or t.Fatalf if this should always be present
						}

						existingStart := appt.StartTime.In(preferredLocation)
						existingEnd := existingStart.Add(time.Duration(apptDur) * time.Minute)
						fmt.Printf("              Comparing with existing appt [%s to %s]", existingStart.Format(time.RFC3339), existingEnd.Format(time.RFC3339))

						if potentialStartTime.Before(existingEnd) && potentialEndTime.After(existingStart) {
							fmt.Printf("                Overlap detected with existing appointment!")
							overlap = true
							break
						}
					}

					if !overlap {
						fmt.Printf("            SUCCESS! Found available slot: StartTime=%s, BranchID=%s, ServiceID=%s",
							potentialStartTime.Format(time.RFC3339), wr.BranchID.String(), selectedServiceID.String())
						slot = &FoundAppointmentSlot{
							StartTimeRFC3339: potentialStartTime.Format(time.RFC3339),
							BranchID:         wr.BranchID.String(),
							ServiceID:        selectedServiceID.String(),
						}
						return slot, true, nil // Found a slot
					} else {
						fmt.Printf("            Slot has overlap, trying next potential start time.")
					}
				} // End time-stepping loop
				fmt.Printf("          Finished time-stepping for service %s in this WR.", selectedServiceID)
			} // End service loop
		} // End work-range loop
	} // End day-offset loop

	fmt.Printf("---- Failed to find any valid appointment slot for Employee ID: %s after checking %d days. ----", employee.Created.ID, slotSearchHorizonDays)
	// The fallback to findNextAvailableSlotRFC3339 is okay for debugging, but the main logic should find it.
	if len(employee.Created.Appointments) > 0 {
		lastApptTime := employee.Created.Appointments[len(employee.Created.Appointments)-1].StartTime
		nextSuggestion, err := findNextAvailableSlotRFC3339(employee, lastApptTime.In(preferredLocation).Format(time.RFC3339))
		if err != nil {
			fmt.Printf("DEBUG: Fallback findNextAvailableSlotRFC3339 (based on last appt %s) failed: %v", lastApptTime.Format(time.RFC3339), err)
		} else {
			fmt.Printf("DEBUG: Fallback findNextAvailableSlotRFC3339 (based on last appt %s) suggests: %s", lastApptTime.Format(time.RFC3339), nextSuggestion)
		}
	} else {
		fmt.Printf("DEBUG: Employee has no existing appointments to base a fallback suggestion on.")
	}

	return slot, false, nil // No valid slot found
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

// findNextAvailableSlot attempts to find the next available start time for an employee
// based on their work schedule, starting from the one provided.
// NOTE: This is a simplified helper for testing and might not cover all edge cases.
// Renamed and modified to handle full timestamps
func findNextAvailableSlotRFC3339(employee *Employee, currentStartTimeRFC3339 string) (nextAvailableTimeRFC3339 string, err error) {

	layoutRFC3339 := time.RFC3339
	start, err := time.Parse(layoutRFC3339, currentStartTimeRFC3339)
	if err != nil {
		return "", fmt.Errorf("failed to parse current start time RFC3339 '%s': %w", currentStartTimeRFC3339, err)
	}

	// --- Service Duration calculation (remains similar) ---
	if len(employee.Services) == 0 {
		return "", fmt.Errorf("employee %s has no services assigned", employee.Created.ID)
	}
	duration := time.Duration(employee.Services[0].Created.Duration) * time.Minute
	nextPossibleStart := start.Add(duration)

	// --- Schedule Check (Now uses the date from input timestamp) ---
	schedule := employee.Created.WorkSchedule
	targetDate := start     // Use the date from the input timestamp
	loc := start.Location() // Preserve the timezone/location

	// Get the day schedule based on the start time's day of the week
	dayOfWeek := targetDate.Weekday()
	var daySchedule []mJSON.WorkRange
	switch dayOfWeek {
	case time.Monday:
		daySchedule = schedule.Monday
	case time.Tuesday:
		daySchedule = schedule.Tuesday
	case time.Wednesday:
		daySchedule = schedule.Wednesday
	case time.Thursday:
		daySchedule = schedule.Thursday
	case time.Friday:
		daySchedule = schedule.Friday
	case time.Saturday:
		daySchedule = schedule.Saturday
	case time.Sunday:
		daySchedule = schedule.Sunday
	}

	if len(daySchedule) == 0 {
		return "", fmt.Errorf("no work schedule found for employee %s on %s", employee.Created.ID, dayOfWeek)
	}

	timeLayout := "15:04" // Layout for parsing schedule times HH:MM

	for _, block := range daySchedule {
		blockStartTimeStr := block.Start
		blockEndTimeStr := block.End

		// Parse the block start/end times *relative to the targetDate's date and location*
		blockStartParsed, err := time.ParseInLocation(timeLayout, blockStartTimeStr, loc)
		if err != nil {
			fmt.Printf("Warn: bad block start time %s: %v", blockStartTimeStr, err)
			continue
		}
		blockEndParsed, err := time.ParseInLocation(timeLayout, blockEndTimeStr, loc)
		if err != nil {
			fmt.Printf("Warn: bad block end time %s: %v", blockEndTimeStr, err)
			continue
		}

		// Construct full datetime objects for the block boundaries on the target date
		blockStartDateTime := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(),
			blockStartParsed.Hour(), blockStartParsed.Minute(), 0, 0, loc)
		blockEndDateTime := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(),
			blockEndParsed.Hour(), blockEndParsed.Minute(), 0, 0, loc)

		// Check if the calculated nextPossibleStart fits within this block
		if !nextPossibleStart.Before(blockStartDateTime) && !nextPossibleStart.Add(duration).After(blockEndDateTime) {
			nextAvailableTimeRFC3339 = nextPossibleStart.Format(layoutRFC3339)
			return nextAvailableTimeRFC3339, nil // Return the RFC3339 string
		}
	}
	// If we reach here, no valid slot was found in the current day's schedule
	return "", fmt.Errorf("no valid slot found after %s for employee %s", currentStartTimeRFC3339, employee.Created.ID)
}

func employeeIsAssignedToBranch(e *Employee, branchID uuid.UUID) bool {
	for _, b := range e.Branches {
		if b.Created.ID == branchID {
			return true
		}
	}
	return false
}

func RescheduleAppointmentRandomly(s int, employee *Employee, company *Company, appointment_id, token string) error {
	preferredLocation := time.UTC // Choose your timezone (e.g., UTC)
	appointmentSlot, found, err := FindValidAppointmentSlot(employee, company, preferredLocation)
	if err != nil {
		return fmt.Errorf("failed to find valid appointment slot: %w", err)
	}
	if !found {
		return fmt.Errorf("no valid appointment slot found for employee %s in company %s", employee.Created.ID.String(), company.Created.ID.String())
	}
	if err := handler.NewHttpClient().
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
	if err := employee.GetById(200); err != nil {
		return err
	}
	if err := company.GetById(200); err != nil {
		return err
	}
	return nil
}
