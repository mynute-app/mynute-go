package e2e_test

import (
	"agenda-kaki-go/core"
	handlerT "agenda-kaki-go/core/test/handlers"
	modelT "agenda-kaki-go/core/test/models"
	utilsT "agenda-kaki-go/core/test/utils"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_Appointment(t *testing.T) {
	var err error
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()
	tt := handlerT.NewTestErrorHandler(t)
	ct := &modelT.Client{}
	tt.Test(ct.Create(200), "Client creation")
	tt.Test(ct.VerifyEmail(200), "Client email verification")
	tt.Test(ct.Login(200), "Client login")
	tt.Test(ct.Update(200, map[string]any{"name": "Updated client Name"}), "Client update")
	tt.Test(ct.GetByEmail(200), "Client get by email")
	cy := &modelT.Company{}
	tt.Test(cy.Set(), "Company setup") // This sets up company, employees (with schedules), branches, services.

	baseEmployee := cy.Owner

	a := []*modelT.Appointment{}

	// --- Test Case 0: Successful creation by client ---
	a = append(a, &modelT.Appointment{})
	// Find a valid slot for the base employee. Using time.Local for preferred location.
	slot0, found0, err := utilsT.FindValidAppointmentSlot(baseEmployee, cy, time.Local)
	tt.Test(err, "Finding valid appointment slot for base employee")
	if !found0 {
		t.Logf("Employee Work Schedule: %+v", baseEmployee.Created.WorkSchedule)
		t.Fatalf("Test setup failed: Could not find any valid appointment slot for employee %s for test case a[0]", baseEmployee.Created.ID)
	}
	// Retrieve the actual Branch and Service objects based on IDs from slot0
	branchForSlot0, err := utilsT.GetBranchByID(cy, slot0.BranchID)
	tt.Test(err, "Getting branch for slot 0")
	serviceForSlot0, err := utilsT.GetServiceByID(cy, slot0.ServiceID)
	tt.Test(err, "Getting service for slot 0")
	a[0].Create(200, ct.X_Auth_Token, nil, &slot0.StartTimeRFC3339, branchForSlot0, baseEmployee, serviceForSlot0, cy, ct)

	// --- Test Case 1: Another successful creation by client ---
	// The employee's appointments list (baseEmployee.Created.Appointments) should have been updated by a[0].Create(),
	// so findValidAppointmentSlot should now find the *next* available slot.
	a = append(a, &modelT.Appointment{})
	slot1, found1, err := utilsT.FindValidAppointmentSlot(baseEmployee, cy, time.Local)
	tt.Test(err, "Finding valid appointment slot for base employee")
	if !found1 {
		t.Logf("Employee Work Schedule: %+v", baseEmployee.Created.WorkSchedule)
		t.Fatalf("Test setup failed: Could not find a second valid appointment slot for employee %s for test case a[1]", baseEmployee.Created.ID)
	}
	branchForSlot1, err := utilsT.GetBranchByID(cy, slot1.BranchID)
	tt.Test(err, "Getting branch for slot 1")
	serviceForSlot1, err := utilsT.GetServiceByID(cy, slot1.ServiceID)
	tt.Test(err, "Getting service for slot 1")
	a[1].Create(200, ct.X_Auth_Token, nil, &slot1.StartTimeRFC3339, branchForSlot1, baseEmployee, serviceForSlot1, cy, ct)

	// --- Test Case 2: Successful creation by company owner ---
	a = append(a, &modelT.Appointment{})
	slot2, found2, err := utilsT.FindValidAppointmentSlot(baseEmployee, cy, time.Local)
	tt.Test(err, "Finding valid appointment slot for base employee")
	if !found2 {
		t.Logf("Employee Work Schedule: %+v", baseEmployee.Created.WorkSchedule)
		t.Fatalf("Test setup failed: Could not find a third valid appointment slot for employee %s for test case a[2]", baseEmployee.Created.ID)
	}
	branchForSlot2, err := utilsT.GetBranchByID(cy, slot2.BranchID)
	tt.Test(err, "Getting branch for slot 2")
	serviceForSlot2, err := utilsT.GetServiceByID(cy, slot2.ServiceID)
	tt.Test(err, "Getting service for slot 2")
	a[2].Create(200, cy.Owner.X_Auth_Token, nil, &slot2.StartTimeRFC3339, branchForSlot2, baseEmployee, serviceForSlot2, cy, ct)

	// --- Test Case 3: Attempt to create conflicting appointment (expects 400) ---
	// This test uses the details of the first successfully created appointment (a[0]) to force a conflict.
	if a[0].Created.ID == uuid.Nil {
		t.Fatalf("Prerequisite failed for Test Case 3: a[0].Created appointment is nil. Cannot test conflict.")
	}
	// The start time for the conflict is the same as a[0]'s start time.
	startTimeForConflict := a[0].Created.StartTime.Format(time.RFC3339)
	// The branch, employee, service must be the same as a[0] to ensure a direct conflict.
	// branchForSlot0, baseEmployee, serviceForSlot0 are already the correct objects.
	a = append(a, &modelT.Appointment{})
	a[3].Create(409, ct.X_Auth_Token, nil, &startTimeForConflict, branchForSlot0, baseEmployee, serviceForSlot0, cy, ct)
}
