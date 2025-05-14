package e2e_test

import (
	"agenda-kaki-go/core"
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/tests/handlers"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

type Appointment struct {
	created model.Appointment
}

func Test_Appointment(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	ct := &Client{}
	ct.Create(t, 200)
	ct.VerifyEmail(t, 200)
	ct.Login(t, 200)
	ct.Update(t, 200, map[string]any{"name": "Updated client Name"})
	ct.GetByEmail(t, 200)
	cy := &Company{}
	cy.Set(t)
	b := cy.branches[0]
	e := cy.employees[0]
	s := cy.services[0]
	a := []*Appointment{}
	a = append(a, &Appointment{})
	a[0].Create(t, 200, ct.auth_token, nil, b, e, s, cy, ct)
	a = append(a, &Appointment{})
	a1StartTime := lib.GenerateDateRFC3339(2027, 10, 28)
	a[1].Create(t, 200, ct.auth_token, &a1StartTime, b, e, s, cy, ct)
	a = append(a, &Appointment{})
	a2StartTime := lib.GenerateDateRFC3339(2027, 10, 27)
	a[2].Create(t, 200, cy.owner.auth_token, &a2StartTime, b, e, s, cy, ct)
	startTimeStr := ct.created.Appointments[0].StartTime.Format(time.RFC3339)
	a = append(a, &Appointment{})
	a[3].Create(t, 400, ct.auth_token, &startTimeStr, b, e, s, cy, ct)
}

func (a *Appointment) Create(t *testing.T, status int, auth_token string, startTime *string, b *Branch, e *Employee, s *Service, cy *Company, ct *Client) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/appointment")
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Company, cy.created.ID.String())
	http.Header(namespace.HeadersKey.Auth, auth_token)
	if startTime == nil {
		tempStartTime := lib.GenerateDateRFC3339(2027, 10, 29)
		startTime = &tempStartTime
	}
	A := DTO.CreateAppointment{
		BranchID:   b.created.ID,
		EmployeeID: e.created.ID,
		ServiceID:  s.created.ID,
		ClientID:   ct.created.ID,
		CompanyID:  cy.created.ID,
		StartTime:  *startTime,
	}
	http.Send(A)
	http.ParseResponse(&a.created)
	b.GetById(t, 200)
	e.GetById(t, 200)
	s.GetById(t, 200)
	cy.GetById(t, 200)
	ct.GetByEmail(t, 200)
	var ClientAppointment mJSON.ClientAppointment
	aCreatedByte, err := json.Marshal(a.created)
	if err != nil {
		t.Fatalf("Failed to marshal appointment: %v", err)
	}
	err = json.Unmarshal(aCreatedByte, &ClientAppointment)
	if err != nil {
		t.Fatalf("Failed to unmarshal appointment: %v", err)
	}
	ct.created.Appointments.Add(&ClientAppointment)
	e.created.Appointments = append(e.created.Appointments, a.created)
	b.created.Appointments = append(b.created.Appointments, a.created)
}

type FoundAppointmentSlot struct {
	StartTimeRFC3339 string
	BranchID         string
	ServiceID        string
}

