package e2e_test

import (
	"agenda-kaki-go/core"
	handlerT "agenda-kaki-go/core/test/handlers"
	modelT "agenda-kaki-go/core/test/models"
	utilsT "agenda-kaki-go/core/test/utils"
	"testing"
)

type permissions_test struct {
	server             *core.Server
	company1           *modelT.Company
	company2           *modelT.Company
	client1            *modelT.Client
	client2            *modelT.Client
	company1_employee1 *modelT.Employee
	company1_employee2 *modelT.Employee
	company2_employee1 *modelT.Employee
	company2_employee2 *modelT.Employee
}

var permissions_test_instance *permissions_test

func load_permissions_test_instance(t *testing.T) {
	if permissions_test_instance == nil {
		Test_Setup_Permissions_Instance(t)
	}
}
func Test_Setup_Permissions_Instance(t *testing.T) {
	tt := handlerT.NewTestErrorHandler(t)

	permissions_test_instance = &permissions_test{}
	permissions_test_instance.server = core.NewServer().Run("parallel")
	if permissions_test_instance.server == nil {
		t.Fatal("Failed to start the server")
	}

	companies, err := utilsT.CreateCompaniesRandomly(2)
	tt.Describe("Creating companies randomly").Test(err)

	permissions_test_instance.company1 = companies[0]
	permissions_test_instance.company2 = companies[1]

	permissions_test_instance.client1 = &modelT.Client{}
	tt.Describe("Setting up client 1").Test(permissions_test_instance.client1.Set())

	permissions_test_instance.client2 = &modelT.Client{}
	tt.Describe("Setting up client 2").Test(permissions_test_instance.client2.Set())

	permissions_test_instance.company1_employee1 = permissions_test_instance.company1.Employees[1]
	permissions_test_instance.company1_employee2 = permissions_test_instance.company1.Employees[2]
	permissions_test_instance.company2_employee1 = permissions_test_instance.company2.Employees[1]
	permissions_test_instance.company2_employee2 = permissions_test_instance.company2.Employees[2]
}

