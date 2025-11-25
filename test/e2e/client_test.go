package e2e_test

import (
	"mynute-go/core"
	"mynute-go/core/src/lib"
	FileBytes "mynute-go/core/src/lib/file_bytes"
	"mynute-go/test/src/handler"
	"mynute-go/test/src/model"
	"testing"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	Created    model.Client
	Auth_token string
}

func Test_Client(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)

	client := &model.Client{}

	tt.Describe("Client creation").Test(client.Create(200))

	tt.Describe("Client get by email").Test(client.GetByEmail(200))

	tt.Describe("Login with password").Test(client.Login(401, "password"))

	tt.Describe("Verify email").Test(client.VerifyEmail(200))

	tt.Describe("Login with email code").Test(client.Login(200, "email_code"))

	tt.Describe("Login by password with invalid password").Test(client.LoginByPassword(401, "invalid_password"))

	tt.Describe("Login with password").Test(client.Login(200, "password"))

	tt.Describe("Client update").Test(client.Update(200, map[string]any{
		"name": "Updated Client Name",
	}))

	tt.Describe("Client update").Test(client.Update(400, map[string]any{
		"name":     "Should Fail Update on Client Name",
		"password": "newpswrd123",
	}))

	new_password := lib.GenerateValidPassword()

	tt.Describe("Client update").Test(client.Update(200, map[string]any{
		"name":     "Should Succeed Update on Client Name",
		"password": new_password,
	}))

	tt.Describe("Client update").Test(client.Update(401, map[string]any{
		"password": "NewPswrd1@!",
	}))

	client.Created.Password = new_password // Update the password in the client model

	// Re-login with new password to get a fresh token
	tt.Describe("Login with new password").Test(client.LoginByPassword(200, new_password))

	// Test password reset by email
	tt.Describe("Reset password by email").Test(client.ResetPasswordByEmail(200))

	// Test that the old password no longer works
	tt.Describe("Login with old password fails").Test(client.LoginByPassword(401, new_password))

	// Test that new password from email works
	tt.Describe("Login with password from email").Test(client.LoginByPassword(200, client.Created.Password))

	tt.Describe("Upload profile image").Test(client.UploadImages(200, map[string][]byte{
		"profile": FileBytes.PNG_FILE_1,
	}, nil))

	tt.Describe("Get profile image").Test(client.GetImage(200, client.Created.Meta.Design.Images.Profile.URL, &FileBytes.PNG_FILE_1))

	tt.Describe("Overwrite profile image").Test(client.UploadImages(200, map[string][]byte{
		"profile": FileBytes.PNG_FILE_3,
	}, nil))

	tt.Describe("Get overwritten profile image").Test(client.GetImage(200, client.Created.Meta.Design.Images.Profile.URL, &FileBytes.PNG_FILE_3))

	img_url := client.Created.Meta.Design.Images.Profile.URL

	tt.Describe("Delete profile image").Test(client.DeleteImages(200, []string{"profile"}, nil))

	tt.Describe("Get deleted profile image").Test(client.GetImage(404, img_url, nil))

	tt.Describe("Upload profile image again logged in with email code").Test(client.UploadImages(200, map[string][]byte{
		"profile": FileBytes.PNG_FILE_1,
	}, nil))

	tt.Describe("Client deletion").Test(client.Delete(200))

	tt.Describe("Get deleted client by email").Test(client.GetByEmail(404))
}

