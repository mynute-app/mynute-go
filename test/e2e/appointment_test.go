package e2e_test

import (
	"mynute-go/core"
	"mynute-go/core/src/config/namespace"
	"mynute-go/test/src/handler"
	"mynute-go/test/src/lib"
	"mynute-go/test/src/model"
	coreModel "mynute-go/core/src/config/db/model"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_Appointment(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()
	tt := handler.NewTestErrorHandler(t)

	appEnv := os.Getenv("APP_ENV")
	if appEnv != "test" {
		t.Fatal("APP_ENV is not set to 'test'. Aborting tests to prevent data loss.")
	}

	ct := &model.Client{}
	cy := &model.Company{}

	tt.Describe("Client creation").Test(ct.Set())
	tt.Describe("Company setup").Test(cy.Set())

	baseEmployee := cy.Employees[1]
	Appointments := []*coreModel.Appointment{}

	sp_location, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		t.Fatalf("Failed to load time zone: %v", err)
	}

	branchCache := &map[string]*coreModel.Branch{}
	serviceCache := &map[string]*coreModel.Service{}

	for _, b := range baseEmployee.Created.Branches {
		var branch coreModel.Branch
		if err := handler.NewHttpClient().
			Header(namespace.HeadersKey.Company, baseEmployee.Company.Created.ID.String()).
			Header(namespace.HeadersKey.Auth, baseEmployee.Company.Owner.X_Auth_Token).
			Method("GET").
			URL("/branch/" + b.ID.String()).
			ExpectedStatus(200).
			Send(nil).
			ParseResponse(&branch).Error; err != nil {
			tt.Describe("Fetching branch for employee appointment setup").Test(err)
		}
		(*branchCache)[b.ID.String()] = &branch
	}

	for _, s := range baseEmployee.Created.Services {
		var service coreModel.Service
		if err := handler.NewHttpClient().
			Header(namespace.HeadersKey.Company, baseEmployee.Company.Created.ID.String()).
			Header(namespace.HeadersKey.Auth, baseEmployee.Company.Owner.X_Auth_Token).
			Method("GET").
			URL("/service/" + s.ID.String()).
			ExpectedStatus(200).
			Send(nil).
			ParseResponse(&service).Error; err != nil {
			tt.Describe("Fetching service for employee appointment setup").Test(err)
		}
		(*serviceCache)[s.ID.String()] = &service
	}

	// --- Test Case 0 ---
	var a0 model.Appointment
	slot0, found0, err := a0.FindValidAppointmentSlot(baseEmployee, sp_location, branchCache, serviceCache)
	tt.Describe("Finding valid appointment slot for base employee - slot0").Test(err)
	if !found0 {
		t.Logf("Employee Work Schedule: %+v", baseEmployee.Created.WorkSchedule)
		t.Fatalf("Setup failed: No valid appointment slot for test case a[0]")
	}

	branch0, err := lib.GetBranchByID(cy, slot0.BranchID)
	tt.Describe("Getting branch for slot0").Test(err)
	service0, err := lib.GetServiceByID(cy, slot0.ServiceID)
	tt.Describe("Getting service for slot0").Test(err)

	tt.Describe("Creating appointment a[0]").Test(
		a0.Create(200, ct.X_Auth_Token, nil, &slot0.StartTimeRFC3339, slot0.TimeZone, branch0, baseEmployee, service0, cy, ct),
	)

	Appointments = append(Appointments, a0.Created)

	// --- Test Case 1 ---
	var a1 model.Appointment
	slot1, found1, err := a1.FindValidAppointmentSlot(baseEmployee, sp_location, branchCache, serviceCache)
	tt.Describe("Finding valid appointment slot for base employee - slot1").Test(err)
	if !found1 {
		t.Logf("Employee Work Schedule: %+v", baseEmployee.Created.WorkSchedule)
		t.Fatalf("Setup failed: No valid appointment slot for test case a[1]")
	}

	branch1, err := lib.GetBranchByID(cy, slot1.BranchID)
	tt.Describe("Getting branch for slot1").Test(err)
	service1, err := lib.GetServiceByID(cy, slot1.ServiceID)
	tt.Describe("Getting service for slot1").Test(err)
	tt.Describe("Creating appointment a[1]").Test(
		a1.Create(200, ct.X_Auth_Token, nil, &slot1.StartTimeRFC3339, slot1.TimeZone, branch1, baseEmployee, service1, cy, ct),
	)
	Appointments = append(Appointments, a1.Created)

	// --- Test Case 2 ---
	var a2 model.Appointment
	slot2, found2, err := a2.FindValidAppointmentSlot(baseEmployee, sp_location, branchCache, serviceCache)
	tt.Describe("Finding valid appointment slot for base employee - slot2").Test(err)
	if !found2 {
		t.Logf("Employee Work Schedule: %+v", baseEmployee.Created.WorkSchedule)
		t.Fatalf("Setup failed: No valid appointment slot for test case a[2]")
	}

	branch2, err := lib.GetBranchByID(cy, slot2.BranchID)
	tt.Describe("Getting branch for slot2").Test(err)
	service2, err := lib.GetServiceByID(cy, slot2.ServiceID)
	tt.Describe("Getting service for slot2").Test(err)
	tt.Describe("Creating appointment a[2]").Test(
		a2.Create(200, cy.Owner.X_Auth_Token, nil, &slot2.StartTimeRFC3339, slot2.TimeZone, branch2, baseEmployee, service2, cy, ct),
	)
	Appointments = append(Appointments, a2.Created)

	// --- Test Case 3 ---
	if Appointments[0].ID == uuid.Nil {
		t.Fatalf("Setup failed: a[0] is nil, cannot test conflict")
	}

	ct1 := &model.Client{}
	tt.Describe("Client creation for conflicting appointments test").Test(ct1.Set())

	var a3 model.Appointment
	tt.Describe("Creating conflicting appointment a[3] with client token").Test(
		a3.Create(400, ct1.X_Auth_Token, nil, &slot0.StartTimeRFC3339, slot0.TimeZone, branch0, baseEmployee, service0, cy, ct1),
	)

	tt.Describe("Creating conflicting appointment a[3] with employee token").Test(
		a3.Create(400, baseEmployee.X_Auth_Token, nil, &slot0.StartTimeRFC3339, slot0.TimeZone, branch0, baseEmployee, service0, cy, ct1),
	)
}
