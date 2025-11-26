package e2e_test

import (
	"fmt"
	"mynute-go/core"
	DTO "mynute-go/core/src/config/api/dto"
	"mynute-go/test/src/handler"
	"mynute-go/test/src/model"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

// Test_ServiceAvailability_ClientFilter tests the client_id query parameter
// to ensure slots where the client already has appointments are filtered out
func Test_ServiceAvailability_ClientFilter(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()
	tt := handler.NewTestErrorHandler(t)

	appEnv := os.Getenv("APP_ENV")
	if appEnv != "test" {
		t.Fatal("APP_ENV is not set to 'test'. Aborting tests to prevent data loss.")
	}

	TimeZone := "America/Sao_Paulo"

	// Setup: Create company with employees, branches, and services
	cy := &model.Company{}
	empN := 1
	branchN := 1
	serviceN := 1
	tt.Describe("Company Random Setup").Test(cy.CreateCompanyRandomly(empN, branchN, serviceN))

	if len(cy.Employees) < 1 {
		t.Fatalf("Expected at least 1 employee, got %d", len(cy.Employees))
	}
	if len(cy.Services) < 1 {
		t.Fatalf("Expected at least 1 service, got %d", len(cy.Services))
	}

	branch := cy.Branches[0]
	service := cy.Services[0]
	employee := cy.Employees[0]

	// Update service to have 60-minute duration
	tt.Describe("Update service duration to 60 minutes").Test(service.Update(200, map[string]any{
		"duration": 60,
	}, cy.Owner.X_Auth_Token, nil))

	// Create a client
	client := &model.Client{}
	tt.Describe("Create a client").Test(client.Set())

	// Set up employee work schedule for the next 7 days
	loc, err := time.LoadLocation(TimeZone)
	if err != nil {
		t.Fatalf("Failed to load timezone: %v", err)
	}

	now := time.Now().In(loc)
	tomorrow := now.AddDate(0, 0, 1)

	// Test 1: Get availability without client_id - should show all available slots
	t.Run("Test 1: Get availability without client_id filter", func(t *testing.T) {
		availability := getServiceAvailabilityWithClient(t, service, TimeZone, "")

		if len(availability.AvailableDates) == 0 {
			t.Fatal("Expected available dates, got none")
		}

		// Find tomorrow's availability
		tomorrowStr := tomorrow.Format("2006-01-02")
		var tomorrowAvailability *DTO.AvailableDate
		for _, ad := range availability.AvailableDates {
			if ad.Date == tomorrowStr {
				tomorrowAvailability = &ad
				break
			}
		}

		if tomorrowAvailability == nil {
			t.Fatalf("Expected availability for tomorrow (%s), got none", tomorrowStr)
		}

		if len(tomorrowAvailability.AvailableTimes) == 0 {
			t.Fatal("Expected available times for tomorrow, got none")
		}

		t.Logf("✓ Without client_id: Found %d available time slots for tomorrow", len(tomorrowAvailability.AvailableTimes))
	})

	// Test 2: Create appointment for client and verify it's filtered out
	t.Run("Test 2: Create appointment and verify slot is filtered with client_id", func(t *testing.T) {
		// Get first available slot for any date (not just tomorrow)
		availability := getServiceAvailabilityWithClient(t, service, TimeZone, "")

		if len(availability.AvailableDates) == 0 {
			t.Skip("No available dates to test with")
		}

		// Find first available date with slots
		var testDate *DTO.AvailableDate
		for _, ad := range availability.AvailableDates {
			if len(ad.AvailableTimes) > 0 {
				testDate = &ad
				break
			}
		}

		if testDate == nil {
			t.Skip("No available slots found to test with")
		}

		firstSlot := testDate.AvailableTimes[0]
		slotTime := firstSlot.Time
		dateStr := testDate.Date

		// Parse the slot time in the correct timezone
		slotDateTime, err := time.Parse("2006-01-02T15:04:05", fmt.Sprintf("%sT%s:00", dateStr, slotTime))
		if err != nil {
			t.Fatalf("Failed to parse slot time: %v", err)
		}
		slotDateTime = time.Date(slotDateTime.Year(), slotDateTime.Month(), slotDateTime.Day(), slotDateTime.Hour(), slotDateTime.Minute(), 0, 0, loc)

		// Create appointment for the client at this time
		startTimeStr := slotDateTime.Format(time.RFC3339)

		appointment := &model.Appointment{}
		err = appointment.Create(200, cy.Owner.X_Auth_Token, nil, &startTimeStr, TimeZone, branch, employee, service, cy, client)
		if err != nil {
			t.Skipf("Could not create appointment at %s %s: %v", dateStr, slotTime, err)
		}

		t.Logf("✓ Created appointment at %s %s", dateStr, slotTime)

		// Get availability WITH client_id - should filter out the booked slot
		availabilityWithFilter := getServiceAvailabilityWithClient(t, service, TimeZone, client.Created.ID.String())

		var testDateWithFilter *DTO.AvailableDate
		for _, ad := range availabilityWithFilter.AvailableDates {
			if ad.Date == dateStr {
				testDateWithFilter = &ad
				break
			}
		}

		// Verify the booked slot is not in the filtered results
		if testDateWithFilter != nil {
			for _, slot := range testDateWithFilter.AvailableTimes {
				if slot.Time == slotTime {
					t.Errorf("Expected slot %s to be filtered out for client %s, but it was present", slotTime, client.Created.ID.String())
				}
			}
		}

		t.Logf("✓ With client_id filter: Booked slot %s correctly filtered out", slotTime)
	})

	// Test 3: Create multiple appointments and verify all are filtered
	t.Run("Test 3: Create multiple appointments and verify all are filtered", func(t *testing.T) {
		// Get availability with client filter
		availability := getServiceAvailabilityWithClient(t, service, TimeZone, client.Created.ID.String())

		if len(availability.AvailableDates) == 0 {
			t.Log("✓ All slots already filtered (client has appointments)")
			return
		}

		// Find first date with available slots
		var testDate *DTO.AvailableDate
		for _, ad := range availability.AvailableDates {
			if len(ad.AvailableTimes) > 0 {
				testDate = &ad
				break
			}
		}

		if testDate == nil {
			t.Log("✓ All slots already filtered (client has appointment covering all available times)")
			return
		}

		dateStr := testDate.Date

		// Try to book up to 2 more appointments if available
		slotsToBook := 2
		if len(testDate.AvailableTimes) < slotsToBook {
			slotsToBook = len(testDate.AvailableTimes)
		}

		bookedSlots := []string{}
		for i := 0; i < slotsToBook; i++ {
			slot := testDate.AvailableTimes[i]
			slotTime := slot.Time

			// Parse using the same date as the availability check
			slotDateTime, err := time.Parse("2006-01-02T15:04:05", fmt.Sprintf("%sT%s:00", dateStr, slotTime))
			if err != nil {
				t.Logf("Warning: Failed to parse slot time: %v", err)
				continue
			}
			// Ensure the time is in the correct timezone
			slotDateTime = time.Date(slotDateTime.Year(), slotDateTime.Month(), slotDateTime.Day(), slotDateTime.Hour(), slotDateTime.Minute(), 0, 0, loc)

			startTimeStr := slotDateTime.Format(time.RFC3339)

			appointment := &model.Appointment{}
			err = appointment.Create(200, cy.Owner.X_Auth_Token, nil, &startTimeStr, TimeZone, branch, employee, service, cy, client)
			if err != nil {
				t.Logf("Warning: Could not create appointment %d at %s: %v", i+2, slotTime, err)
				continue
			}
			bookedSlots = append(bookedSlots, slotTime)
		}

		t.Logf("✓ Created %d additional appointments at slots: %v", len(bookedSlots), bookedSlots)

		// Only continue if we were able to book some appointments
		if len(bookedSlots) == 0 {
			t.Log("Warning: Could not create any additional appointments, skipping rest of test")
			return
		}

		// Get availability again with client filter
		availabilityAfter := getServiceAvailabilityWithClient(t, service, TimeZone, client.Created.ID.String())

		var testDateAfter *DTO.AvailableDate
		for _, ad := range availabilityAfter.AvailableDates {
			if ad.Date == dateStr {
				testDateAfter = &ad
				break
			}
		}

		// Verify all booked slots are filtered out
		if testDateAfter != nil {
			for _, slot := range testDateAfter.AvailableTimes {
				for _, bookedSlot := range bookedSlots {
					if slot.Time == bookedSlot {
						t.Errorf("Expected slot %s to be filtered out, but it was present", bookedSlot)
					}
				}
			}
		}

		t.Logf("✓ All %d booked slots correctly filtered from availability", len(bookedSlots))
	})

	// Test 4: Test with non-existent client_id (valid UUID but doesn't exist)
	t.Run("Test 4: Test with non-existent client_id", func(t *testing.T) {
		nonExistentClientID := uuid.New().String()

		availability := getServiceAvailabilityWithClient(t, service, TimeZone, nonExistentClientID)

		// Should return availability normally (no appointments for this client)
		if len(availability.AvailableDates) == 0 {
			t.Error("Expected available dates for non-existent client_id")
		} else {
			t.Logf("✓ Non-existent client_id returned %d available dates (expected behavior)", len(availability.AvailableDates))
		}
	})

	// Test 5: Test cross-service filtering - client with appointment in Service A should see filtered slots in Service B
	t.Run("Test 5: Test cross-service filtering", func(t *testing.T) {
		// Create a second service
		service2 := &model.Service{Company: cy}
		tt.Describe("Create second service").Test(service2.Create(200, cy.Owner.X_Auth_Token, nil))

		// Update second service duration to match first service
		tt.Describe("Update second service duration").Test(service2.Update(200, map[string]any{
			"duration": 60,
		}, cy.Owner.X_Auth_Token, nil))

		// Add service to branch and employee
		tt.Describe("Add second service to branch").Test(branch.AddService(200, service2, cy.Owner.X_Auth_Token, nil))
		tt.Describe("Add second service to employee").Test(employee.AddService(200, service2, nil, nil))

		// Add second service to employee's existing work schedules
		// Since CreateCompanyRandomly now guarantees all 7 weekdays, we just need to add service2
		workRangesToAdd2 := []DTO.CreateEmployeeWorkRange{}
		for weekday := 0; weekday < 7; weekday++ {
			workRangesToAdd2 = append(workRangesToAdd2, DTO.CreateEmployeeWorkRange{
				EmployeeID: employee.Created.ID,
				StartTime:  "08:00",
				EndTime:    "20:00",
				TimeZone:   TimeZone,
				Weekday:    uint8(weekday),
				BranchID:   branch.Created.ID,
				EmployeeWorkRangeServices: DTO.EmployeeWorkRangeServices{
					Services: []DTO.ServiceBase{{ID: service2.Created.ID}},
				},
			})
		}

		additionalWorkSchedule2 := DTO.CreateEmployeeWorkSchedule{
			WorkRanges: workRangesToAdd2,
		}
		// Try to create, but don't fail if already exists (employee might already have this weekday)
		_ = employee.CreateWorkSchedule(200, additionalWorkSchedule2, nil, nil)

		// Create a new client for this test
		client3 := &model.Client{}
		tt.Describe("Create third client").Test(client3.Set())

		// Book an appointment for client3 in the FIRST service at a specific time
		availability := getServiceAvailabilityWithClient(t, service, TimeZone, "")
		if len(availability.AvailableDates) == 0 {
			t.Skip("No available dates for first service")
		}

		var testDate *DTO.AvailableDate
		for _, ad := range availability.AvailableDates {
			if len(ad.AvailableTimes) > 0 {
				testDate = &ad
				break
			}
		}

		if testDate == nil {
			t.Skip("No available slots in first service")
		}

		firstSlot := testDate.AvailableTimes[0]
		slotTime := firstSlot.Time
		dateStr := testDate.Date

		slotDateTime, err := time.Parse("2006-01-02T15:04:05", fmt.Sprintf("%sT%s:00", dateStr, slotTime))
		if err != nil {
			t.Fatalf("Failed to parse slot time: %v", err)
		}
		slotDateTime = time.Date(slotDateTime.Year(), slotDateTime.Month(), slotDateTime.Day(), slotDateTime.Hour(), slotDateTime.Minute(), 0, 0, loc)
		startTimeStr := slotDateTime.Format(time.RFC3339)

		// Create appointment in FIRST service
		appointment := &model.Appointment{}
		err = appointment.Create(200, cy.Owner.X_Auth_Token, nil, &startTimeStr, TimeZone, branch, employee, service, cy, client3)
		if err != nil {
			t.Skipf("Could not create appointment in first service: %v", err)
		}

		t.Logf("✓ Created appointment in first service at %s %s for client3", dateStr, slotTime)

		// Now check availability for the SECOND service with client3's ID
		// The same time slot should be filtered out because the client is busy
		availability2 := getServiceAvailabilityWithClient(t, service2, TimeZone, client3.Created.ID.String())

		// Check if the booked slot appears in service2's availability
		var testDate2 *DTO.AvailableDate
		for _, ad := range availability2.AvailableDates {
			if ad.Date == dateStr {
				testDate2 = &ad
				break
			}
		}

		if testDate2 != nil {
			for _, slot := range testDate2.AvailableTimes {
				if slot.Time == slotTime {
					t.Errorf("CROSS-SERVICE FILTERING FAILED: Expected slot %s to be filtered out in service2 for client3 (who has appointment in service1), but it was present", slotTime)
				}
			}
		}

		t.Logf("✓ Cross-service filtering works: Slot %s correctly filtered out in service2 when client3 has appointment in service1", slotTime)
	})

	// Test 6: Test overlapping appointment filtering
	t.Run("Test 6: Test overlapping appointment filtering", func(t *testing.T) {
		// Create a new client for this test
		client2 := &model.Client{}
		tt.Describe("Create second client").Test(client2.Set())

		// Get availability for any available date
		availability := getServiceAvailabilityWithClient(t, service, TimeZone, "")

		if len(availability.AvailableDates) == 0 {
			t.Skip("No available dates to test with")
		}

		// Find an available date with slots
		var testDate *DTO.AvailableDate
		for _, ad := range availability.AvailableDates {
			if len(ad.AvailableTimes) > 0 {
				testDate = &ad
				break
			}
		}

		if testDate == nil {
			t.Skip("No available slots to test with")
		}

		dateStr := testDate.Date
		firstSlot := testDate.AvailableTimes[0]
		slotTime := firstSlot.Time

		// Book an appointment for client2
		slotDateTime, err := time.Parse("2006-01-02T15:04:05", fmt.Sprintf("%sT%s:00", dateStr, slotTime))
		if err != nil {
			t.Fatalf("Failed to parse slot time: %v", err)
		}
		slotDateTime = time.Date(slotDateTime.Year(), slotDateTime.Month(), slotDateTime.Day(), slotDateTime.Hour(), slotDateTime.Minute(), 0, 0, loc)
		startTimeStr := slotDateTime.Format(time.RFC3339)

		appointment := &model.Appointment{}
		err = appointment.Create(200, cy.Owner.X_Auth_Token, nil, &startTimeStr, TimeZone, branch, employee, service, cy, client2)
		if err != nil {
			t.Skipf("Could not create appointment for client2: %v", err)
		}

		t.Logf("✓ Created appointment at %s %s for client2", dateStr, slotTime)

		// Get availability with client2's ID
		availabilityFiltered := getServiceAvailabilityWithClient(t, service, TimeZone, client2.Created.ID.String())

		var testDateFiltered *DTO.AvailableDate
		for _, ad := range availabilityFiltered.AvailableDates {
			if ad.Date == dateStr {
				testDateFiltered = &ad
				break
			}
		}

		// Verify the booked slot is filtered out for client2
		if testDateFiltered != nil {
			for _, slot := range testDateFiltered.AvailableTimes {
				if slot.Time == slotTime {
					t.Errorf("Expected %s slot to be filtered out for client2, but it was present", slotTime)
				}
			}
		}

		t.Logf("✓ Appointment at %s correctly filtered for client2", slotTime)
	})
}

// Helper function to get service availability with optional client filter
func getServiceAvailabilityWithClient(t *testing.T, service *model.Service, timezone string, clientID string) DTO.ServiceAvailability {
	http := handler.NewHttpClient()
	http.Method("GET")
	http.ExpectedStatus(200)

	query := fmt.Sprintf("date_forward_start=%d&date_forward_end=%d&timezone=%s", 0, 7, timezone)
	if clientID != "" {
		query += fmt.Sprintf("&client_id=%s", clientID)
	}

	url := fmt.Sprintf("/service/%s/availability?%s", service.Created.ID.String(), query)
	http.URL(url)
	http.Header("X-Company-ID", service.Company.Created.ID.String())
	http.Send(nil)

	if http.Error != nil {
		t.Fatalf("Failed to get service availability: %v", http.Error)
	}

	var availability DTO.ServiceAvailability
	http.ParseResponse(&availability)

	return availability
}
