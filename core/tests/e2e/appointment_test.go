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
	cy.Set(t) // This sets up company, employees (with schedules), branches, services.

	// We will primarily use one employee for these tests.
	// The findValidAppointmentSlot function will determine suitable branch and service.
	if len(cy.employees) == 0 {
		t.Fatalf("Test setup failed: No employees created by cy.Set(t)")
	}
	baseEmployee := cy.employees[0]

	a := []*Appointment{}

	// --- Test Case 0: Successful creation by client ---
	a = append(a, &Appointment{})
	// Find a valid slot for the base employee. Using time.Local for preferred location.
	slot0, found0 := findValidAppointmentSlot(t, baseEmployee, cy, time.Local)
	if !found0 {
		t.Fatalf("Test setup failed: Could not find any valid appointment slot for employee %s for test case a[0]", baseEmployee.created.ID)
	}
	// Retrieve the actual Branch and Service objects based on IDs from slot0
	branchForSlot0 := getBranchByID(t, cy, slot0.BranchID)
	serviceForSlot0 := getServiceByID(t, cy, slot0.ServiceID)
	a[0].Create(t, 200, ct.auth_token, &slot0.StartTimeRFC3339, branchForSlot0, baseEmployee, serviceForSlot0, cy, ct)

	// --- Test Case 1: Another successful creation by client ---
	// The employee's appointments list (baseEmployee.created.Appointments) should have been updated by a[0].Create(),
	// so findValidAppointmentSlot should now find the *next* available slot.
	a = append(a, &Appointment{})
	slot1, found1 := findValidAppointmentSlot(t, baseEmployee, cy, time.Local)
	if !found1 {
		t.Fatalf("Test setup failed: Could not find a second valid appointment slot for employee %s for test case a[1]", baseEmployee.created.ID)
	}
	branchForSlot1 := getBranchByID(t, cy, slot1.BranchID)
	serviceForSlot1 := getServiceByID(t, cy, slot1.ServiceID)
	a[1].Create(t, 200, ct.auth_token, &slot1.StartTimeRFC3339, branchForSlot1, baseEmployee, serviceForSlot1, cy, ct)

	// --- Test Case 2: Successful creation by company owner ---
	a = append(a, &Appointment{})
	slot2, found2 := findValidAppointmentSlot(t, baseEmployee, cy, time.Local)
	if !found2 {
		t.Fatalf("Test setup failed: Could not find a third valid appointment slot for employee %s for test case a[2]", baseEmployee.created.ID)
	}
	branchForSlot2 := getBranchByID(t, cy, slot2.BranchID)
	serviceForSlot2 := getServiceByID(t, cy, slot2.ServiceID)
	a[2].Create(t, 200, cy.owner.auth_token, &slot2.StartTimeRFC3339, branchForSlot2, baseEmployee, serviceForSlot2, cy, ct)

	// --- Test Case 3: Attempt to create conflicting appointment (expects 400) ---
	// This test uses the details of the first successfully created appointment (a[0]) to force a conflict.
	if a[0].created.ID == uuid.Nil {
		t.Fatalf("Prerequisite failed for Test Case 3: a[0].created appointment is nil. Cannot test conflict.")
	}
	// The start time for the conflict is the same as a[0]'s start time.
	startTimeForConflict := a[0].created.StartTime.Format(time.RFC3339)
	// The branch, employee, service must be the same as a[0] to ensure a direct conflict.
	// branchForSlot0, baseEmployee, serviceForSlot0 are already the correct objects.
	a = append(a, &Appointment{})
	a[3].Create(t, 409, ct.auth_token, &startTimeForConflict, branchForSlot0, baseEmployee, serviceForSlot0, cy, ct)
}

// Helper functions to retrieve Branch and Service objects from Company test setup data
func getBranchByID(t *testing.T, company *Company, branchIDStr string) *Branch {
	t.Helper()
	branchUUID, err := uuid.Parse(branchIDStr)
	if err != nil {
		t.Fatalf("Invalid Branch ID string from slot finder: %s, error: %v", branchIDStr, err)
	}
	for _, br := range company.branches {
		if br.created.ID == branchUUID {
			return br
		}
	}
	t.Fatalf("Test setup error: Branch with ID %s (found by slot finder) not in company.branches", branchIDStr)
	return nil
}