func Test_Owner_x_Appointments(t *testing.T) {
	load_permissions_test_instance(t)
	tt := handlerT.NewTestErrorHandler(t)

	client1 := permissions_test_instance.client1
	client2 := permissions_test_instance.client2
	company1 := permissions_test_instance.company1
	company2 := permissions_test_instance.company2
	company1_owner := company1.Owner
	company2_owner := company2.Owner
	company1_employee1 := company1.Employees[1]
	company1_employee2 := company1.Employees[2]
	company2_employee1 := company2.Employees[1]
	company2_employee2 := company2.Employees[2]

	tt.Describe("Company1 Owner creates appointment for employee1 and client1").
		Test(utilsT.CreateAppointmentRandomly(200, company1, client1, company1_employee1, company1_owner.X_Auth_Token, company1.Created.ID.String(), nil))

	tt.Describe("Company1 Owner creates appointment for employee1 and client2").
		Test(utilsT.CreateAppointmentRandomly(200, company1, client2, company1_employee1, company1_owner.X_Auth_Token, company1.Created.ID.String(), nil))

	tt.Describe("Company1 Owner creates appointment for employee2 and client1").
		Test(utilsT.CreateAppointmentRandomly(200, company1, client1, company1_employee2, company1_owner.X_Auth_Token, company1.Created.ID.String(), nil))

	tt.Describe("Company1 Owner creates appointment for employee2 and client2").
		Test(utilsT.CreateAppointmentRandomly(200, company1, client2, company1_employee2, company1_owner.X_Auth_Token, company1.Created.ID.String(), nil))

	tt.Describe("Company2 Owner creates appointment for employee1 and client1").
		Test(utilsT.CreateAppointmentRandomly(200, company2, client1, company2_employee1, company2_owner.X_Auth_Token, company2.Created.ID.String(), nil))

	tt.Describe("Company2 Owner creates appointment for employee1 and client2").
		Test(utilsT.CreateAppointmentRandomly(200, company2, client2, company2_employee1, company2_owner.X_Auth_Token, company2.Created.ID.String(), nil))

	tt.Describe("Company2 Owner creates appointment for employee2 and client1").
		Test(utilsT.CreateAppointmentRandomly(200, company2, client1, company2_employee2, company2_owner.X_Auth_Token, company2.Created.ID.String(), nil))

	tt.Describe("Company2 Owner creates appointment for employee2 and client2").
		Test(utilsT.CreateAppointmentRandomly(200, company2, client2, company2_employee2, company2_owner.X_Auth_Token, company2.Created.ID.String(), nil))

	tt.Describe("Company1 Owner creates appointment for employee1 and client2 but using company2 ID").
		Test(utilsT.CreateAppointmentRandomly(403, company2, client2, company1_employee1, company1_owner.X_Auth_Token, company2.Created.ID.String(), nil))

	tt.Describe("Company2 Owner creates appointment for employee2 and client2 but using company1 ID").
		Test(utilsT.CreateAppointmentRandomly(403, company1, client2, company2_employee2, company2_owner.X_Auth_Token, company1.Created.ID.String(), nil))

	tt.Describe("Company1 Owner gets appointment of employee1 from company1").
		Test(utilsT.GetAppointment(200, company1_employee1.Created.Appointments[0].ID.String(), company1.Created.ID.String(), company1_owner.X_Auth_Token, nil))

	tt.Describe("Company1 Owner gets appointment of employee1 from company1 but using company2 ID").
		Test(utilsT.GetAppointment(403, company1_employee1.Created.Appointments[0].ID.String(), company2.Created.ID.String(), company1_owner.X_Auth_Token, nil))

	tt.Describe("Company1 Owner gets appointment of employee1 from company2").
		Test(utilsT.GetAppointment(403, company2_employee1.Created.Appointments[0].ID.String(), company1.Created.ID.String(), company1_owner.X_Auth_Token, nil))

	tt.Describe("Company2 Owner gets appointment of employee2 from company2 but using company1 ID").
		Test(utilsT.GetAppointment(403, company2_employee2.Created.Appointments[0].ID.String(), company2.Created.ID.String(), company1_owner.X_Auth_Token, nil))

	tt.Describe("Company2 Owner gets appointment of employee2 from company1").
		Test(utilsT.GetAppointment(403, company1_employee2.Created.Appointments[0].ID.String(), company2.Created.ID.String(), company2_owner.X_Auth_Token, nil))

	tt.Describe("Company1 Owner reschedules appointment of employee1 from company1").
		Test(utilsT.RescheduleAppointmentRandomly(200, company1_employee1, company1, company1_employee1.Created.Appointments[0].ID.String(), company1_owner.X_Auth_Token))

	tt.Describe("Company1 Owner reschedules appointment of employee1 from company2").
		Test(utilsT.RescheduleAppointmentRandomly(403, company1_employee1, company1, company2_employee1.Created.Appointments[0].ID.String(), company1_owner.X_Auth_Token))

	tt.Describe("Company1 Owner reschedules appointment of employee1 from company1 but using company2 ID").
		Test(utilsT.RescheduleAppointmentRandomly(403, company1_employee1, company2, company1_employee1.Created.Appointments[0].ID.String(), company1_owner.X_Auth_Token))

	tt.Describe("Company1 Owner cancels appointment of employee1 from company1").
		Test(utilsT.CancelAppointment(200, company1_employee1.Created.Appointments[0].ID.String(), company1.Created.ID.String(), company1_owner.X_Auth_Token))

	tt.Describe("Company1 Owner cancels appointment of employee1 from company2").
		Test(utilsT.CancelAppointment(403, company2_employee1.Created.Appointments[0].ID.String(), company1.Created.ID.String(), company1_owner.X_Auth_Token))

	tt.Describe("Company1 Owner cancels appointment of employee1 from company1 but using company2 ID").
		Test(utilsT.CancelAppointment(403, company1_employee1.Created.Appointments[0].ID.String(), company2.Created.ID.String(), company1_owner.X_Auth_Token))
}