func Test_Client_DoubleBooking_Prevention(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)

	// Create client
	client := &model.Client{}
	tt.Describe("Client creation").Test(client.Create(200))
	tt.Describe("Client get by email").Test(client.GetByEmail(200))
	tt.Describe("Verify email").Test(client.VerifyEmail(200))
	tt.Describe("Login with email code").Test(client.Login(200, "email_code"))

	// Create company with employees, branches, and services
	company := &model.Company{}
	empN := 2
	branchN := 1
	serviceN := 1
	tt.Describe("Company creation").Test(company.CreateCompanyRandomly(empN, branchN, serviceN))

	if len(company.Services) == 0 {
		t.Fatal("No services created")
	}
	if len(company.Branches) == 0 {
		t.Fatal("No branches created")
	}
	if len(company.Employees) == 0 {
		t.Fatal("No employees created")
	}

	service := company.Services[0]

	TimeZone := "America/Sao_Paulo"
	client_public_id := client.Created.ID.String()

	// Find a valid appointment slot
	slot, err := service.FindValidRandomAppointmentSlot(TimeZone, &client_public_id)
	if err != nil {
		t.Fatalf("Failed to find valid appointment slot: %v", err)
	}

	// Find the branch and employee for this slot
	var slotBranch *model.Branch
	for _, b := range company.Branches {
		if b.Created.ID.String() == slot.BranchID {
			slotBranch = b
			break
		}
	}
	if slotBranch == nil {
		t.Fatal("Failed to find branch for slot")
	}

	var slotEmployee *model.Employee
	for _, e := range company.Employees {
		if e.Created.ID.String() == slot.EmployeeID {
			slotEmployee = e
			break
		}
	}
	if slotEmployee == nil {
		t.Fatal("Failed to find employee for slot")
	}

	// Test 1: Create first appointment successfully
	appointment1 := &model.Appointment{}
	tt.Describe("Create first appointment successfully").
		Test(appointment1.Create(200, client.X_Auth_Token, nil, &slot.StartTimeRFC3339, slot.TimeZone, slotBranch, slotEmployee, service, company, client))

	// Test 2: Try to create overlapping appointment - should fail with 400
	appointment2 := &model.Appointment{}
	tt.Describe("Prevent creating overlapping appointment for same client").
		Test(appointment2.Create(400, client.X_Auth_Token, nil, &slot.StartTimeRFC3339, slot.TimeZone, slotBranch, slotEmployee, service, company, client))

	// Test 3: Cancel first appointment
	tt.Describe("Cancel first appointment").
		Test(appointment1.Cancel(200, client.X_Auth_Token, nil))

	// Test 4: After cancellation, should be able to create appointment in the same time slot
	appointment3 := &model.Appointment{}
	tt.Describe("Create appointment in same slot after cancellation").
		Test(appointment3.Create(200, client.X_Auth_Token, nil, &slot.StartTimeRFC3339, slot.TimeZone, slotBranch, slotEmployee, service, company, client))

	// Test 5: Create another appointment in a different time slot
	slot2, err := service.FindValidRandomAppointmentSlot(TimeZone, &client_public_id)
	if err != nil {
		t.Fatalf("Failed to find second valid appointment slot: %v", err)
	}

	// Find the branch and employee for the second slot
	var slotBranch2 *model.Branch
	for _, b := range company.Branches {
		if b.Created.ID.String() == slot2.BranchID {
			slotBranch2 = b
			break
		}
	}
	if slotBranch2 == nil {
		t.Fatal("Failed to find branch for second slot")
	}

	var slotEmployee2 *model.Employee
	for _, e := range company.Employees {
		if e.Created.ID.String() == slot2.EmployeeID {
			slotEmployee2 = e
			break
		}
	}
	if slotEmployee2 == nil {
		t.Fatal("Failed to find employee for second slot")
	}

	appointment4 := &model.Appointment{}
	tt.Describe("Create second appointment in different time slot").
		Test(appointment4.Create(200, client.X_Auth_Token, nil, &slot2.StartTimeRFC3339, slot2.TimeZone, slotBranch2, slotEmployee2, service, company, client))

	// Test 6: Try to create a third appointment that overlaps with appointment3 - should fail
	appointment5 := &model.Appointment{}
	tt.Describe("Prevent creating third overlapping appointment").
		Test(appointment5.Create(400, client.X_Auth_Token, nil, &slot.StartTimeRFC3339, slot.TimeZone, slotBranch, slotEmployee, service, company, client))

	// Clean up
	tt.Describe("Delete client").Test(client.Delete(200))
}