func getServiceByID(t *testing.T, company *Company, serviceIDStr string) *Service {
	t.Helper()
	serviceUUID, err := uuid.Parse(serviceIDStr)
	if err != nil {
		t.Fatalf("Invalid Service ID string from slot finder: %s, error: %v", serviceIDStr, err)
	}
	for _, serv := range company.services { // Assuming company.services holds all services
		if serv.created.ID == serviceUUID {
			return serv
		}
	}
	// It's possible findValidAppointmentSlot finds services associated directly with employee/branch,
	// ensure company.services is comprehensive or adjust where to look for the service object.
	// For now, assuming company.services is the master list for the test.
	t.Fatalf("Test setup error: Service with ID %s (found by slot finder) not in company.services", serviceIDStr)
	return nil
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
// Constants for findValidAppointmentSlot logic
const (
	slotSearchHorizonDays = 14            // Example: search up to 2 weeks ahead
	slotSearchTimeStep    = 15 * time.Minute // Example: check every 15 minutes
)

func findValidAppointmentSlot(t *testing.T, employee *Employee, company *Company, preferredLocation *time.Location) (slot FoundAppointmentSlot, found bool) {
	t.Helper()

	if preferredLocation == nil {
		t.Fatalf("preferredLocation is nil; timezone must be explicitly passed")
	}
	t.Logf("---- Starting findValidAppointmentSlot for Employee ID: %s ----", employee.created.ID)
	t.Logf("PreferredLocation: %s, SlotSearchHorizon: %d days, TimeStep: %v", preferredLocation.String(), slotSearchHorizonDays, slotSearchTimeStep)


	schedule := employee.created.WorkSchedule
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
	t.Logf("Searching from %s (now in preferred TZ is %s)", searchStartDate.Format(time.RFC3339), nowInPreferredLocation.Format(time.RFC3339))


	for dayOffset := range slotSearchHorizonDays {
		currentSearchDate := searchStartDate.AddDate(0, 0, dayOffset)
		currentWeekday := currentSearchDate.Weekday()
		t.Logf("  Checking Date: %s (Weekday: %s, DayOffset: %d)", currentSearchDate.Format("2006-01-02"), currentWeekday, dayOffset)

		dayScheduleRanges, hasScheduleForDay := weekdaySchedules[currentWeekday]
		if !hasScheduleForDay || len(dayScheduleRanges) == 0 {
			t.Logf("    No schedule ranges for this day.")
			continue
		}
		t.Logf("    Found %d schedule range(s) for %s", len(dayScheduleRanges), currentWeekday)

		for wrIdx := range dayScheduleRanges {
			wr := &dayScheduleRanges[wrIdx]
			t.Logf("      Processing WorkRange #%d: Start='%s', End='%s', BranchID='%s'", wrIdx, wr.Start, wr.End, wr.BranchID)

			if wr.Start == "" || wr.End == "" || wr.BranchID == uuid.Nil {
				t.Logf("        Skipping WR: Invalid data (Start, End, or BranchID missing).")
				continue
			}

			var targetBranch *Branch
			branchObjFound := false
			for k := range company.branches {
				if company.branches[k].created.ID == wr.BranchID {
					targetBranch = company.branches[k]
					branchObjFound = true
					break
				}
			}
			if !branchObjFound {
				t.Logf("        Skipping WR: Branch object for ID %s not found in company.branches.", wr.BranchID)
				continue
			}
			if !employeeIsAssignedToBranch(employee, wr.BranchID) {
                t.Logf("        Skipping WR: Employee %s not assigned to Branch %s (ID: %s).", employee.created.ID, targetBranch.created.Name, wr.BranchID)
				continue
			}
			t.Logf("        Branch '%s' (ID: %s) is valid and employee is assigned.", targetBranch.created.Name, wr.BranchID)


			employeeServiceIDs := make(map[uuid.UUID]bool)
			if len(employee.services) == 0 {
				t.Logf("        Skipping WR: Employee has no services assigned.")
				continue
			}
			for _, service := range employee.services {
				employeeServiceIDs[service.created.ID] = true
			}
			t.Logf("        Employee services: %v", employeeServiceIDs)


			if len(targetBranch.services) == 0 {
				t.Logf("        Skipping WR: Branch has no services assigned.")
				continue
			}
			branchServiceIDs := make(map[uuid.UUID]bool)
			for _, service := range targetBranch.services {
				branchServiceIDs[service.created.ID] = true
			}
			t.Logf("        Branch services: %v", branchServiceIDs)


			validServiceIDs := []uuid.UUID{}
			for empSvcID := range employeeServiceIDs {
				if branchServiceIDs[empSvcID] {
					validServiceIDs = append(validServiceIDs, empSvcID)
				}
			}
			if len(validServiceIDs) == 0 {
				t.Logf("        Skipping WR: No common services between employee and branch.")
				continue
			}
			t.Logf("        Found %d common/valid services: %v", len(validServiceIDs), validServiceIDs)


			workRangeStartDateTime, err := parseTimeWithLocation(t, currentSearchDate, wr.Start, preferredLocation)
			if err != nil {
				t.Logf("        Skipping WR: Error parsing WorkRange Start Time '%s': %v", wr.Start, err)
				continue
			}
			workRangeEndDateTime, err := parseTimeWithLocation(t, currentSearchDate, wr.End, preferredLocation)
			if err != nil {
				t.Logf("        Skipping WR: Error parsing WorkRange End Time '%s': %v", wr.End, err)
				continue
			}
			if !workRangeStartDateTime.Before(workRangeEndDateTime) {
				t.Logf("        Skipping WR: WorkRange Start (%s) is not before End (%s).", workRangeStartDateTime, workRangeEndDateTime)
				continue
			}
			t.Logf("        Parsed WorkRange Times: Start=%s, End=%s", workRangeStartDateTime.Format(time.RFC3339), workRangeEndDateTime.Format(time.RFC3339))


			for _, selectedServiceID := range validServiceIDs {
				t.Logf("          Trying Service ID: %s", selectedServiceID)
				var selectedServiceDurationMinutes uint
				serviceFound := false
				for _, s := range company.services { // Search in company's master list of services
					if s.created.ID == selectedServiceID {
						selectedServiceDurationMinutes = s.created.Duration
						serviceFound = true
						break
					}
				}
				if !serviceFound {
					t.Logf("          Skipping Service: Definition for ID %s not found in company.services. This is a test setup error.", selectedServiceID)
					continue // Should probably be t.Fatalf
				}
				if selectedServiceDurationMinutes == 0 {
					t.Logf("          Skipping Service: Duration is 0 minutes.")
					continue
				}
				t.Logf("          Service Duration: %d minutes", selectedServiceDurationMinutes)

				serviceDuration := time.Duration(selectedServiceDurationMinutes) * time.Minute

				for potentialStartTime := workRangeStartDateTime; potentialStartTime.Before(workRangeEndDateTime); potentialStartTime = potentialStartTime.Add(slotSearchTimeStep) {
					potentialEndTime := potentialStartTime.Add(serviceDuration)
					t.Logf("            Testing Slot: PotentialStart=%s, PotentialEnd=%s", potentialStartTime.Format(time.RFC3339), potentialEndTime.Format(time.RFC3339))

					// If the service cannot be completed by the end of the work range
					if potentialEndTime.After(workRangeEndDateTime) {
						t.Logf("              Slot ends after work range. Breaking from time-stepping for this service in this WR.")
						// No further time steps for THIS service in THIS work range will fit.
						break // Break from the time-stepping loop (for current service)
					}

					// Skip past slots
					if potentialStartTime.Before(nowInPreferredLocation) {
						t.Logf("              Skipping slot: Potential start time is in the past.")
						continue
					}

					overlap := false
					t.Logf("            Checking for overlaps with %d existing appointments...", len(employee.created.Appointments))
					for _, appt := range employee.created.Appointments {
						var apptDur uint
						apptServiceFound := false
						if appt.Service != nil && appt.Service.ID != uuid.Nil { // Ensure Service object and its ID are valid
							apptDur = appt.Service.Duration
							apptServiceFound = true
						} else { // Fallback to ServiceID
							for _, s := range company.services {
								if s.created.ID == appt.ServiceID {
									apptDur = s.created.Duration
									apptServiceFound = true
									break
								}
							}
						}
						if !apptServiceFound || apptDur == 0 {
                            t.Logf("              Warning: Could not determine duration for existing appt ID %s (ServiceID: %s). Skipping overlap with it.", appt.ID, appt.ServiceID)
                            continue // Or t.Fatalf if this should always be present
                        }

						existingStart := appt.StartTime.In(preferredLocation)
						existingEnd := existingStart.Add(time.Duration(apptDur) * time.Minute)
						t.Logf("              Comparing with existing appt [%s to %s]", existingStart.Format(time.RFC3339), existingEnd.Format(time.RFC3339))

						if potentialStartTime.Before(existingEnd) && potentialEndTime.After(existingStart) {
							t.Logf("                Overlap detected with existing appointment!")
							overlap = true
							break
						}
					}

					if !overlap {
						t.Logf("            SUCCESS! Found available slot: StartTime=%s, BranchID=%s, ServiceID=%s",
							potentialStartTime.Format(time.RFC3339), wr.BranchID.String(), selectedServiceID.String())
						slot = FoundAppointmentSlot{
							StartTimeRFC3339: potentialStartTime.Format(time.RFC3339),
							BranchID:         wr.BranchID.String(),
							ServiceID:        selectedServiceID.String(),
						}
						return slot, true // Found a slot
					} else {
						t.Logf("            Slot has overlap, trying next potential start time.")
					}
				} // End time-stepping loop
				t.Logf("          Finished time-stepping for service %s in this WR.", selectedServiceID)
			} // End service loop
		} // End work-range loop
	} // End day-offset loop

	t.Logf("---- Failed to find any valid appointment slot for Employee ID: %s after checking %d days. ----", employee.created.ID, slotSearchHorizonDays)
	// The fallback to findNextAvailableSlotRFC3339 is okay for debugging, but the main logic should find it.
	if len(employee.created.Appointments) > 0 {
		lastApptTime := employee.created.Appointments[len(employee.created.Appointments)-1].StartTime
		nextSuggestion := findNextAvailableSlotRFC3339(t, employee, lastApptTime.In(preferredLocation).Format(time.RFC3339))
		t.Logf("DEBUG: Fallback findNextAvailableSlotRFC3339 (based on last appt %s) suggests: %s", lastApptTime.Format(time.RFC3339), nextSuggestion)
	} else {
		t.Logf("DEBUG: Employee has no existing appointments to base a fallback suggestion on.")
	}

	return slot, false
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
// func findNextWeekday(startAfter time.Time, targetWeekday time.Weekday) time.Time {
// 	// Start checking from the day *after* startAfter
// 	currentDate := startAfter.AddDate(0, 0, 1)
// 	for i := 0; i < 7; i++ { // Loop a maximum of 7 days to find the next occurrence
// 		if currentDate.Weekday() == targetWeekday {
// 			return currentDate
// 		}
// 		currentDate = currentDate.AddDate(0, 0, 1)
// 	}
// 	// Should realistically always find within 7 days, but fallback if needed
// 	// This indicates a potential logic error if reached.
// 	// Depending on test strictness, you might want to panic or t.Fatal here.
// 	return startAfter.AddDate(0, 0, 7) // Fallback: return a week later
// }

// // createRFC3339Timestamp creates an RFC3339 formatted timestamp string for a given date and HH:MM time.
// func createRFC3339Timestamp(t *testing.T, targetDate time.Time, timeStr string, loc *time.Location) string {
// 	t.Helper()
// 	layout := "15:04"                                             // HH:MM format
// 	parsedTime, err := time.ParseInLocation(layout, timeStr, loc) // Use ParseInLocation
// 	if err != nil {
// 		t.Fatalf("Failed to parse time string '%s': %v", timeStr, err)
// 	}

// 	// Combine the target date with the parsed time components
// 	appointmentDateTime := time.Date(
// 		targetDate.Year(),
// 		targetDate.Month(),
// 		targetDate.Day(),
// 		parsedTime.Hour(),
// 		parsedTime.Minute(),
// 		0, 0, // Seconds and Nanoseconds
// 		loc, // Use the specified location
// 	)
// 	return appointmentDateTime.Format(time.RFC3339) // Format to RFC3339
// }

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

func employeeIsAssignedToBranch(e *Employee, branchID uuid.UUID) bool {
	for _, b := range e.branches {
		if b.created.ID == branchID {
			return true
		}
	}
	return false
}