// Searches for the first available and valid appointment slot
// for a given employee within a company, considering branch, service, and existing appointment overlaps.
// It returns the necessary details or found=false if no valid slot is found.
func findValidAppointmentSlot(t *testing.T, employee *Employee, company *Company, preferredLocation *time.Location) (slot FoundAppointmentSlot, found bool) {
	t.Helper()

	// 1. Iterate through Employee's Schedule
	schedule := employee.created.WorkSchedule // Assumes .created is model.Employee
	weekdaySchedules := map[time.Weekday][]mJSON.WorkRange{
		time.Monday:    schedule.Monday,
		time.Tuesday:   schedule.Tuesday,
		time.Wednesday: schedule.Wednesday,
		time.Thursday:  schedule.Thursday,
		time.Friday:    schedule.Friday,
		time.Saturday:  schedule.Saturday,
		time.Sunday:    schedule.Sunday,
	}
	checkOrder := []time.Weekday{
		time.Monday, time.Tuesday, time.Wednesday, time.Thursday,
		time.Friday, time.Saturday, time.Sunday,
	}

	for _, day := range checkOrder { // Loop through days of the week
		daySchedule := weekdaySchedules[day]
		if len(daySchedule) == 0 {
			continue // No work ranges for this day
		}

		for i := range daySchedule { // Loop through work ranges within the day
			wr := &daySchedule[i] // Pointer to mJSON.WorkRange

			// 2. Check if WorkRange is potentially valid
			if wr.Start == "" || wr.BranchID == uuid.Nil {
				continue // Skip invalid ranges (no start time or no branch ID)
			}
			t.Logf("Checking potential slot: Time=%s on Weekday=%s at Branch=%s", wr.Start, day, wr.BranchID)

			// 3. Find the target Branch object from the setup data
			var targetBranch *Branch // Use your test Branch helper type
			branchFound := false
			for k := range company.branches { // Assumes company.branches []*Branch
				// Assumes Branch helper has .created field of type model.Branch
				if company.branches[k].created.ID == wr.BranchID {
					targetBranch = company.branches[k]
					branchFound = true
					break
				}
			}
			if !branchFound {
				t.Logf("Warning: Branch object with ID %s (from schedule) not found in company setup data. Skipping this work range.", wr.BranchID)
				continue // Skip this work range if branch data is inconsistent
			}
			t.Logf("Found Branch Object: %s", targetBranch.created.Name)

			// 4. Determine Valid Service IDs (Employee + Branch intersection)
			// 4a. Employee services map
			employeeServiceIDs := make(map[uuid.UUID]bool)
			// Assumes employee.services is []*Service helper type
			if len(employee.services) == 0 {
				// If employee must have services, could continue to next work range here.
				// Or maybe this work range is unusable if no services can be offered.
				t.Logf("Debug: Employee %s has no services assigned in test setup. Skipping this work range.", employee.created.ID)
				continue // Skip this work range as employee can't do anything
			}
			for _, service := range employee.services {
				// Assumes service.created is model.Service
				employeeServiceIDs[service.created.ID] = true
			}

			// 4b. Branch services map (Requires branch.services []*Service helper type in setup)
			branchServiceIDs := make(map[uuid.UUID]bool)
			if len(targetBranch.services) == 0 {
				// Strict check: branch must explicitly list offered services
				t.Logf("Warning: Branch %s (%s) has no specific services assigned in setup. Cannot validate service for this work range. Skipping.", targetBranch.created.Name, targetBranch.created.ID)
				continue // Skip this work range if branch definition is incomplete
			} else {
				for _, service := range targetBranch.services {
					branchServiceIDs[service.created.ID] = true
				}
			}

			// 4c. Find intersection
			validServiceIDs := []uuid.UUID{}
			for empSvcID := range employeeServiceIDs {
				if branchServiceIDs[empSvcID] {
					validServiceIDs = append(validServiceIDs, empSvcID)
				}
			}

			// If no common service for this employee/branch combination
			if len(validServiceIDs) == 0 {
				t.Logf("No common services found for employee %s and branch %s. Skipping this work range.", employee.created.ID, targetBranch.created.ID)
				continue // Skip this work range
			}
			t.Logf("Found %d potential common services for this slot: %v", len(validServiceIDs), validServiceIDs)

			// --- Loop through potentially valid services for this specific WorkRange ---
			for _, selectedServiceID := range validServiceIDs {
				t.Logf("  Attempting to validate service %s for slot %s on %s...", selectedServiceID, wr.Start, day)

				// 5. Get Duration for the *selected* service
				var selectedServiceDuration uint // Assuming duration is uint (minutes)
				serviceDurationFound := false
				// Search within the company's services list (most reliable place for master data)
				// Assumes company.services holds []*Service helpers
				for _, s := range company.services {
					if s.created.ID == selectedServiceID { // Assumes s.created is model.Service
						selectedServiceDuration = s.created.Duration // Assumes Duration field exists
						serviceDurationFound = true
						t.Logf("  Duration for service %s is %d minutes", selectedServiceID, selectedServiceDuration)
						break
					}
				}
				if !serviceDurationFound {
					// This is critical. If the service ID was valid enough to be in the intersection,
					// its definition should exist in the company setup. Failure implies bad test data.
					t.Fatalf("Critical Error: Service definition/duration for ID %s not found in company setup data. Check SetupRandomized.", selectedServiceID)
					// No point continuing if we can't get duration
				}

				// 6. Calculate Potential Appointment Time Range
				targetDate := findNextWeekday(time.Now(), day)
				potentialStartTime, err := parseTimeWithLocation(t, targetDate, wr.Start, preferredLocation)
				if err != nil {
					// If we can't parse the start time from the schedule, skip this work range entry
					t.Logf("Error parsing potential start time %s: %v. Skipping this WorkRange entry.", wr.Start, err)
					// This `continue` jumps to the next iteration of the *innermost* loop (validServiceIDs).
					// We should probably break out of the service loop and try the next *work range*?
					// Let's continue to the next service first, maybe another service works. If time parsing fails for *all*, we move on naturally.
					continue // Continue to the next potential service ID for this work range
				}
				potentialEndTime := potentialStartTime.Add(time.Duration(selectedServiceDuration) * time.Minute)
				t.Logf("  Potential new appt range: [%s, %s)", potentialStartTime.Format(time.RFC3339), potentialEndTime.Format(time.RFC3339))

				// 7. *** OVERLAP CHECK ***
				hasOverlap := false
				// Assumes employee.created is model.Employee and has Appointments loaded
				// Assumes model.Appointment has StartTime (time.Time) and Service preloaded or a Duration field
				if len(employee.created.Appointments) > 0 {
					t.Logf("  Checking against %d existing appointments for overlaps...", len(employee.created.Appointments))
					for _, existingAppt := range employee.created.Appointments {
						var existingApptDuration uint
						// Determine duration of the existing appointment
						if existingAppt.Service != nil { // Prefer preloaded Service
							existingApptDuration = existingAppt.Service.Duration
						} else {
							// Fallback: Need to look up service duration based on existingAppt.ServiceID
							// Requires searching company.services again, adds complexity
							lookupDurationOk := false
							for _, s := range company.services {
								if s.created.ID == existingAppt.ServiceID {
									existingApptDuration = s.created.Duration
									lookupDurationOk = true
									break
								}
							}
							if !lookupDurationOk {
								t.Logf("  Warning: Could not find service %s to determine duration for existing appointment %s. Skipping overlap check against this one.", existingAppt.ServiceID, existingAppt.ID)
								continue // Skip check against this specific appointment
							}
						}

						// Perform comparison in the consistent preferred location
						existingStartTime := existingAppt.StartTime.In(preferredLocation)
						existingEndTime := existingStartTime.Add(time.Duration(existingApptDuration) * time.Minute)

						// The overlap condition: (newStart < existingEnd) && (newEnd > existingStart)
						if potentialStartTime.Before(existingEndTime) && potentialEndTime.After(existingStartTime) {
							hasOverlap = true
							t.Logf("  Overlap detected! Slot conflicts with existing [%s, %s) (Appt ID: %s)",
								existingStartTime.Format(time.RFC1123), existingEndTime.Format(time.RFC1123),
								existingAppt.ID,
							)
							break // Stop checking other existing appointments; this service/time combo is invalid
						}
					} // End overlap check loop
				} else {
					t.Logf("  No existing appointments found for employee. Skipping overlap check.")
				}

				// 8. If NO overlap was detected for THIS SERVICE at THIS TIME
				if !hasOverlap {
					t.Logf("  Success! No overlap found for service %s at %s", selectedServiceID, potentialStartTime.Format(time.RFC3339))

					// We found a fully valid slot! Populate and return.
					slot.StartTimeRFC3339 = potentialStartTime.Format(time.RFC3339)
					slot.BranchID = wr.BranchID.String()
					slot.ServiceID = selectedServiceID.String()
					found = true // Mark as found

					t.Logf("----> Found valid slot details: Time=%s, Branch=%s, Service=%s", slot.StartTimeRFC3339, slot.BranchID, slot.ServiceID)
					return slot, found // Return the first fully valid slot found
				} else {
					// If overlap WAS detected, loop continues to try the next valid service ID for the *same* work range slot
					t.Logf("  Overlap detected for service %s, trying next valid service if available for this slot.", selectedServiceID)
				}

			} // --- End loop through valid services for this WorkRange slot ---
			// If we finish the service loop without returning, it means all valid services for *this specific* wr.Start time had overlaps.
			t.Logf("Finished checking all %d valid services for slot %s on %s. None were free of overlaps.", len(validServiceIDs), wr.Start, day)

		} // --- End loop WorkRanges within a day ---
	} // --- End loop days of week ---

	// If all loops complete without returning, no valid slot was found
	t.Logf("Failed to find any valid, non-overlapping appointment slot for employee %s meeting all criteria.", employee.created.ID)
	return slot, false // found is already false
}

