package e2e_test

import (
	"mynute-go/src"
	"mynute-go/src/config/db/model"
	"mynute-go/src/config/namespace"
	handlerT "mynute-go/test/src/handlers"
	modelT "mynute-go/test/src/models"
	utilsT "mynute-go/test/src/utils"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_Appointment(t *testing.T) {
	server := src.NewServer().Run("parallel")
	defer server.Shutdown()
	tt := handlerT.NewTestErrorHandler(t)

	appEnv := os.Getenv("APP_ENV")
	if appEnv != "test" {
		t.Fatal("APP_ENV is not set to 'test'. Aborting tests to prevent data loss.")
	}

	ct := &modelT.Client{}
	cy := &modelT.Company{}

	tt.Describe("Client creation").Test(ct.Set())
	tt.Describe("Company setup").Test(cy.Set())

	baseEmployee := cy.Employees[1]
	Appointments := []*modelT.Appointment{}

	sp_location, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		t.Fatalf("Failed to load time zone: %v", err)
	}

	branchCache := &map[string]*model.Branch{}
	serviceCache := &map[string]*model.Service{}

	for _, b := range baseEmployee.Created.Branches {
		var branch model.Branch
		if err := handlerT.NewHttpClient().
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
		var service model.Service
		if err := handlerT.NewHttpClient().
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
	var a0 modelT.Appointment
	slot0, found0, err := a0.FindValidAppointmentSlot(baseEmployee, sp_location, branchCache, serviceCache)
	tt.Describe("Finding valid appointment slot for base employee - slot0").Test(err)
	if !found0 {
		t.Logf("Employee Work Schedule: %+v", baseEmployee.Created.WorkSchedule)
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
	slot1, found1, err := a1.FindValidAppointmentSlot(baseEmployee, sp_location, branchCache, serviceCache)
	tt.Describe("Finding valid appointment slot for base employee - slot1").Test(err)
	if !found1 {
		t.Logf("Employee Work Schedule: %+v", baseEmployee.Created.WorkSchedule)
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
	slot2, found2, err := a2.FindValidAppointmentSlot(baseEmployee, sp_location, branchCache, serviceCache)
	tt.Describe("Finding valid appointment slot for base employee - slot2").Test(err)
	if !found2 {
		t.Logf("Employee Work Schedule: %+v", baseEmployee.Created.WorkSchedule)
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

	ct1 := &modelT.Client{}
	tt.Describe("Client creation for conflicting appointments test").Test(ct1.Set())

	var a3 modelT.Appointment
	tt.Describe("Creating conflicting appointment a[3] with client token").Test(
		a3.Create(400, ct1.X_Auth_Token, nil, &slot0.StartTimeRFC3339, slot0.TimeZone, branch0, baseEmployee, service0, cy, ct1),
	)

	tt.Describe("Creating conflicting appointment a[3] with employee token").Test(
		a3.Create(400, baseEmployee.X_Auth_Token, nil, &slot0.StartTimeRFC3339, slot0.TimeZone, branch0, baseEmployee, service0, cy, ct1),
	)
}
