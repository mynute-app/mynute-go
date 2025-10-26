package e2e_test

import (
	"fmt"
	"mynute-go/core"
	"mynute-go/core/src/lib"
	"mynute-go/test/src/handler"
	testModel "mynute-go/test/src/model"
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

	TimeZone := "America/Sao_Paulo"

	ct := &testModel.Client{}
	tt.Describe("Client creation").Test(ct.Set())

	var cys []*testModel.Company
	companiesToCreate := 2

	for i := range companiesToCreate {
		_ = &testModel.Company{}
		min := 2
		max := 4
		empN := lib.GenerateRandomIntFromRange(min, max)
		branchN := lib.GenerateRandomIntFromRange(min, max)
		serviceN := lib.GenerateRandomIntFromRange(min, max)
		cy := &testModel.Company{}
		tt.Describe(fmt.Sprintf("Company Random Setup [%d]", i)).Test(cy.CreateCompanyRandomly(empN, branchN, serviceN))
		cys = append(cys, cy)
	}

	Appointments := []*testModel.Appointment{}

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

	// Test Case 0: Creating multiple appointments for a single customer.

	appointmentsToCreate := 50
	client_public_id := ct.Created.ID.String()
	for _, cy := range cys {
		for _, service := range cy.Services {
			for i := range appointmentsToCreate {
				slot, err := service.FindValidRandomAppointmentSlot(TimeZone, &client_public_id)
				tt.Describe(fmt.Sprintf("Finding valid appointment slot for service - a[%d]", i)).Test(err)

				slotBranch, err := GetBranchByID(slot.BranchID, cy)
				tt.Describe(fmt.Sprintf("Getting branch for slot[%d]", i)).Test(err)
				slotEmployee, err := GetEmployeeByID(slot.EmployeeID, cy)
				tt.Describe(fmt.Sprintf("Getting employee for slot[%d]", i)).Test(err)

				var appointment testModel.Appointment
				appointment_creation_error := appointment.Create(200, ct.X_Auth_Token, nil, &slot.StartTimeRFC3339, slot.TimeZone, slotBranch, slotEmployee, service, cy, ct)
				tt.Describe(fmt.Sprintf("Creating appointment a[%d]", i)).Test(appointment_creation_error)

				Appointments = append(Appointments, &appointment)
			}
		}
	}

	// Test Case 2: Equal conflicting appointment creation for employee at slot 0

	if Appointments[0].Created.ID == uuid.Nil {
		t.Fatalf("Setup failed: a[0] is nil, cannot test conflict")
	}

	slot0 := Appointments[0]
	slot0StartTimeRFC3339 := slot0.Created.StartTime.Format(time.RFC3339)

	var a2 testModel.Appointment
	a2_creation_error := a2.Create(400, slot0.Employee.X_Auth_Token, nil, &slot0StartTimeRFC3339, slot0.Created.TimeZone, slot0.Branch, slot0.Employee, slot0.Service, slot0.Company, ct)
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