// Helper to parse HH:MM or HH:MM:SS time string into a full time.Time on a specific date/location
func parseTimeWithLocation(t *testing.T, targetDate time.Time, timeStr string, loc *time.Location) (time.Time, error) {
	t.Helper()
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

// findNextWeekday finds the first date occurrence of targetWeekday strictly after startAfter.
func findNextWeekday(startAfter time.Time, targetWeekday time.Weekday) time.Time {
	// Start checking from the day *after* startAfter
	currentDate := startAfter.AddDate(0, 0, 1)
	for i := 0; i < 7; i++ { // Loop a maximum of 7 days to find the next occurrence
		if currentDate.Weekday() == targetWeekday {
			return currentDate
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}
	// Should realistically always find within 7 days, but fallback if needed
	// This indicates a potential logic error if reached.
	// Depending on test strictness, you might want to panic or t.Fatal here.
	return startAfter.AddDate(0, 0, 7) // Fallback: return a week later
}

// createRFC3339Timestamp creates an RFC3339 formatted timestamp string for a given date and HH:MM time.
func createRFC3339Timestamp(t *testing.T, targetDate time.Time, timeStr string, loc *time.Location) string {
	t.Helper()
	layout := "15:04"                                             // HH:MM format
	parsedTime, err := time.ParseInLocation(layout, timeStr, loc) // Use ParseInLocation
	if err != nil {
		t.Fatalf("Failed to parse time string '%s': %v", timeStr, err)
	}

	// Combine the target date with the parsed time components
	appointmentDateTime := time.Date(
		targetDate.Year(),
		targetDate.Month(),
		targetDate.Day(),
		parsedTime.Hour(),
		parsedTime.Minute(),
		0, 0, // Seconds and Nanoseconds
		loc, // Use the specified location
	)
	return appointmentDateTime.Format(time.RFC3339) // Format to RFC3339
}

// findNextAvailableSlot attempts to find the next available start time for an employee
// based on their work schedule, starting from the one provided.
// NOTE: This is a simplified helper for testing and might not cover all edge cases.
// Renamed and modified to handle full timestamps
func findNextAvailableSlotRFC3339(t *testing.T, employee *Employee, currentStartTimeRFC3339 string) string {
	t.Helper()
	layoutRFC3339 := time.RFC3339
	start, err := time.Parse(layoutRFC3339, currentStartTimeRFC3339)
	if err != nil {
		t.Fatalf("Failed to parse current start time RFC3339 '%s': %v", currentStartTimeRFC3339, err)
	}

	// --- Service Duration calculation (remains similar) ---
	if len(employee.services) == 0 {
		t.Fatal("Employee has no services, cannot determine next slot based on duration.")
	}
	duration := time.Duration(employee.services[0].created.Duration) * time.Minute
	nextPossibleStart := start.Add(duration)

	// --- Schedule Check (Now uses the date from input timestamp) ---
	schedule := employee.created.WorkSchedule
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
		t.Logf("Warning: No work schedule found for employee %s on %s. Cannot find next slot.", employee.created.ID, dayOfWeek)
		return currentStartTimeRFC3339 // Fallback
	}

	timeLayout := "15:04" // Layout for parsing schedule times HH:MM

	for _, block := range daySchedule {
		blockStartTimeStr := block.Start
		blockEndTimeStr := block.End

		// Parse the block start/end times *relative to the targetDate's date and location*
		blockStartParsed, err := time.ParseInLocation(timeLayout, blockStartTimeStr, loc)
		if err != nil {
			t.Logf("Warn: bad block start time %s: %v", blockStartTimeStr, err)
			continue
		}
		blockEndParsed, err := time.ParseInLocation(timeLayout, blockEndTimeStr, loc)
		if err != nil {
			t.Logf("Warn: bad block end time %s: %v", blockEndTimeStr, err)
			continue
		}

		// Construct full datetime objects for the block boundaries on the target date
		blockStartDateTime := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(),
			blockStartParsed.Hour(), blockStartParsed.Minute(), 0, 0, loc)
		blockEndDateTime := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(),
			blockEndParsed.Hour(), blockEndParsed.Minute(), 0, 0, loc)

		// Check if the calculated nextPossibleStart fits within this block
		if !nextPossibleStart.Before(blockStartDateTime) && !nextPossibleStart.Add(duration).After(blockEndDateTime) {
			return nextPossibleStart.Format(layoutRFC3339) // Return the RFC3339 string
		}
	}

	t.Logf("Warning: Could not find next available slot after %s for employee %s. Returning original time.", currentStartTimeRFC3339, employee.created.ID)
	return currentStartTimeRFC3339 // Fallback
}