func Test_Employee_x_Appointments(t *testing.T) {
	load_permissions_test_instance(t)
	tt := handlerT.NewTestErrorHandler(t)

	client1 := permissions_test_instance.client1
	client2 := permissions_test_instance.client2
	company1 := permissions_test_instance.company1
	company2 := permissions_test_instance.company2
	company1_employee1 := company1.Employees[1]
	company1_employee2 := company1.Employees[2]
	company2_employee1 := company2.Employees[1]
	company2_employee2 := company2.Employees[2]

	tt.Describe("Employee1 from company1 creates appointment for client1").
		Test(utilsT.CreateAppointmentRandomly(200, company1, client1, company1_employee1, company1_employee1.X_Auth_Token, company1.Created.ID.String(), nil))

	tt.Describe("Employee1 from company1 creates appointment for client2").
		Test(utilsT.CreateAppointmentRandomly(200, company1, client2, company1_employee1, company1_employee1.X_Auth_Token, company1.Created.ID.String(), nil))

	tt.Describe("Employee2 from company2 creates appointment for client1").
		Test(utilsT.CreateAppointmentRandomly(200, company2, client1, company2_employee2, company2_employee2.X_Auth_Token, company2.Created.ID.String(), nil))

	tt.Describe("Employee2 from company2 creates appointment for client2").
		Test(utilsT.CreateAppointmentRandomly(200, company2, client2, company2_employee2, company2_employee2.X_Auth_Token, company2.Created.ID.String(), nil))

	tt.Describe("Employee1 from company1 tries to create for another employee in company1").
		Test(utilsT.CreateAppointmentRandomly(403, company1, client1, company1_employee2, company1_employee1.X_Auth_Token, company1.Created.ID.String(), nil))

	tt.Describe("Employee1 from company1 tries to create for employee2 in company2").
		Test(utilsT.CreateAppointmentRandomly(403, company2, client1, company2_employee2, company1_employee1.X_Auth_Token, company2.Created.ID.String(), nil))

	tt.Describe("Employee1 from company1 tries to get appointment from employee2 in company2").
		Test(utilsT.GetAppointment(403, company2_employee2.Created.Appointments[0].ID.String(), company2.Created.ID.String(), company1_employee1.X_Auth_Token, nil))

	tt.Describe("Employee1 from company2 tries to get appointment from employee2 in company2").
		Test(utilsT.GetAppointment(403, company2_employee2.Created.Appointments[0].ID.String(), company1.Created.ID.String(), company2_employee1.X_Auth_Token, nil))

	tt.Describe("Employee1 from company1 gets their own appointment").
		Test(utilsT.GetAppointment(200, company1_employee1.Created.Appointments[0].ID.String(), company1.Created.ID.String(), company1_employee1.X_Auth_Token, nil))

	tt.Describe("Employee1 from company1 reschedules their own appointment").
		Test(utilsT.RescheduleAppointmentRandomly(200, company1_employee1, company1, company1_employee1.Created.Appointments[0].ID.String(), company1_employee1.X_Auth_Token))

	tt.Describe("Employee1 from company1 tries to reschedule another employee's appointment in company1").
		Test(utilsT.RescheduleAppointmentRandomly(403, company1_employee2, company1, company1_employee2.Created.Appointments[0].ID.String(), company1_employee1.X_Auth_Token))

	tt.Describe("Employee1 from company1 tries to reschedule appointment from company2").
		Test(utilsT.RescheduleAppointmentRandomly(403, company2_employee2, company1, company2_employee2.Created.Appointments[0].ID.String(), company1_employee1.X_Auth_Token))

	tt.Describe("Employee1 from company1 cancels their own appointment").
		Test(utilsT.CancelAppointment(200, company1_employee1.Created.Appointments[0].ID.String(), company1.Created.ID.String(), company1_employee1.X_Auth_Token))

	tt.Describe("Employee1 from company1 tries to cancel another employee's appointment in company1").
		Test(utilsT.CancelAppointment(403, company1_employee2.Created.Appointments[0].ID.String(), company1.Created.ID.String(), company1_employee1.X_Auth_Token))

	tt.Describe("Employee1 from company1 tries to cancel appointment from company2").
		Test(utilsT.CancelAppointment(403, company2_employee2.Created.Appointments[0].ID.String(), company1.Created.ID.String(), company1_employee1.X_Auth_Token))
}

func Test_Shutdown_Permissions_Instance(t *testing.T) {
	if permissions_test_instance != nil {
		permissions_test_instance.server.Shutdown()
		permissions_test_instance = nil
	}
}
