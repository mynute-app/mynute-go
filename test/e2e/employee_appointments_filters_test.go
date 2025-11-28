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

// Test_EmployeeAppointments_Filters tests the branch_id and service_id query parameters
// for the GET /employee/{id}/appointments endpoint
func Test_EmployeeAppointments_Filters(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()
	tt := handler.NewTestErrorHandler(t)

	appEnv := os.Getenv("APP_ENV")
	if appEnv != "test" {
		t.Fatal("APP_ENV is not set to 'test'. Aborting tests to prevent data loss.")
	}

	TimeZone := "America/Sao_Paulo"

	// Setup: Create company with 1 employee, 2 branches, and 2 services
	cy := &model.Company{}
	empN := 1
	branchN := 2
	serviceN := 2
	tt.Describe("Company Random Setup").Test(cy.CreateCompanyRandomly(empN, branchN, serviceN))

	if len(cy.Employees) < 1 {
		t.Fatalf("Expected at least 1 employee, got %d", len(cy.Employees))
	}
	if len(cy.Branches) < 2 {
		t.Fatalf("Expected at least 2 branches, got %d", len(cy.Branches))
	}
	if len(cy.Services) < 2 {
		t.Fatalf("Expected at least 2 services, got %d", len(cy.Services))
	}

	branch1 := cy.Branches[0]
	branch2 := cy.Branches[1]
	service1 := cy.Services[0]
	service2 := cy.Services[1]
	employee := cy.Employees[0]

	// Update services to have 60-minute duration
	tt.Describe("Update service1 duration").Test(service1.Update(200, map[string]any{
		"duration": 60,
	}, cy.Owner.X_Auth_Token, nil))

	tt.Describe("Update service2 duration").Test(service2.Update(200, map[string]any{
		"duration": 60,
	}, cy.Owner.X_Auth_Token, nil))

	// Ensure all services are assigned to all entities (to guarantee common services for work schedules)
	// Silently ignore "already has service" errors (400 status)
	_ = branch1.AddService(200, service1, cy.Owner.X_Auth_Token, nil) // Ignore if already assigned
	_ = branch1.AddService(200, service2, cy.Owner.X_Auth_Token, nil) // Ignore if already assigned
	_ = branch2.AddService(200, service1, cy.Owner.X_Auth_Token, nil) // Ignore if already assigned
	_ = branch2.AddService(200, service2, cy.Owner.X_Auth_Token, nil) // Ignore if already assigned

	_ = employee.AddService(200, service1, nil, nil) // Ignore if already assigned
	_ = employee.AddService(200, service2, nil, nil) // Ignore if already assigned

	// Refresh employee data to see which branches they're actually assigned to
	tt.Describe("Refresh employee data").Test(employee.GetById(200, nil, nil))

	// Determine which branches the employee is assigned to
	employeeBranchIDs := make(map[string]bool)
	for _, b := range employee.Created.Branches {
		employeeBranchIDs[b.ID.String()] = true
	}

	// Check if employee has work schedules at both branches
	hasBranch1 := employeeBranchIDs[branch1.Created.ID.String()]
	hasBranch2 := employeeBranchIDs[branch2.Created.ID.String()]

	if !hasBranch1 {
		t.Fatalf("Employee is not assigned to branch1, cannot proceed with test")
	}

	t.Logf("Employee assigned to branch1: %v, branch2: %v", hasBranch1, hasBranch2)

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
	// 1. Branch1, Service1, Client1, Tomorrow 10:00
	// 2. Branch1, Service2, Client2, Tomorrow 11:00
	// 3. Branch2, Service1, Client3, Tomorrow 10:00
	// 4. Branch2, Service2, Client1, Tomorrow 14:00

	appointments := []struct {
		Branch  *model.Branch
		Service *model.Service
		Client  *model.Client
		Time    string
		Name    string
	}{
		{branch1, service1, client1, "10:00", "Appt1: Branch1-Svc1-Client1"},
		{branch1, service2, client2, "11:00", "Appt2: Branch1-Svc2-Client2"},
		{branch2, service1, client3, "10:00", "Appt3: Branch2-Svc1-Client3"},
		{branch2, service2, client1, "14:00", "Appt4: Branch2-Svc2-Client1"},
	}

	for _, apt := range appointments {
		slotDateTime, err := time.ParseInLocation("2006-01-02T15:04:05", fmt.Sprintf("%sT%s:00", tomorrowStr, apt.Time), loc)
		if err != nil {
			t.Fatalf("Failed to parse time for %s: %v", apt.Name, err)
		}
		startTimeStr := slotDateTime.Format(time.RFC3339)

		appointment := &model.Appointment{}
		err = appointment.Create(200, cy.Owner.X_Auth_Token, nil, &startTimeStr, TimeZone, apt.Branch, employee, apt.Service, cy, apt.Client)
		if err != nil {
			t.Fatalf("Could not create %s: %v", apt.Name, err)
		}
		t.Logf("✓ Created %s at %s", apt.Name, apt.Time)
	}

	tomorrowFormatted := tomorrow.Format("02/01/2006")

	// Test 1: Get all appointments for the employee (no filters)
	t.Run("Test 1: Get all employee appointments without filters", func(t *testing.T) {
		http := handler.NewHttpClient()
		http.Method("GET")
		http.ExpectedStatus(200)
		query := fmt.Sprintf("page=1&page_size=10&timezone=%s&start_date=%s&end_date=%s", TimeZone, tomorrowFormatted, tomorrowFormatted)
		url := fmt.Sprintf("/employee/%s/appointments?%s", employee.Created.ID.String(), query)
		http.URL(url)
		http.Header("X-Company-ID", cy.Created.ID.String())
		http.Header("X-Auth-Token", cy.Owner.X_Auth_Token)
		http.Send(nil)

		if http.Error != nil {
			t.Fatalf("Failed to get employee appointments: %v", http.Error)
		}

		var appointmentList DTO.AppointmentList
		http.ParseResponse(&appointmentList)

		if appointmentList.TotalCount != 4 {
			t.Errorf("Expected 4 total appointments, got %d", appointmentList.TotalCount)
		}

		t.Logf("✓ Found %d appointments without filters", appointmentList.TotalCount)
	})

	// Test 2: Filter by branch_id (branch1) - should return 2 appointments
	t.Run("Test 2: Filter by branch_id (branch1)", func(t *testing.T) {
		http := handler.NewHttpClient()
		http.Method("GET")
		http.ExpectedStatus(200)
		query := fmt.Sprintf("page=1&page_size=10&timezone=%s&start_date=%s&end_date=%s&branch_id=%s",
			TimeZone, tomorrowFormatted, tomorrowFormatted, branch1.Created.ID.String())
		url := fmt.Sprintf("/employee/%s/appointments?%s", employee.Created.ID.String(), query)
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
			t.Errorf("Expected 2 appointments for branch1, got %d", appointmentList.TotalCount)
		}

		// Verify all appointments belong to branch1
		for _, apt := range appointmentList.Appointments {
			if apt.BranchID.String() != branch1.Created.ID.String() {
				t.Errorf("Expected appointment to belong to branch1, got branch_id: %s", apt.BranchID.String())
			}
		}

		t.Logf("✓ Found %d appointments for branch1", appointmentList.TotalCount)
	})

	// Test 3: Filter by service_id (service1) - should return 2 appointments
	t.Run("Test 3: Filter by service_id (service1)", func(t *testing.T) {
		http := handler.NewHttpClient()
		http.Method("GET")
		http.ExpectedStatus(200)
		query := fmt.Sprintf("page=1&page_size=10&timezone=%s&start_date=%s&end_date=%s&service_id=%s",
			TimeZone, tomorrowFormatted, tomorrowFormatted, service1.Created.ID.String())
		url := fmt.Sprintf("/employee/%s/appointments?%s", employee.Created.ID.String(), query)
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

	// Test 4: Filter by both branch_id and service_id - should return 1 appointment
	t.Run("Test 4: Filter by both branch_id (branch1) and service_id (service1)", func(t *testing.T) {
		http := handler.NewHttpClient()
		http.Method("GET")
		http.ExpectedStatus(200)
		query := fmt.Sprintf("page=1&page_size=10&timezone=%s&start_date=%s&end_date=%s&branch_id=%s&service_id=%s",
			TimeZone, tomorrowFormatted, tomorrowFormatted, branch1.Created.ID.String(), service1.Created.ID.String())
		url := fmt.Sprintf("/employee/%s/appointments?%s", employee.Created.ID.String(), query)
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
			t.Errorf("Expected 1 appointment for branch1+service1, got %d", appointmentList.TotalCount)
		}

		// Verify the appointment matches both filters
		if len(appointmentList.Appointments) > 0 {
			apt := appointmentList.Appointments[0]
			if apt.BranchID.String() != branch1.Created.ID.String() {
				t.Errorf("Expected appointment to belong to branch1")
			}
			if apt.ServiceID.String() != service1.Created.ID.String() {
				t.Errorf("Expected appointment to belong to service1")
			}
		}

		t.Logf("✓ Found %d appointment for branch1+service1 combination", appointmentList.TotalCount)
	})

	// Test 5: Filter with non-matching combination - should return 0 appointments
	t.Run("Test 5: Filter by branch_id (branch2) and service_id (service2) with date mismatch", func(t *testing.T) {
		// Query for yesterday (no appointments)
		yesterday := now.AddDate(0, 0, -1).Format("02/01/2006")

		http := handler.NewHttpClient()
		http.Method("GET")
		http.ExpectedStatus(200)
		query := fmt.Sprintf("page=1&page_size=10&timezone=%s&start_date=%s&end_date=%s&branch_id=%s&service_id=%s",
			TimeZone, yesterday, yesterday, branch2.Created.ID.String(), service2.Created.ID.String())
		url := fmt.Sprintf("/employee/%s/appointments?%s", employee.Created.ID.String(), query)
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

	// Test 6: Invalid branch_id format - should return 400
	t.Run("Test 6: Invalid branch_id format", func(t *testing.T) {
		http := handler.NewHttpClient()
		http.Method("GET")
		http.ExpectedStatus(400)
		query := fmt.Sprintf("page=1&page_size=10&timezone=%s&branch_id=invalid-uuid", TimeZone)
		url := fmt.Sprintf("/employee/%s/appointments?%s", employee.Created.ID.String(), query)
		http.URL(url)
		http.Header("X-Company-ID", cy.Created.ID.String())
		http.Header("X-Auth-Token", cy.Owner.X_Auth_Token)
		http.Send(nil)

		if http.Error == nil {
			t.Error("Expected error for invalid branch_id format")
		}

		t.Logf("✓ Correctly rejected invalid branch_id format")
	})

	// Test 7: Invalid service_id format - should return 400
	t.Run("Test 7: Invalid service_id format", func(t *testing.T) {
		http := handler.NewHttpClient()
		http.Method("GET")
		http.ExpectedStatus(400)
		query := fmt.Sprintf("page=1&page_size=10&timezone=%s&service_id=not-a-uuid", TimeZone)
		url := fmt.Sprintf("/employee/%s/appointments?%s", employee.Created.ID.String(), query)
		http.URL(url)
		http.Header("X-Company-ID", cy.Created.ID.String())
		http.Header("X-Auth-Token", cy.Owner.X_Auth_Token)
		http.Send(nil)

		if http.Error == nil {
			t.Error("Expected error for invalid service_id format")
		}

		t.Logf("✓ Correctly rejected invalid service_id format")
	})

	// Test 8: Filter by branch2 - should return 2 appointments
	t.Run("Test 8: Filter by branch_id (branch2)", func(t *testing.T) {
		http := handler.NewHttpClient()
		http.Method("GET")
		http.ExpectedStatus(200)
		query := fmt.Sprintf("page=1&page_size=10&timezone=%s&start_date=%s&end_date=%s&branch_id=%s",
			TimeZone, tomorrowFormatted, tomorrowFormatted, branch2.Created.ID.String())
		url := fmt.Sprintf("/employee/%s/appointments?%s", employee.Created.ID.String(), query)
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
			t.Errorf("Expected 2 appointments for branch2, got %d", appointmentList.TotalCount)
		}

		t.Logf("✓ Found %d appointments for branch2", appointmentList.TotalCount)
	})

	// Test 9: Pagination with filters
	t.Run("Test 9: Pagination with branch_id filter", func(t *testing.T) {
		http := handler.NewHttpClient()
		http.Method("GET")
		http.ExpectedStatus(200)
		query := fmt.Sprintf("page=1&page_size=1&timezone=%s&start_date=%s&end_date=%s&branch_id=%s",
			TimeZone, tomorrowFormatted, tomorrowFormatted, branch1.Created.ID.String())
		url := fmt.Sprintf("/employee/%s/appointments?%s", employee.Created.ID.String(), query)
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

	// Test 10: Combine all filters (branch + service + date + cancelled)
	t.Run("Test 10: Combine branch_id, service_id, date, and cancelled filters", func(t *testing.T) {
		http := handler.NewHttpClient()
		http.Method("GET")
		http.ExpectedStatus(200)
		query := fmt.Sprintf("page=1&page_size=10&timezone=%s&start_date=%s&end_date=%s&branch_id=%s&service_id=%s&cancelled=false",
			TimeZone, tomorrowFormatted, tomorrowFormatted, branch2.Created.ID.String(), service2.Created.ID.String())
		url := fmt.Sprintf("/employee/%s/appointments?%s", employee.Created.ID.String(), query)
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
			t.Errorf("Expected 1 appointment with all filters, got %d", appointmentList.TotalCount)
		}

		if len(appointmentList.Appointments) > 0 {
			apt := appointmentList.Appointments[0]
			if apt.BranchID.String() != branch2.Created.ID.String() {
				t.Errorf("Expected branch2")
			}
			if apt.ServiceID.String() != service2.Created.ID.String() {
				t.Errorf("Expected service2")
			}
			if apt.IsCancelled {
				t.Errorf("Expected non-cancelled appointment")
			}
		}

		t.Logf("✓ All filters combined correctly")
	})

	// Test 11: Verify ServiceInfo is returned correctly
	t.Run("Test 11: Verify ServiceInfo is populated correctly", func(t *testing.T) {
		http := handler.NewHttpClient()
		http.Method("GET")
		http.ExpectedStatus(200)
		query := fmt.Sprintf("page=1&page_size=10&timezone=%s&start_date=%s&end_date=%s", TimeZone, tomorrowFormatted, tomorrowFormatted)
		url := fmt.Sprintf("/employee/%s/appointments?%s", employee.Created.ID.String(), query)
		http.URL(url)
		http.Header("X-Company-ID", cy.Created.ID.String())
		http.Header("X-Auth-Token", cy.Owner.X_Auth_Token)
		http.Send(nil)

		if http.Error != nil {
			t.Fatalf("Failed to get employee appointments: %v", http.Error)
		}

		var appointmentList DTO.AppointmentList
		http.ParseResponse(&appointmentList)

		// Check that ServiceInfo is populated
		if len(appointmentList.ServiceInfo) == 0 {
			t.Error("Expected ServiceInfo to be populated, but it's empty")
		}

		// We should have 2 unique services (service1 and service2)
		if len(appointmentList.ServiceInfo) != 2 {
			t.Errorf("Expected 2 unique services in ServiceInfo, got %d", len(appointmentList.ServiceInfo))
		}

		// Create a map to quickly lookup services
		serviceMap := make(map[string]DTO.ServiceBasicInfo)
		for _, svc := range appointmentList.ServiceInfo {
			serviceMap[svc.ID.String()] = svc

			// Verify service has required fields
			if svc.ID.String() == "" {
				t.Error("Service ID should not be empty")
			}
			if svc.Name == "" {
				t.Error("Service Name should not be empty")
			}
			if svc.Duration == 0 {
				t.Error("Service Duration should not be zero")
			}
		}

		// Verify that all appointments have corresponding service info
		for _, apt := range appointmentList.Appointments {
			if _, exists := serviceMap[apt.ServiceID.String()]; !exists {
				t.Errorf("Appointment has service_id %s but it's not in ServiceInfo", apt.ServiceID.String())
			}
		}

		// Verify specific services are present
		service1Found := false
		service2Found := false
		for _, svc := range appointmentList.ServiceInfo {
			if svc.ID.String() == service1.Created.ID.String() {
				service1Found = true
				if svc.Duration != 60 {
					t.Errorf("Expected service1 duration to be 60, got %d", svc.Duration)
				}
			}
			if svc.ID.String() == service2.Created.ID.String() {
				service2Found = true
				if svc.Duration != 60 {
					t.Errorf("Expected service2 duration to be 60, got %d", svc.Duration)
				}
			}
		}

		if !service1Found {
			t.Error("Expected to find service1 in ServiceInfo")
		}
		if !service2Found {
			t.Error("Expected to find service2 in ServiceInfo")
		}

		// Verify ClientInfo is also populated (regression check)
		if len(appointmentList.ClientInfo) == 0 {
			t.Error("Expected ClientInfo to be populated, but it's empty")
		}

		// We should have 3 unique clients
		if len(appointmentList.ClientInfo) != 3 {
			t.Errorf("Expected 3 unique clients in ClientInfo, got %d", len(appointmentList.ClientInfo))
		}

		// Verify EmployeeInfo is also populated
		if len(appointmentList.EmployeeInfo) == 0 {
			t.Error("Expected EmployeeInfo to be populated, but it's empty")
		}

		// We should have 1 unique employee (only employee1 is being queried)
		if len(appointmentList.EmployeeInfo) != 1 {
			t.Errorf("Expected 1 unique employee in EmployeeInfo, got %d", len(appointmentList.EmployeeInfo))
		}

		// Create a map to quickly lookup employees
		employeeMap := make(map[string]DTO.EmployeeBasicInfo)
		for _, emp := range appointmentList.EmployeeInfo {
			employeeMap[emp.ID.String()] = emp

			// Verify employee has required fields
			if emp.ID.String() == "" {
				t.Error("Employee ID should not be empty")
			}
			if emp.Name == "" {
				t.Error("Employee Name should not be empty")
			}
			if emp.Email == "" {
				t.Error("Employee Email should not be empty")
			}
		}

		// Verify that all appointments have corresponding employee info
		for _, apt := range appointmentList.Appointments {
			if _, exists := employeeMap[apt.EmployeeID.String()]; !exists {
				t.Errorf("Appointment has employee_id %s but it's not in EmployeeInfo", apt.EmployeeID.String())
			}
		}

		// Verify the employee is employee
		if len(appointmentList.EmployeeInfo) > 0 {
			emp := appointmentList.EmployeeInfo[0]
			if emp.ID.String() != employee.Created.ID.String() {
				t.Errorf("Expected EmployeeInfo to contain employee (ID: %s), got ID: %s", employee.Created.ID.String(), emp.ID.String())
			}
		}

		t.Logf("✓ ServiceInfo correctly populated with %d services", len(appointmentList.ServiceInfo))
		t.Logf("✓ ClientInfo correctly populated with %d clients", len(appointmentList.ClientInfo))
		t.Logf("✓ EmployeeInfo correctly populated with %d employees", len(appointmentList.EmployeeInfo))
	})
}
