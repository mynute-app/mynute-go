package e2e_test

import (
	"fmt"
	"mynute-go/core"
	"mynute-go/test/src/handler"
	testModel "mynute-go/test/src/model"
	"os"
	"testing"

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

	ct := &testModel.Client{}
	cy := &testModel.Company{}

	tt.Describe("Client creation").Test(ct.Set())
	tt.Describe("Company setup").Test(cy.Set())

	Appointments := []*testModel.Appointment{}

	service, err := cy.GetRandomService()
	tt.Describe("Getting random service from company services").Test(err)

	GetBranchByID := func(branchID string, company *testModel.Company) (*testModel.Branch, error) {
		for _, branch := range company.Branches {
			if branch.Created.ID.String() == branchID {
				return branch, nil
			}
		}
		return nil, fmt.Errorf("branch not found")
	}

	GetEmployeeByID := func(employeeID string, company *testModel.Company) (*testModel.Employee, error) {
		for _, employee := range company.Employees {
			if employee.Created.ID.String() == employeeID {
				return employee, nil
			}
		}
		return nil, fmt.Errorf("employee not found")
	}

	// Test Case 0: Successful appointment creation
	slot0, err := service.FindValidRandomAppointmentSlot()
	tt.Describe("Finding valid appointment slot for service - a[0]").Test(err)

	slot0Branch, err := GetBranchByID(slot0.BranchID, cy)
	tt.Describe("Getting branch for slot0").Test(err)
	slot0Employee, err := GetEmployeeByID(slot0.EmployeeID, cy)
	tt.Describe("Getting employee for slot0").Test(err)

	var a0 testModel.Appointment
	a0_creation_error := a0.Create(200, ct.X_Auth_Token, nil, &slot0.StartTimeRFC3339, slot0.TimeZone, slot0Branch, slot0Employee, service, cy, ct)
	tt.Describe("Creating appointment a[0]").Test(a0_creation_error)

	Appointments = append(Appointments, &a0)

	// Test Case 1: Successful appointment creation
	slot1, err := service.FindValidRandomAppointmentSlot()
	tt.Describe("Finding valid appointment slot for service - a[1]").Test(err)

	slot1Branch, err := GetBranchByID(slot1.BranchID, cy)
	tt.Describe("Getting branch for slot1").Test(err)
	slot1Employee, err := GetEmployeeByID(slot1.EmployeeID, cy)
	tt.Describe("Getting employee for slot1").Test(err)

	var a1 testModel.Appointment
	a1_creation_error := a1.Create(200, ct.X_Auth_Token, nil, &slot1.StartTimeRFC3339, slot1.TimeZone, slot1Branch, slot1Employee, service, cy, ct)
	tt.Describe("Creating appointment a[1]").Test(a1_creation_error)

	Appointments = append(Appointments, &a1)

	// Test Case 2: Equal conflicting appointment creation for employee at slot 0

	if Appointments[0].Created.ID == uuid.Nil {
		t.Fatalf("Setup failed: a[0] is nil, cannot test conflict")
	}

	var a2 testModel.Appointment
	a2_creation_error := a2.Create(400, slot0Employee.X_Auth_Token, nil, &slot0.StartTimeRFC3339, slot0.TimeZone, slot0Branch, slot0Employee, service, cy, ct)
	tt.Describe("Creating conflicting appointment a[2] with employee token").Test(a2_creation_error)

}

// func Test_Appointment(t *testing.T) {
// 	server := src.NewServer().Run("parallel")
// 	defer server.Shutdown()
// 	tt := handlerT.NewTestErrorHandler(t)

// 	appEnv := os.Getenv("APP_ENV")
// 	if appEnv != "test" {
// 		t.Fatal("APP_ENV is not set to 'test'. Aborting tests to prevent data loss.")
// 	}

// 	ct := &testModel.Client{}
// 	cy := &testModel.Company{}

// 	tt.Describe("Client creation").Test(ct.Set())
// 	tt.Describe("Company setup").Test(cy.Set())

// 	baseEmployee := cy.Employees[1]
// 	Appointments := []*testModel.Appointment{}

// 	sp_location, err := time.LoadLocation("America/Sao_Paulo")
// 	if err != nil {
// 		t.Fatalf("Failed to load time zone: %v", err)
// 	}

// 	branchCache := &map[string]*model.Branch{}
// 	serviceCache := &map[string]*model.Service{}

// 	for _, b := range baseEmployee.Created.Branches {
// 		var branch model.Branch
// 		if err := handlerT.NewHttpClient().
// 			Header(namespace.HeadersKey.Company, baseEmployee.Company.Created.ID.String()).
// 			Header(namespace.HeadersKey.Auth, baseEmployee.Company.Owner.X_Auth_Token).
// 			Method("GET").
// 			URL("/branch/" + b.ID.String()).
// 			ExpectedStatus(200).
// 			Send(nil).
// 			ParseResponse(&branch).Error; err != nil {
// 			tt.Describe("Fetching branch for employee appointment setup").Test(err)
// 		}
// 		(*branchCache)[b.ID.String()] = &branch
// 	}

