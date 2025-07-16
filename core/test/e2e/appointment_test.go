package e2e_test

import (
	"mynute-go/core"
	handlerT "mynute-go/core/test/handlers"
	modelT "mynute-go/core/test/models"
	utilsT "mynute-go/core/test/utils"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_Appointment(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()
	tt := handlerT.NewTestErrorHandler(t)

	ct := &modelT.Client{}
	cy := &modelT.Company{}

	tt.Describe("Client creation").Test(ct.Create(200))
	tt.Describe("Client email verification").Test(ct.VerifyEmail(200))
	tt.Describe("Client login").Test(ct.Login(200))
	tt.Describe("Client get by email").Test(ct.GetByEmail(200))
	tt.Describe("Company setup").Test(cy.Set())

	baseEmployee := cy.Employees[1]
	Appointments := []*modelT.Appointment{}

	// --- Test Case 0 ---
	var a0 modelT.Appointment
	slot0, found0, err := a0.FindValidAppointmentSlot(baseEmployee, time.Now().Location())
	tt.Describe("Finding valid appointment slot for base employee - slot0").Test(err)
	if !found0 {
		t.Logf("Employee Work Schedule: %+v", baseEmployee.Created.EmployeeWorkSchedule)
		t.Fatalf("Setup failed: No valid appointment slot for test case a[0]")
	}

	branch0, err := utilsT.GetBranchByID(cy, slot0.BranchID)
	tt.Describe("Getting branch for slot0").Test(err)
	service0, err := utilsT.GetServiceByID(cy, slot0.ServiceID)
	tt.Describe("Getting service for slot0").Test(err)

	tt.Describe("Creating appointment a[0]").Test(
		a0.Create(200, ct.X_Auth_Token, nil, &slot0.StartTimeRFC3339, slot0.TimeZone, branch0, baseEmployee, service0, cy, ct),
	)

	Appointments = append(Appointments, &a0)

	// --- Test Case 1 ---
	var a1 modelT.Appointment
	slot1, found1, err := a1.FindValidAppointmentSlot(baseEmployee, time.Local)
	tt.Describe("Finding valid appointment slot for base employee - slot1").Test(err)
	if !found1 {
		t.Logf("Employee Work Schedule: %+v", baseEmployee.Created.EmployeeWorkSchedule)
		t.Fatalf("Setup failed: No valid appointment slot for test case a[1]")
	}

	branch1, err := utilsT.GetBranchByID(cy, slot1.BranchID)
	tt.Describe("Getting branch for slot1").Test(err)
	service1, err := utilsT.GetServiceByID(cy, slot1.ServiceID)
	tt.Describe("Getting service for slot1").Test(err)
	tt.Describe("Creating appointment a[1]").Test(
		a1.Create(200, ct.X_Auth_Token, nil, &slot1.StartTimeRFC3339, slot1.TimeZone, branch1, baseEmployee, service1, cy, ct),
	)
	Appointments = append(Appointments, &a1)

	// --- Test Case 2 ---
	var a2 modelT.Appointment
	slot2, found2, err := a2.FindValidAppointmentSlot(baseEmployee, time.Local)
	tt.Describe("Finding valid appointment slot for base employee - slot2").Test(err)
	if !found2 {
		t.Logf("Employee Work Schedule: %+v", baseEmployee.Created.EmployeeWorkSchedule)
		t.Fatalf("Setup failed: No valid appointment slot for test case a[2]")
	}

	branch2, err := utilsT.GetBranchByID(cy, slot2.BranchID)
	tt.Describe("Getting branch for slot2").Test(err)
	service2, err := utilsT.GetServiceByID(cy, slot2.ServiceID)
	tt.Describe("Getting service for slot2").Test(err)
	tt.Describe("Creating appointment a[2]").Test(
		a2.Create(200, cy.Owner.X_Auth_Token, nil, &slot2.StartTimeRFC3339, slot2.TimeZone, branch2, baseEmployee, service2, cy, ct),
	)
	Appointments = append(Appointments, &a2)

	// --- Test Case 3 ---
	if Appointments[0].Created.ID == uuid.Nil {
		t.Fatalf("Setup failed: a[0] is nil, cannot test conflict")
	}
	startTimeConflict := Appointments[0].Created.StartTime.Format(time.RFC3339)
	timeZoneConflict := Appointments[0].Created.TimeZone

	var a3 modelT.Appointment
	tt.Describe("Creating conflicting appointment a[3]").Test(
		a3.Create(409, ct.X_Auth_Token, nil, &startTimeConflict, timeZoneConflict, branch0, baseEmployee, service0, cy, ct),
	)
}