func Test_Client_Appointments_Pagination_And_Filters(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)

	// Create client
	client := &model.Client{}
	tt.Describe("Client creation").Test(client.Create(200))
	tt.Describe("Client get by email").Test(client.GetByEmail(200))
	tt.Describe("Verify email").Test(client.VerifyEmail(200))
	tt.Describe("Login with email code").Test(client.Login(200, "email_code"))

	// Create company with employees, branches, and services
	company := &model.Company{}
	empN := 2
	branchN := 1
	serviceN := 1
	tt.Describe("Company creation").Test(company.CreateCompanyRandomly(empN, branchN, serviceN))

	companyID := company.Created.ID.String()

	if len(company.Services) == 0 {
		t.Fatal("No services created")
	}
	if len(company.Branches) == 0 {
		t.Fatal("No branches created")
	}
	if len(company.Employees) == 0 {
		t.Fatal("No employees created")
	}

	service := company.Services[0]
	TimeZone := "America/Sao_Paulo"
	client_public_id := client.Created.ID.String()

	// Create 6 appointments across different dates and times
	appointments := []*model.Appointment{}
	targetAppointments := 6

	for i := 0; i < targetAppointments; i++ {
		slot, err := service.FindValidRandomAppointmentSlot(TimeZone, &client_public_id)
		if err != nil {
			t.Logf("Warning: Could not find slot for appointment %d: %v", i+1, err)
			continue
		}

		var slotBranch *model.Branch
		for _, b := range company.Branches {
			if b.Created.ID.String() == slot.BranchID {
				slotBranch = b
				break
			}
		}
		if slotBranch == nil {
			t.Logf("Warning: Could not find branch for slot %d", i+1)
			continue
		}

		var slotEmployee *model.Employee
		for _, e := range company.Employees {
			if e.Created.ID.String() == slot.EmployeeID {
				slotEmployee = e
				break
			}
		}
		if slotEmployee == nil {
			t.Logf("Warning: Could not find employee for slot %d", i+1)
			continue
		}

		appointment := &model.Appointment{}
		err = appointment.Create(200, client.X_Auth_Token, nil, &slot.StartTimeRFC3339, slot.TimeZone, slotBranch, slotEmployee, service, company, client)
		if err != nil {
			t.Logf("Warning: Could not create appointment %d: %v", i+1, err)
			continue
		}
		appointments = append(appointments, appointment)
	}

	if len(appointments) < 3 {
		t.Fatalf("Need at least 3 appointments for testing, only created %d", len(appointments))
	}

	t.Logf("Created %d appointments for testing", len(appointments))

	// Test 1: Get appointments with pagination - page 1 with page_size 2
	result1, err := client.GetAppointments(200, 1, 2, "", "", "", TimeZone, nil, &companyID)
	tt.Describe("Get appointments page 1 with page_size 2").Test(err)
	if err == nil {
		if len(result1.Appointments) > 2 {
			t.Errorf("Expected at most 2 appointments on page 1, got %d", len(result1.Appointments))
		}
		if result1.Page != 1 {
			t.Errorf("Expected page 1, got %d", result1.Page)
		}
		if result1.PageSize != 2 {
			t.Errorf("Expected page_size 2, got %d", result1.PageSize)
		}
		if result1.TotalCount != len(appointments) {
			t.Errorf("Expected total_count %d, got %d", len(appointments), result1.TotalCount)
		}
		t.Logf("✓ Pagination page 1: got %d appointments, total %d", len(result1.Appointments), result1.TotalCount)
	}

	// Test 2: Get appointments with pagination - page 2 with page_size 2
	result2, err := client.GetAppointments(200, 2, 2, "", "", "", TimeZone, nil, &companyID)
	tt.Describe("Get appointments page 2 with page_size 2").Test(err)
	if err == nil {
		if len(result2.Appointments) > 2 {
			t.Errorf("Expected at most 2 appointments on page 2, got %d", len(result2.Appointments))
		}
		if result2.Page != 2 {
			t.Errorf("Expected page 2, got %d", result2.Page)
		}
		t.Logf("✓ Pagination page 2: got %d appointments", len(result2.Appointments))
	}

	// Test 3: Get all appointments (default pagination)
	resultAll, err := client.GetAppointments(200, 1, 10, "", "", "", TimeZone, nil, &companyID)
	tt.Describe("Get all appointments with default pagination").Test(err)
	if err == nil {
		if len(resultAll.Appointments) != len(appointments) {
			t.Errorf("Expected %d appointments, got %d", len(appointments), len(resultAll.Appointments))
		}
		if len(resultAll.ClientInfo) == 0 {
			t.Error("Expected ClientInfo to be populated")
		} else if resultAll.ClientInfo[0].ID.String() != client.Created.ID.String() {
			t.Error("Expected ClientInfo[0].ID to match client ID")
		}
		t.Logf("✓ Got all %d appointments with ClientInfo populated", len(resultAll.Appointments))
	}

	// Test 4: Cancel one appointment to test cancelled filter
	if len(appointments) > 0 {
		tt.Describe("Cancel first appointment").Test(appointments[0].Cancel(200, client.X_Auth_Token, nil))

		// Test 5: Filter by cancelled=false (should exclude the cancelled one)
		resultNonCancelled, err := client.GetAppointments(200, 1, 10, "", "", "false", TimeZone, nil, &companyID)
		tt.Describe("Get only non-cancelled appointments").Test(err)
		if err == nil {
			expected := len(appointments) - 1
			if len(resultNonCancelled.Appointments) != expected {
				t.Errorf("Expected %d non-cancelled appointments, got %d", expected, len(resultNonCancelled.Appointments))
			}
			for _, apt := range resultNonCancelled.Appointments {
				if apt.IsCancelled {
					t.Error("Found cancelled appointment when filtering cancelled=false")
					break
				}
			}
			t.Logf("✓ Cancelled filter=false: got %d non-cancelled appointments", len(resultNonCancelled.Appointments))
		}

		// Test 6: Filter by cancelled=true (should only get the cancelled one)
		resultCancelled, err := client.GetAppointments(200, 1, 10, "", "", "true", TimeZone, nil, &companyID)
		tt.Describe("Get only cancelled appointments").Test(err)
		if err == nil {
			if len(resultCancelled.Appointments) != 1 {
				t.Errorf("Expected 1 cancelled appointment, got %d", len(resultCancelled.Appointments))
			} else if !resultCancelled.Appointments[0].IsCancelled {
				t.Error("Expected cancelled appointment, but got non-cancelled")
			}
			t.Logf("✓ Cancelled filter=true: got %d cancelled appointment", len(resultCancelled.Appointments))
		}
	}

	// Test 7: Date range filtering
	if len(appointments) >= 3 {
		// Get the earliest appointment's date (skip if nil or Created is nil)
		var earliestApt *model.Appointment
		for _, apt := range appointments {
			if apt != nil && apt.Created != nil && apt.Created.ID != uuid.Nil {
				if earliestApt == nil || apt.Created.StartTime.Before(earliestApt.Created.StartTime) {
					earliestApt = apt
				}
			}
		}

		if earliestApt == nil {
			t.Skip("No valid appointments for date range test")
		}

		// Format dates in DD/MM/YYYY format
		startDate := earliestApt.Created.StartTime.Format("02/01/2006")
		endDate := earliestApt.Created.StartTime.AddDate(0, 0, 1).Format("02/01/2006") // Next day

		resultDateRange, err := client.GetAppointments(200, 1, 10, startDate, endDate, "", TimeZone, nil, &companyID)
		tt.Describe("Get appointments within date range").Test(err)
		if err == nil {
			if len(resultDateRange.Appointments) < 1 {
				t.Errorf("Expected at least 1 appointment in date range, got %d", len(resultDateRange.Appointments))
			}
			// Verify all appointments are within the date range
			for _, apt := range resultDateRange.Appointments {
				aptStartTime, err := time.Parse(time.RFC3339, apt.StartTime)
				if err != nil {
					t.Errorf("Failed to parse appointment start time %s: %v", apt.StartTime, err)
					continue
				}
				if aptStartTime.Before(earliestApt.Created.StartTime) {
					t.Error("Found appointment before start_date")
					break
				}
			}
			t.Logf("✓ Date range filter (%s to %s): got %d appointments", startDate, endDate, len(resultDateRange.Appointments))
		}
	}

	// Test 8: Invalid date format should return 400
	_, err = client.GetAppointments(400, 1, 10, "2024-01-15", "", "", TimeZone, nil, &companyID)
	tt.Describe("Reject invalid start_date format").Test(err)
	if err == nil {
		t.Logf("✓ Invalid date format correctly rejected")
	}

	// Test 9: Missing timezone should return 400
	_, err = client.GetAppointments(400, 1, 10, "", "", "", "", nil, &companyID)
	tt.Describe("Reject missing timezone parameter").Test(err)
	if err == nil {
		t.Logf("✓ Missing timezone correctly rejected")
	}

	// Test 10: Invalid cancelled value should return 400
	_, err = client.GetAppointments(400, 1, 10, "", "", "invalid", TimeZone, nil, &companyID)
	tt.Describe("Reject invalid cancelled parameter").Test(err)
	if err == nil {
		t.Logf("✓ Invalid cancelled value correctly rejected")
	}

	// Test 11: Date range exceeding 90 days should return 400
	if len(appointments) > 0 {
		var firstApt *model.Appointment
		for _, apt := range appointments {
			if apt != nil && apt.Created != nil && apt.Created.ID != uuid.Nil {
				firstApt = apt
				break
			}
		}

		if firstApt != nil {
			startDate := firstApt.Created.StartTime.Format("02/01/2006")
			endDate := firstApt.Created.StartTime.AddDate(0, 0, 91).Format("02/01/2006") // 91 days later

			_, err = client.GetAppointments(400, 1, 10, startDate, endDate, "", TimeZone, nil, &companyID)
			tt.Describe("Reject date range exceeding 90 days").Test(err)
			if err == nil {
				t.Logf("✓ Date range > 90 days correctly rejected")
			}
		}
	}

	// Test 12: End date before start date should return 400
	if len(appointments) > 0 {
		var firstApt *model.Appointment
		for _, apt := range appointments {
			if apt != nil && apt.Created != nil && apt.Created.ID != uuid.Nil {
				firstApt = apt
				break
			}
		}

		if firstApt != nil {
			startDate := firstApt.Created.StartTime.Format("02/01/2006")
			endDate := firstApt.Created.StartTime.AddDate(0, 0, -1).Format("02/01/2006") // 1 day before

			_, err = client.GetAppointments(400, 1, 10, startDate, endDate, "", TimeZone, nil, &companyID)
			tt.Describe("Reject end_date before start_date").Test(err)
			if err == nil {
				t.Logf("✓ end_date before start_date correctly rejected")
			}
		}
	}

	// Test 13: Combined filters - non-cancelled appointments with date range and pagination
	if len(appointments) >= 3 {
		var firstApt *model.Appointment
		for _, apt := range appointments {
			if apt != nil && apt.Created != nil && apt.Created.ID != uuid.Nil {
				firstApt = apt
				break
			}
		}

		if firstApt != nil {
			startDate := firstApt.Created.StartTime.AddDate(0, 0, -7).Format("02/01/2006") // 7 days before
			endDate := firstApt.Created.StartTime.AddDate(0, 0, 30).Format("02/01/2006")   // 30 days after

			resultCombined, err := client.GetAppointments(200, 1, 3, startDate, endDate, "false", TimeZone, nil, &companyID)
			tt.Describe("Get appointments with combined filters").Test(err)
			if err == nil {
				// All returned appointments should be non-cancelled
				for _, apt := range resultCombined.Appointments {
					if apt.IsCancelled {
						t.Error("Found cancelled appointment when filtering cancelled=false")
						break
					}
				}
				t.Logf("✓ Combined filters: got %d appointments (page_size=3, cancelled=false, date range)", len(resultCombined.Appointments))
			}
		}
	}

	// Clean up
	tt.Describe("Delete client").Test(client.Delete(200))
}