// 	for _, s := range baseEmployee.Created.Services {
// 		var service model.Service
// 		if err := handlerT.NewHttpClient().
// 			Header(namespace.HeadersKey.Company, baseEmployee.Company.Created.ID.String()).
// 			Header(namespace.HeadersKey.Auth, baseEmployee.Company.Owner.X_Auth_Token).
// 			Method("GET").
// 			URL("/service/" + s.ID.String()).
// 			ExpectedStatus(200).
// 			Send(nil).
// 			ParseResponse(&service).Error; err != nil {
// 			tt.Describe("Fetching service for employee appointment setup").Test(err)
// 		}
// 		(*serviceCache)[s.ID.String()] = &service
// 	}

// 	// --- Test Case 0 ---
// 	var a0 testModel.Appointment
// 	slot0, found0, err := a0.FindValidAppointmentSlot(baseEmployee, sp_location, branchCache, serviceCache)
// 	tt.Describe("Finding valid appointment slot for base employee - slot0").Test(err)
// 	if !found0 {
// 		t.Logf("Employee Work Schedule: %+v", baseEmployee.Created.WorkSchedule)
// 		t.Fatalf("Setup failed: No valid appointment slot for test case a[0]")
// 	}

// 	branch0, err := utilsT.GetBranchByID(cy, slot0.BranchID)
// 	tt.Describe("Getting branch for slot0").Test(err)
// 	service0, err := utilsT.GetServiceByID(cy, slot0.ServiceID)
// 	tt.Describe("Getting service for slot0").Test(err)

// 	tt.Describe("Creating appointment a[0]").Test(
// 		a0.Create(200, ct.X_Auth_Token, nil, &slot0.StartTimeRFC3339, slot0.TimeZone, branch0, baseEmployee, service0, cy, ct),
// 	)

// 	Appointments = append(Appointments, &a0)

// 	// --- Test Case 1 ---
// 	var a1 testModel.Appointment
// 	slot1, found1, err := a1.FindValidAppointmentSlot(baseEmployee, sp_location, branchCache, serviceCache)
// 	tt.Describe("Finding valid appointment slot for base employee - slot1").Test(err)
// 	if !found1 {
// 		t.Logf("Employee Work Schedule: %+v", baseEmployee.Created.WorkSchedule)
// 		t.Fatalf("Setup failed: No valid appointment slot for test case a[1]")
// 	}

// 	branch1, err := utilsT.GetBranchByID(cy, slot1.BranchID)
// 	tt.Describe("Getting branch for slot1").Test(err)
// 	service1, err := utilsT.GetServiceByID(cy, slot1.ServiceID)
// 	tt.Describe("Getting service for slot1").Test(err)
// 	tt.Describe("Creating appointment a[1]").Test(
// 		a1.Create(200, ct.X_Auth_Token, nil, &slot1.StartTimeRFC3339, slot1.TimeZone, branch1, baseEmployee, service1, cy, ct),
// 	)
// 	Appointments = append(Appointments, &a1)

// 	// --- Test Case 2 ---
// 	var a2 testModel.Appointment
// 	slot2, found2, err := a2.FindValidAppointmentSlot(baseEmployee, sp_location, branchCache, serviceCache)
// 	tt.Describe("Finding valid appointment slot for base employee - slot2").Test(err)
// 	if !found2 {
// 		t.Logf("Employee Work Schedule: %+v", baseEmployee.Created.WorkSchedule)
// 		t.Fatalf("Setup failed: No valid appointment slot for test case a[2]")
// 	}

// 	branch2, err := utilsT.GetBranchByID(cy, slot2.BranchID)
// 	tt.Describe("Getting branch for slot2").Test(err)
// 	service2, err := utilsT.GetServiceByID(cy, slot2.ServiceID)
// 	tt.Describe("Getting service for slot2").Test(err)
// 	tt.Describe("Creating appointment a[2]").Test(
// 		a2.Create(200, cy.Owner.X_Auth_Token, nil, &slot2.StartTimeRFC3339, slot2.TimeZone, branch2, baseEmployee, service2, cy, ct),
// 	)
// 	Appointments = append(Appointments, &a2)

// 	// --- Test Case 3 ---
// 	if Appointments[0].Created.ID == uuid.Nil {
// 		t.Fatalf("Setup failed: a[0] is nil, cannot test conflict")
// 	}

// 	ct1 := &testModel.Client{}
// 	tt.Describe("Client creation for conflicting appointments test").Test(ct1.Set())

// 	var a3 testModel.Appointment
// 	tt.Describe("Creating conflicting appointment a[3] with client token").Test(
// 		a3.Create(400, ct1.X_Auth_Token, nil, &slot0.StartTimeRFC3339, slot0.TimeZone, branch0, baseEmployee, service0, cy, ct1),
// 	)

// 	tt.Describe("Creating conflicting appointment a[3] with employee token").Test(
// 		a3.Create(400, baseEmployee.X_Auth_Token, nil, &slot0.StartTimeRFC3339, slot0.TimeZone, branch0, baseEmployee, service0, cy, ct1),
// 	)
// }
