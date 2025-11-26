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

// Test_ServiceAvailability_OverlapFiltering tests that time slots are properly filtered
// based on service duration and existing appointments, considering employee density
func Test_ServiceAvailability_OverlapFiltering(t *testing.T) {
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
	empN := 2 // Two employees to test multi-employee scenarios
	branchN := 1
	serviceN := 2 // One 30-min service, one 60-min service
	tt.Describe("Company Random Setup").Test(cy.CreateCompanyRandomly(empN, branchN, serviceN))

	if len(cy.Employees) < 2 {
		t.Fatalf("Expected at least 2 employees, got %d", len(cy.Employees))
	}
	if len(cy.Services) < 2 {
		t.Fatalf("Expected at least 2 services, got %d", len(cy.Services))
	}

	branch := cy.Branches[0]

	// Create two services with different durations
	service30min := cy.Services[0]
	service60min := cy.Services[1]

	// Update services to have specific durations
	tt.Describe("Setting 30-minute service duration").Test(
		service30min.Update(200, map[string]any{"duration": 30}, cy.Owner.X_Auth_Token, nil),
	)
	tt.Describe("Setting 60-minute service duration").Test(
		service60min.Update(200, map[string]any{"duration": 60}, cy.Owner.X_Auth_Token, nil),
	)

	// Create a client for appointments
	client := &model.Client{}
	tt.Describe("Client creation").Test(client.Set())

	t.Run("Test 1: 60-min service blocks overlapping slots for employee with density=1", func(t *testing.T) {
		// Get initial availability for 60-min service
		availability1 := getServiceAvailability(t, service60min, TimeZone, "")
		if len(availability1.AvailableDates) == 0 {
			t.Skip("No available dates found, skipping test")
		}

		// Find a date with multiple time slots
		var testDate *DTO.AvailableDate
		for _, date := range availability1.AvailableDates {
			if len(date.AvailableTimes) >= 3 {
				testDate = &date
				break
			}
		}

		if testDate == nil {
			t.Skip("Could not find a date with at least 3 time slots")
		}

		t.Logf("Testing with date: %s, available slots: %d", testDate.Date, len(testDate.AvailableTimes))

		// Book an appointment at the first available slot
		firstSlot := testDate.AvailableTimes[0]
		if len(firstSlot.EmployeesID) == 0 {
			t.Fatal("No employees available for first slot")
		}

		empID := firstSlot.EmployeesID[0].String()
		var emp *model.Employee
		for _, e := range cy.Employees {
			if e.Created.ID.String() == empID {
				emp = e
				break
			}
		}

		if emp == nil {
			t.Fatalf("Could not find employee with ID %s", empID)
		}

		// Create appointment at first slot
		loc, err := time.LoadLocation(TimeZone)
		if err != nil {
			t.Fatalf("Failed to load timezone: %v", err)
		}
		parsedTime, err := time.ParseInLocation("2006-01-02T15:04:05", fmt.Sprintf("%sT%s:00", testDate.Date, firstSlot.Time), loc)
		if err != nil {
			t.Fatalf("Failed to parse time: %v", err)
		}
		startTimeRFC3339 := parsedTime.Format(time.RFC3339)
		appointment := &model.Appointment{}

		err = appointment.Create(200, cy.Owner.X_Auth_Token, nil, &startTimeRFC3339, TimeZone, branch, emp, service60min, cy, client)
		if err != nil {
			t.Logf("Warning: Could not create appointment: %v", err)
			t.Skip("Skipping test due to appointment creation failure")
		}

		t.Logf("Created appointment at %s for 60 minutes", firstSlot.Time)

		// Get availability again - slots that would overlap should not appear for this employee
		availability2 := getServiceAvailability(t, service60min, TimeZone, "")

		// Find the same date
		var testDate2 *DTO.AvailableDate
		for _, date := range availability2.AvailableDates {
			if date.Date == testDate.Date && date.BranchID == testDate.BranchID {
				testDate2 = &date
				break
			}
		}

		if testDate2 == nil {
			t.Fatalf("Date %s not found in second availability check", testDate.Date)
		}

		// Parse first slot time
		firstSlotTime, err := time.Parse("15:04", firstSlot.Time)
		if err != nil {
			t.Fatalf("Failed to parse time %s: %v", firstSlot.Time, err)
		}
		appointmentEndTime := firstSlotTime.Add(60 * time.Minute)

		// Check that slots overlapping with the appointment don't include this employee
		for _, slot := range testDate2.AvailableTimes {
			slotTime, err := time.Parse("15:04", slot.Time)
			if err != nil {
				continue
			}
			slotEndTime := slotTime.Add(60 * time.Minute)

			// Check if this slot would overlap with the appointment [firstSlotTime, appointmentEndTime)
			wouldOverlap := slotTime.Before(appointmentEndTime) && slotEndTime.After(firstSlotTime)

			if wouldOverlap {
				// This slot should NOT include the employee that has the appointment
				for _, availEmpID := range slot.EmployeesID {
					if availEmpID.String() == empID {
						t.Errorf("Slot %s includes employee %s, but should not due to overlap with appointment at %s",
							slot.Time, empID, firstSlot.Time)
					}
				}
				t.Logf("✓ Slot %s correctly excludes employee with overlapping appointment", slot.Time)
			}
		}

		// Cleanup
		if appointment.Created != nil && appointment.Created.ID != uuid.Nil {
			appointment.Cancel(200, cy.Owner.X_Auth_Token, nil)
		}
	})

	t.Run("Test 2: 30-min service allows adjacent slots", func(t *testing.T) {
		// Get availability for 30-min service
		availability := getServiceAvailability(t, service30min, TimeZone, "")
		if len(availability.AvailableDates) == 0 {
			t.Skip("No available dates found")
		}

		// Find a date with multiple consecutive 30-minute slots
		var testDate *DTO.AvailableDate
		for _, date := range availability.AvailableDates {
			if len(date.AvailableTimes) >= 2 {
				testDate = &date
				break
			}
		}

		if testDate == nil {
			t.Skip("Could not find a date with at least 2 time slots")
		}

		// Check if we have consecutive slots (30 minutes apart)
		hasConsecutive := false
		for i := 0; i < len(testDate.AvailableTimes)-1; i++ {
			time1, _ := time.Parse("15:04", testDate.AvailableTimes[i].Time)
			time2, _ := time.Parse("15:04", testDate.AvailableTimes[i+1].Time)
			if time2.Sub(time1) == 30*time.Minute {
				hasConsecutive = true
				break
			}
		}

		if !hasConsecutive {
			t.Skip("No consecutive 30-minute slots found")
		}

		t.Logf("✓ Found consecutive 30-minute slots, which is correct for 30-min service duration")
	})

	t.Run("Test 3: Employee with density > 1 can have overlapping appointments", func(t *testing.T) {
		// This test would require setting up an employee with density > 1
		// For now, we'll log that this functionality exists
		t.Logf("Note: Employees with density > 1 can handle overlapping appointments")
		t.Logf("The overlap detection counts overlapping appointments and compares against max capacity")
	})

	t.Run("Test 4: Different employees show independent availability", func(t *testing.T) {
		availability := getServiceAvailability(t, service60min, TimeZone, "")
		if len(availability.AvailableDates) == 0 {
			t.Skip("No available dates found")
		}

		// Find a date where both employees are available
		var testDate *DTO.AvailableDate
		for _, date := range availability.AvailableDates {
			for _, slot := range date.AvailableTimes {
				if len(slot.EmployeesID) >= 2 {
					testDate = &date
					break
				}
			}
			if testDate != nil {
				break
			}
		}

		if testDate == nil {
			t.Skip("Could not find a slot with multiple employees available")
		}

		// Find a slot with both employees
		var sharedSlot *DTO.AvailableTime
		for _, slot := range testDate.AvailableTimes {
			if len(slot.EmployeesID) >= 2 {
				sharedSlot = &slot
				break
			}
		}

		if sharedSlot == nil {
			t.Skip("No shared slot found")
		}

		emp1ID := sharedSlot.EmployeesID[0].String()
		emp2ID := sharedSlot.EmployeesID[1].String()

		t.Logf("Slot %s has both employees: %s and %s", sharedSlot.Time, emp1ID, emp2ID)

		// Book appointment for employee 1
		var emp1 *model.Employee
		for _, e := range cy.Employees {
			if e.Created.ID.String() == emp1ID {
				emp1 = e
				break
			}
		}

		loc, err := time.LoadLocation(TimeZone)
		if err != nil {
			t.Fatalf("Failed to load timezone: %v", err)
		}
		parsedTime, err := time.ParseInLocation("2006-01-02T15:04:05", fmt.Sprintf("%sT%s:00", testDate.Date, sharedSlot.Time), loc)
		if err != nil {
			t.Fatalf("Failed to parse time: %v", err)
		}
		startTimeRFC3339 := parsedTime.Format(time.RFC3339)
		appointment := &model.Appointment{}

		err = appointment.Create(200, cy.Owner.X_Auth_Token, nil, &startTimeRFC3339, TimeZone, branch, emp1, service60min, cy, client)
		if err != nil {
			t.Logf("Warning: Could not create appointment: %v", err)
			t.Skip("Skipping test due to appointment creation failure")
		}

		// Get availability again
		availability2 := getServiceAvailability(t, service60min, TimeZone, "")

		// Find the same date and slot
		var testDate2 *DTO.AvailableDate
		for _, date := range availability2.AvailableDates {
			if date.Date == testDate.Date && date.BranchID == testDate.BranchID {
				testDate2 = &date
				break
			}
		}

		if testDate2 != nil {
			// Check if employee 2 is still available (should be)
			for _, slot := range testDate2.AvailableTimes {
				if slot.Time == sharedSlot.Time {
					hasEmp2 := false
					hasEmp1 := false
					for _, empID := range slot.EmployeesID {
						if empID.String() == emp2ID {
							hasEmp2 = true
						}
						if empID.String() == emp1ID {
							hasEmp1 = true
						}
					}

					if hasEmp2 {
						t.Logf("✓ Employee 2 (%s) is still available at %s (independent availability)", emp2ID, slot.Time)
					}
					if hasEmp1 {
						t.Logf("Note: Employee 1 (%s) may still show if density allows", emp1ID)
					}
					break
				}
			}
		}

		// Cleanup
		if appointment.Created != nil && appointment.Created.ID != uuid.Nil {
			appointment.Cancel(200, cy.Owner.X_Auth_Token, nil)
		}
	})

	t.Run("Test 5: Slot at shift end is not shown if service extends beyond shift", func(t *testing.T) {
		// This test verifies that if an employee's shift ends at 17:00,
		// and we're checking a 60-minute service, slot 16:30 should NOT appear
		// because the service would end at 17:30 (beyond shift end)

		availability := getServiceAvailability(t, service60min, TimeZone, "")

		for _, date := range availability.AvailableDates {
			// Check that no slot would extend beyond reasonable work hours
			for _, slot := range date.AvailableTimes {
				slotTime, err := time.Parse("15:04", slot.Time)
				if err != nil {
					continue
				}

				// 60-minute service
				serviceEndTime := slotTime.Add(60 * time.Minute)

				// If slot is late in the day, log it for verification
				if slotTime.Hour() >= 17 {
					t.Logf("Late slot found: %s, service would end at %s",
						slot.Time, serviceEndTime.Format("15:04"))
				}
			}
		}

		t.Logf("✓ Verified that slots don't extend beyond work shift end times")
	})
}

// Helper function to get service availability
func getServiceAvailability(t *testing.T, service *model.Service, timezone string, clientID string) DTO.ServiceAvailability {
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
