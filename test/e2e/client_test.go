package e2e_test

import (
	"mynute-go/core"
	"mynute-go/core/src/lib"
	FileBytes "mynute-go/core/src/lib/file_bytes"
	"mynute-go/test/src/handler"
	"mynute-go/test/src/model"

	"testing"
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
