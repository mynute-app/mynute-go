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
)

// Test_BranchAppointments_Filters tests the employee_id and service_id query parameters
// for the GET /branch/{id}/appointments endpoint
func Test_BranchAppointments_Filters(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()
	tt := handler.NewTestErrorHandler(t)

	appEnv := os.Getenv("APP_ENV")
	if appEnv != "test" {
		t.Fatal("APP_ENV is not set to 'test'. Aborting tests to prevent data loss.")
	}

	TimeZone := "America/Sao_Paulo"

	// Setup: Create company with 2 employees, 1 branch, and 2 services
	cy := &model.Company{}
	empN := 2
	branchN := 1
	serviceN := 2
	tt.Describe("Company Random Setup").Test(cy.CreateCompanyRandomly(empN, branchN, serviceN))

	if len(cy.Employees) < 2 {
		t.Fatalf("Expected at least 2 employees, got %d", len(cy.Employees))
	}
	if len(cy.Services) < 2 {
		t.Fatalf("Expected at least 2 services, got %d", len(cy.Services))
	}

	branch := cy.Branches[0]
	service1 := cy.Services[0]
	service2 := cy.Services[1]
	employee1 := cy.Employees[0]
	employee2 := cy.Employees[1]

	// Update services to have 60-minute duration
	tt.Describe("Update service1 duration").Test(service1.Update(200, map[string]any{
		"duration": 60,
	}, cy.Owner.X_Auth_Token, nil))

	tt.Describe("Update service2 duration").Test(service2.Update(200, map[string]any{
		"duration": 60,
	}, cy.Owner.X_Auth_Token, nil))

	// Ensure all services are assigned to all entities (to guarantee common services for work schedules)
	// Silently ignore "already has service" errors (400 status)
	_ = branch.AddService(200, service1, cy.Owner.X_Auth_Token, nil) // Ignore if already assigned
	_ = branch.AddService(200, service2, cy.Owner.X_Auth_Token, nil) // Ignore if already assigned

	_ = employee1.AddService(200, service1, nil, nil) // Ignore if already assigned
	_ = employee1.AddService(200, service2, nil, nil) // Ignore if already assigned
	_ = employee2.AddService(200, service1, nil, nil) // Ignore if already assigned
	_ = employee2.AddService(200, service2, nil, nil) // Ignore if already assigned

	// Create 3 different clients
	client1 := &model.Client{}
	tt.Describe("Create client1").Test(client1.Set())

	client2 := &model.Client{}
	tt.Describe("Create client2").Test(client2.Set())

	client3 := &model.Client{}
	tt.Describe("Create client3").Test(client3.Set())

	loc, err := time.LoadLocation(TimeZone)
	if err != nil {
		t.Fatalf("Failed to load timezone: %v", err)
	}

	now := time.Now().In(loc)
	tomorrow := now.AddDate(0, 0, 1)
	tomorrowStr := tomorrow.Format("2006-01-02")

	// Create appointments with different combinations:
	// 1. Employee1, Service1, Client1, Tomorrow 10:00
	// 2. Employee1, Service2, Client2, Tomorrow 11:00
	// 3. Employee2, Service1, Client3, Tomorrow 10:00
	// 4. Employee2, Service2, Client1, Tomorrow 14:00

	appointments := []struct {
		Employee *model.Employee
		Service  *model.Service
		Client   *model.Client
		Time     string
		Name     string
	}{
		{employee1, service1, client1, "10:00", "Appt1: Emp1-Svc1-Client1"},
		{employee1, service2, client2, "11:00", "Appt2: Emp1-Svc2-Client2"},
		{employee2, service1, client3, "10:00", "Appt3: Emp2-Svc1-Client3"},
		{employee2, service2, client1, "14:00", "Appt4: Emp2-Svc2-Client1"},
	}

	for _, apt := range appointments {
		slotDateTime, err := time.ParseInLocation("2006-01-02T15:04:05", fmt.Sprintf("%sT%s:00", tomorrowStr, apt.Time), loc)
		if err != nil {
			t.Fatalf("Failed to parse time for %s: %v", apt.Name, err)
		}
		startTimeStr := slotDateTime.Format(time.RFC3339)

		appointment := &model.Appointment{}
		err = appointment.Create(200, cy.Owner.X_Auth_Token, nil, &startTimeStr, TimeZone, branch, apt.Employee, apt.Service, cy, apt.Client)
		if err != nil {
			t.Fatalf("Could not create %s: %v", apt.Name, err)
		}
		t.Logf("✓ Created %s at %s", apt.Name, apt.Time)
	}

	tomorrowFormatted := tomorrow.Format("02/01/2006")

	// Test 1: Get all appointments for the branch (no filters)
	t.Run("Test 1: Get all branch appointments without filters", func(t *testing.T) {
		http := handler.NewHttpClient()
		http.Method("GET")
		http.ExpectedStatus(200)
		query := fmt.Sprintf("page=1&page_size=10&timezone=%s&start_date=%s&end_date=%s", TimeZone, tomorrowFormatted, tomorrowFormatted)
		url := fmt.Sprintf("/branch/%s/appointments?%s", branch.Created.ID.String(), query)
		http.URL(url)
		http.Header("X-Company-ID", cy.Created.ID.String())
		http.Header("X-Auth-Token", cy.Owner.X_Auth_Token)
		http.Send(nil)

		if http.Error != nil {
			t.Fatalf("Failed to get branch appointments: %v", http.Error)
		}

		var appointmentList DTO.AppointmentList
		http.ParseResponse(&appointmentList)

		if appointmentList.TotalCount != 4 {
			t.Errorf("Expected 4 total appointments, got %d", appointmentList.TotalCount)
		}

		t.Logf("✓ Found %d appointments without filters", appointmentList.TotalCount)
	})

	// Test 2: Filter by employee_id (employee1) - should return 2 appointments
	t.Run("Test 2: Filter by employee_id (employee1)", func(t *testing.T) {
		http := handler.NewHttpClient()
		http.Method("GET")
		http.ExpectedStatus(200)
		query := fmt.Sprintf("page=1&page_size=10&timezone=%s&start_date=%s&end_date=%s&employee_id=%s",
			TimeZone, tomorrowFormatted, tomorrowFormatted, employee1.Created.ID.String())
		url := fmt.Sprintf("/branch/%s/appointments?%s", branch.Created.ID.String(), query)
		http.URL(url)
		http.Header("X-Company-ID", cy.Created.ID.String())
		http.Header("X-Auth-Token", cy.Owner.X_Auth_Token)
		http.Send(nil)

		if http.Error != nil {
			t.Fatalf("Failed to get filtered appointments: %v", http.Error)
		}

		var appointmentList DTO.AppointmentList
		http.ParseResponse(&appointmentList)

		if appointmentList.TotalCount != 2 {
			t.Errorf("Expected 2 appointments for employee1, got %d", appointmentList.TotalCount)
		}

		// Verify all appointments belong to employee1
		for _, apt := range appointmentList.Appointments {
			if apt.EmployeeID.String() != employee1.Created.ID.String() {
				t.Errorf("Expected appointment to belong to employee1, got employee_id: %s", apt.EmployeeID.String())
			}
		}

		t.Logf("✓ Found %d appointments for employee1", appointmentList.TotalCount)
	})

	// Test 3: Filter by service_id (service1) - should return 2 appointments
	t.Run("Test 3: Filter by service_id (service1)", func(t *testing.T) {
		http := handler.NewHttpClient()
		http.Method("GET")
		http.ExpectedStatus(200)
		query := fmt.Sprintf("page=1&page_size=10&timezone=%s&start_date=%s&end_date=%s&service_id=%s",
			TimeZone, tomorrowFormatted, tomorrowFormatted, service1.Created.ID.String())
		url := fmt.Sprintf("/branch/%s/appointments?%s", branch.Created.ID.String(), query)
		http.URL(url)
		http.Header("X-Company-ID", cy.Created.ID.String())
		http.Header("X-Auth-Token", cy.Owner.X_Auth_Token)
		http.Send(nil)

		if http.Error != nil {
			t.Fatalf("Failed to get filtered appointments: %v", http.Error)
		}

		var appointmentList DTO.AppointmentList
		http.ParseResponse(&appointmentList)

		if appointmentList.TotalCount != 2 {
			t.Errorf("Expected 2 appointments for service1, got %d", appointmentList.TotalCount)
		}

		// Verify all appointments belong to service1
		for _, apt := range appointmentList.Appointments {
			if apt.ServiceID.String() != service1.Created.ID.String() {
				t.Errorf("Expected appointment to belong to service1, got service_id: %s", apt.ServiceID.String())
			}
		}

		t.Logf("✓ Found %d appointments for service1", appointmentList.TotalCount)
	})

	// Test 4: Filter by both employee_id and service_id - should return 1 appointment
	t.Run("Test 4: Filter by both employee_id (employee1) and service_id (service1)", func(t *testing.T) {
		http := handler.NewHttpClient()
		http.Method("GET")
		http.ExpectedStatus(200)
		query := fmt.Sprintf("page=1&page_size=10&timezone=%s&start_date=%s&end_date=%s&employee_id=%s&service_id=%s",
			TimeZone, tomorrowFormatted, tomorrowFormatted, employee1.Created.ID.String(), service1.Created.ID.String())
		url := fmt.Sprintf("/branch/%s/appointments?%s", branch.Created.ID.String(), query)
		http.URL(url)
		http.Header("X-Company-ID", cy.Created.ID.String())
		http.Header("X-Auth-Token", cy.Owner.X_Auth_Token)
		http.Send(nil)

		if http.Error != nil {
			t.Fatalf("Failed to get filtered appointments: %v", http.Error)
		}

		var appointmentList DTO.AppointmentList
		http.ParseResponse(&appointmentList)

		if appointmentList.TotalCount != 1 {
			t.Errorf("Expected 1 appointment for employee1+service1, got %d", appointmentList.TotalCount)
		}

		// Verify the appointment matches both filters
		if len(appointmentList.Appointments) > 0 {
			apt := appointmentList.Appointments[0]
			if apt.EmployeeID.String() != employee1.Created.ID.String() {
				t.Errorf("Expected appointment to belong to employee1")
			}
			if apt.ServiceID.String() != service1.Created.ID.String() {
				t.Errorf("Expected appointment to belong to service1")
			}
		}

		t.Logf("✓ Found %d appointment for employee1+service1 combination", appointmentList.TotalCount)
	})

	// Test 5: Filter with non-matching combination - should return 0 appointments
	t.Run("Test 5: Filter by employee_id (employee2) and service_id (service2) with date mismatch", func(t *testing.T) {
		// Query for yesterday (no appointments)
		yesterday := now.AddDate(0, 0, -1).Format("02/01/2006")

		http := handler.NewHttpClient()
		http.Method("GET")
		http.ExpectedStatus(200)
		query := fmt.Sprintf("page=1&page_size=10&timezone=%s&start_date=%s&end_date=%s&employee_id=%s&service_id=%s",
			TimeZone, yesterday, yesterday, employee2.Created.ID.String(), service2.Created.ID.String())
		url := fmt.Sprintf("/branch/%s/appointments?%s", branch.Created.ID.String(), query)
		http.URL(url)
		http.Header("X-Company-ID", cy.Created.ID.String())
		http.Header("X-Auth-Token", cy.Owner.X_Auth_Token)
		http.Send(nil)

		if http.Error != nil {
			t.Fatalf("Failed to get filtered appointments: %v", http.Error)
		}

		var appointmentList DTO.AppointmentList
		http.ParseResponse(&appointmentList)

		if appointmentList.TotalCount != 0 {
			t.Errorf("Expected 0 appointments for yesterday, got %d", appointmentList.TotalCount)
		}

		t.Logf("✓ Correctly returned 0 appointments for non-matching date")
	})

	// Test 6: Invalid employee_id format - should return 400
	t.Run("Test 6: Invalid employee_id format", func(t *testing.T) {
		http := handler.NewHttpClient()
		http.Method("GET")
		http.ExpectedStatus(400)
		query := fmt.Sprintf("page=1&page_size=10&timezone=%s&employee_id=invalid-uuid", TimeZone)
		url := fmt.Sprintf("/branch/%s/appointments?%s", branch.Created.ID.String(), query)
		http.URL(url)
		http.Header("X-Company-ID", cy.Created.ID.String())
		http.Header("X-Auth-Token", cy.Owner.X_Auth_Token)
		http.Send(nil)

		if http.Error == nil {
			t.Error("Expected error for invalid employee_id format")
		}

		t.Logf("✓ Correctly rejected invalid employee_id format")
	})

	// Test 7: Invalid service_id format - should return 400
	t.Run("Test 7: Invalid service_id format", func(t *testing.T) {
		http := handler.NewHttpClient()
		http.Method("GET")
		http.ExpectedStatus(400)
		query := fmt.Sprintf("page=1&page_size=10&timezone=%s&service_id=not-a-uuid", TimeZone)
		url := fmt.Sprintf("/branch/%s/appointments?%s", branch.Created.ID.String(), query)
		http.URL(url)
		http.Header("X-Company-ID", cy.Created.ID.String())
		http.Header("X-Auth-Token", cy.Owner.X_Auth_Token)
		http.Send(nil)

		if http.Error == nil {
			t.Error("Expected error for invalid service_id format")
		}

		t.Logf("✓ Correctly rejected invalid service_id format")
	})

	// Test 8: Filter by employee2 - should return 2 appointments
	t.Run("Test 8: Filter by employee_id (employee2)", func(t *testing.T) {
		http := handler.NewHttpClient()
		http.Method("GET")
		http.ExpectedStatus(200)
		query := fmt.Sprintf("page=1&page_size=10&timezone=%s&start_date=%s&end_date=%s&employee_id=%s",
			TimeZone, tomorrowFormatted, tomorrowFormatted, employee2.Created.ID.String())
		url := fmt.Sprintf("/branch/%s/appointments?%s", branch.Created.ID.String(), query)
		http.URL(url)
		http.Header("X-Company-ID", cy.Created.ID.String())
		http.Header("X-Auth-Token", cy.Owner.X_Auth_Token)
		http.Send(nil)

		if http.Error != nil {
			t.Fatalf("Failed to get filtered appointments: %v", http.Error)
		}

		var appointmentList DTO.AppointmentList
		http.ParseResponse(&appointmentList)

		if appointmentList.TotalCount != 2 {
			t.Errorf("Expected 2 appointments for employee2, got %d", appointmentList.TotalCount)
		}

		t.Logf("✓ Found %d appointments for employee2", appointmentList.TotalCount)
	})

	// Test 9: Pagination with filters
	t.Run("Test 9: Pagination with employee_id filter", func(t *testing.T) {
		http := handler.NewHttpClient()
		http.Method("GET")
		http.ExpectedStatus(200)
		query := fmt.Sprintf("page=1&page_size=1&timezone=%s&start_date=%s&end_date=%s&employee_id=%s",
			TimeZone, tomorrowFormatted, tomorrowFormatted, employee1.Created.ID.String())
		url := fmt.Sprintf("/branch/%s/appointments?%s", branch.Created.ID.String(), query)
		http.URL(url)
		http.Header("X-Company-ID", cy.Created.ID.String())
		http.Header("X-Auth-Token", cy.Owner.X_Auth_Token)
		http.Send(nil)

		if http.Error != nil {
			t.Fatalf("Failed to get filtered appointments: %v", http.Error)
		}

		var appointmentList DTO.AppointmentList
		http.ParseResponse(&appointmentList)

		if appointmentList.TotalCount != 2 {
			t.Errorf("Expected total count of 2, got %d", appointmentList.TotalCount)
		}

		if len(appointmentList.Appointments) != 1 {
			t.Errorf("Expected 1 appointment on page 1, got %d", len(appointmentList.Appointments))
		}

		if appointmentList.Page != 1 {
			t.Errorf("Expected page 1, got %d", appointmentList.Page)
		}

		if appointmentList.PageSize != 1 {
			t.Errorf("Expected page size 1, got %d", appointmentList.PageSize)
		}

		t.Logf("✓ Pagination works correctly with filters (page 1/2, showing 1 of 2)")
	})
}
