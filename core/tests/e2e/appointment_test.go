package e2e_test

import (
	"agenda-kaki-go/core"
	models_test "agenda-kaki-go/core/tests/models"
	utils_test "agenda-kaki-go/core/tests/utils"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_Appointment(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	ct := &models_test.Client{}
	ct.Create(t, 200)
	ct.VerifyEmail(t, 200)
	ct.Login(t, 200)
	ct.Update(t, 200, map[string]any{"name": "Updated client Name"})
	ct.GetByEmail(t, 200)
	cy := &models_test.Company{}
	cy.Set(t) // This sets up company, employees (with schedules), branches, services.

	// We will primarily use one employee for these tests.
	// The findValidAppointmentSlot function will determine suitable branch and service.
	if len(cy.Employees) == 0 {
		t.Fatalf("Test setup failed: No employees created by cy.Set(t)")
	}
	baseEmployee := cy.Employees[0]

	a := []*models_test.Appointment{}

	// --- Test Case 0: Successful creation by client ---
	a = append(a, &models_test.Appointment{})
	// Find a valid slot for the base employee. Using time.Local for preferred location.
	slot0, found0 := utils_test.FindValidAppointmentSlot(t, baseEmployee, cy, time.Local)
	if !found0 {
		t.Fatalf("Test setup failed: Could not find any valid appointment slot for employee %s for test case a[0]", baseEmployee.Created.ID)
	}
	// Retrieve the actual Branch and Service objects based on IDs from slot0
	branchForSlot0 := utils_test.GetBranchByID(t, cy, slot0.BranchID)
	serviceForSlot0 := utils_test.GetServiceByID(t, cy, slot0.ServiceID)
	a[0].Create(t, 200, ct.Auth_token, &slot0.StartTimeRFC3339, branchForSlot0, baseEmployee, serviceForSlot0, cy, ct)

	// --- Test Case 1: Another successful creation by client ---
	// The employee's appointments list (baseEmployee.Created.Appointments) should have been updated by a[0].Create(),
	// so findValidAppointmentSlot should now find the *next* available slot.
	a = append(a, &models_test.Appointment{})
	slot1, found1 := utils_test.FindValidAppointmentSlot(t, baseEmployee, cy, time.Local)
	if !found1 {
		t.Fatalf("Test setup failed: Could not find a second valid appointment slot for employee %s for test case a[1]", baseEmployee.Created.ID)
	}
	branchForSlot1 := utils_test.GetBranchByID(t, cy, slot1.BranchID)
	serviceForSlot1 := utils_test.GetServiceByID(t, cy, slot1.ServiceID)
	a[1].Create(t, 200, ct.Auth_token, &slot1.StartTimeRFC3339, branchForSlot1, baseEmployee, serviceForSlot1, cy, ct)

	// --- Test Case 2: Successful creation by company owner ---
	a = append(a, &models_test.Appointment{})
	slot2, found2 := utils_test.FindValidAppointmentSlot(t, baseEmployee, cy, time.Local)
	if !found2 {
		t.Fatalf("Test setup failed: Could not find a third valid appointment slot for employee %s for test case a[2]", baseEmployee.Created.ID)
	}
	branchForSlot2 := utils_test.GetBranchByID(t, cy, slot2.BranchID)
	serviceForSlot2 := utils_test.GetServiceByID(t, cy, slot2.ServiceID)
	a[2].Create(t, 200, cy.Owner.Auth_token, &slot2.StartTimeRFC3339, branchForSlot2, baseEmployee, serviceForSlot2, cy, ct)

	// --- Test Case 3: Attempt to create conflicting appointment (expects 400) ---
	// This test uses the details of the first successfully created appointment (a[0]) to force a conflict.
	if a[0].Created.ID == uuid.Nil {
		t.Fatalf("Prerequisite failed for Test Case 3: a[0].Created appointment is nil. Cannot test conflict.")
	}
	// The start time for the conflict is the same as a[0]'s start time.
	startTimeForConflict := a[0].Created.StartTime.Format(time.RFC3339)
	// The branch, employee, service must be the same as a[0] to ensure a direct conflict.
	// branchForSlot0, baseEmployee, serviceForSlot0 are already the correct objects.
	a = append(a, &models_test.Appointment{})
	a[3].Create(t, 409, ct.Auth_token, &startTimeForConflict, branchForSlot0, baseEmployee, serviceForSlot0, cy, ct)
}