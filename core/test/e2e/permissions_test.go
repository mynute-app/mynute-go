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
	c1ID := company1.Created.ID.String()
	company2 := permissions_test_instance.company2
	c2ID := company2.Created.ID.String()
	company1_owner := company1.Owner
	company2_owner := company2.Owner
	company1_employee1 := company1.Employees[1]
	company1_employee2 := company1.Employees[2]
	company2_employee1 := company2.Employees[1]
	company2_employee2 := company2.Employees[2]

	var a_cy1_e1_ct1 modelT.Appointment
	tt.Describe("Company1 Owner creates appointment for employee1 and client1").
		Test(a_cy1_e1_ct1.CreateRandomly(200, company1, client1, company1_employee1, company1_owner.X_Auth_Token, c1ID))

	var a_cy1_e1_ct2 modelT.Appointment
	tt.Describe("Company1 Owner creates appointment for employee1 and client2").
		Test(a_cy1_e1_ct2.CreateRandomly(200, company1, client2, company1_employee1, company1_owner.X_Auth_Token, c1ID))

	var a_cy1_e2_ct1 modelT.Appointment
	tt.Describe("Company1 Owner creates appointment for employee2 and client1").
		Test(a_cy1_e2_ct1.CreateRandomly(200, company1, client1, company1_employee2, company1_owner.X_Auth_Token, c1ID))

	var a_cy1_e2_ct2 modelT.Appointment
	tt.Describe("Company1 Owner creates appointment for employee2 and client2").
		Test(a_cy1_e2_ct2.CreateRandomly(200, company1, client2, company1_employee2, company1_owner.X_Auth_Token, c1ID))

	var a_cy2_e1_ct1 modelT.Appointment
	tt.Describe("Company2 Owner creates appointment for employee1 and client1").
		Test(a_cy2_e1_ct1.CreateRandomly(200, company2, client1, company2_employee1, company2_owner.X_Auth_Token, c2ID))

	var a_cy2_e1_ct2 modelT.Appointment
	tt.Describe("Company2 Owner creates appointment for employee1 and client2").
		Test(a_cy2_e1_ct2.CreateRandomly(200, company2, client2, company2_employee1, company2_owner.X_Auth_Token, c2ID))

	var a_cy2_e2_ct1 modelT.Appointment
	tt.Describe("Company2 Owner creates appointment for employee2 and client1").
		Test(a_cy2_e2_ct1.CreateRandomly(200, company2, client1, company2_employee2, company2_owner.X_Auth_Token, c2ID))

	var a_cy2_e2_ct2 modelT.Appointment
	tt.Describe("Company2 Owner creates appointment for employee2 and client2").
		Test(a_cy2_e2_ct2.CreateRandomly(200, company2, client2, company2_employee2, company2_owner.X_Auth_Token, c2ID))

	var a_fail modelT.Appointment
	tt.Describe("Company1 owner unauthorized to create appointment for employee1 and client2 using company2 ID").
		Test(a_fail.CreateRandomly(403, company2, client2, company1_employee1, company1_owner.X_Auth_Token, c2ID))

	tt.Describe("Company2 owner unauthorized to create appointment for employee2 and client2 using company1 ID").
		Test(a_fail.CreateRandomly(403, company1, client2, company2_employee2, company2_owner.X_Auth_Token, c1ID))

	tt.Describe("Company1 Owner gets appointment of employee1 from company1").
		Test(a_cy1_e1_ct1.GetById(200, company1.Owner.X_Auth_Token, nil))

	tt.Describe("Company1 Owner unauthorized to get appointment of employee1 from company1 using company2 ID as header").
		Test(a_cy1_e1_ct1.GetById(403, company1.Owner.X_Auth_Token, &c2ID))

	tt.Describe("Company1 Owner unauthorized to get appointment of employee1 from company2").
		Test(a_cy2_e1_ct1.GetById(403, company1.Owner.X_Auth_Token, nil))

	tt.Describe("Company2 Owner unauthorized to get appointment of employee2 from company2 using company1 ID as header").
		Test(a_cy2_e2_ct1.GetById(403, company2.Owner.X_Auth_Token, &c1ID))

	tt.Describe("Company2 Owner unauthorized to get appointment of employee2 from company1").
		Test(a_cy1_e2_ct1.GetById(403, company2.Owner.X_Auth_Token, nil))

	tt.Describe("Company1 Owner reschedules appointment of employee1 from company1").
		Test(a_cy1_e1_ct1.RescheduleRandomly(200, company1.Owner.X_Auth_Token, nil))

	tt.Describe("Company1 Owner unauthorized to reschedule appointment of employee1 from company2").
		Test(a_cy2_e1_ct1.RescheduleRandomly(403, company1.Owner.X_Auth_Token, nil))

	tt.Describe("Company1 Owner unauthorized to reschedule appointment of employee1 from company1 using company2 ID").
		Test(a_cy1_e1_ct1.RescheduleRandomly(403, company1.Owner.X_Auth_Token, &c2ID))

	tt.Describe("Company1 Owner unauthorized to cancel appointment of employee1 from company2").
		Test(a_cy2_e1_ct1.Cancel(403, company1.Owner.X_Auth_Token, nil))

	tt.Describe("Company1 Owner unauthorized to cancel appointment of employee1 from company1 using company2 ID").
		Test(a_cy2_e1_ct1.Cancel(403, company1.Owner.X_Auth_Token, &c2ID))

	tt.Describe("Company1 Owner cancels appointment of employee1 from company1 with client 2").
		Test(a_cy1_e1_ct2.Cancel(200, company1.Owner.X_Auth_Token, nil))

	tt.Describe("Company1 Owner creates appointment for employee1 and client1 in company1").
		Test(a_cy1_e1_ct2.CreateRandomly(200, company1, client1, company1_employee1, company1_owner.X_Auth_Token, c1ID))

}

func Test_Employee_x_Appointments(t *testing.T) {
	load_permissions_test_instance(t)
	tt := handlerT.NewTestErrorHandler(t)

	client1 := permissions_test_instance.client1
	client2 := permissions_test_instance.client2
	company1 := permissions_test_instance.company1
	c1ID := company1.Created.ID.String()
	company2 := permissions_test_instance.company2
	c2ID := company2.Created.ID.String()
	company1_employee1 := company1.Employees[1]
	company1_employee2 := company1.Employees[2]
	company2_employee1 := company2.Employees[1]
	company2_employee2 := company2.Employees[2]

	var a_cy1_e1_ct1 modelT.Appointment
	tt.Describe("Company1 Employee1 creates appointment for client1").
		Test(a_cy1_e1_ct1.CreateRandomly(200, company1, client1, company1_employee1, company1_employee1.X_Auth_Token, c1ID))

	var a_cy1_e1_ct2 modelT.Appointment
	tt.Describe("Company1 Employee1 creates appointment for client2").
		Test(a_cy1_e1_ct2.CreateRandomly(200, company1, client2, company1_employee1, company1_employee1.X_Auth_Token, c1ID))

	var a_cy1_e2_ct1 modelT.Appointment
	tt.Describe("Company1 Employee2 creates appointment for client1").
		Test(a_cy1_e2_ct1.CreateRandomly(200, company1, client1, company1_employee2, company1_employee2.X_Auth_Token, c1ID))

	var a_cy1_e2_ct2 modelT.Appointment
	tt.Describe("Company1 Employee2 creates appointment for client2").
		Test(a_cy1_e2_ct2.CreateRandomly(200, company1, client2, company1_employee2, company1_employee2.X_Auth_Token, c1ID))

	var a_cy2_e1_ct1 modelT.Appointment
	tt.Describe("Company2 Employee1 creates appointment for client1").
		Test(a_cy2_e1_ct1.CreateRandomly(200, company2, client1, company2_employee1, company2_employee1.X_Auth_Token, c2ID))

	var a_cy2_e1_ct2 modelT.Appointment
	tt.Describe("Company2 Employee1 creates appointment for client2").
		Test(a_cy2_e1_ct2.CreateRandomly(200, company2, client2, company2_employee1, company2_employee1.X_Auth_Token, c2ID))

	var a_cy2_e2_ct1 modelT.Appointment
	tt.Describe("Company2 Employee2 creates appointment for client1").
		Test(a_cy2_e2_ct1.CreateRandomly(200, company2, client1, company2_employee2, company2_employee2.X_Auth_Token, c2ID))

	var a_cy2_e2_ct2 modelT.Appointment
	tt.Describe("Company2 Employee2 creates appointment for client2").
		Test(a_cy2_e2_ct2.CreateRandomly(200, company2, client2, company2_employee2, company2_employee2.X_Auth_Token, c2ID))

	var a_fail modelT.Appointment
	tt.Describe("Company1 Employee1 unauthorized to create appointment for employee2 in company1").
		Test(a_fail.CreateRandomly(403, company1, client1, company1_employee2, company1_employee1.X_Auth_Token, c1ID))

	tt.Describe("Company1 Employee1 unauthorized to create appointment for employee2 in company2").
		Test(a_fail.CreateRandomly(403, company2, client1, company2_employee2, company1_employee1.X_Auth_Token, c2ID))

	tt.Describe("Company1 Employee1 gets their own appointment").
		Test(a_cy1_e1_ct1.GetById(200, company1_employee1.X_Auth_Token, nil))

	tt.Describe("Company1 Employee1 unauthorized to get their own appointment using company2 ID").
		Test(a_cy1_e1_ct1.GetById(403, company1_employee1.X_Auth_Token, &c2ID))

	tt.Describe("Company1 Employee1 unauthorized to get appointment from employee2 in company1").
		Test(a_cy1_e2_ct1.GetById(403, company1_employee1.X_Auth_Token, nil))

	tt.Describe("Company1 Employee1 unauthorized to get appointment from employee2 in company2").
		Test(a_cy2_e2_ct1.GetById(403, company1_employee1.X_Auth_Token, nil))

	tt.Describe("Company1 Employee1 unauthorized to get appointment from employee2 in company1 using company2 ID").
		Test(a_cy1_e2_ct1.GetById(403, company1_employee1.X_Auth_Token, &c2ID))

	tt.Describe("Company1 Employee1 unauthorized to get appointment from employee2 in company2 using company1 ID").
		Test(a_cy2_e2_ct1.GetById(403, company1_employee1.X_Auth_Token, &c1ID))

	tt.Describe("Company1 Employee1 reschedules their own appointment").
		Test(a_cy1_e1_ct1.RescheduleRandomly(200, company1_employee1.X_Auth_Token, nil))

	tt.Describe("Company1 Employee1 unauthorized to reschedule their own appointment using company2 ID").
		Test(a_cy1_e1_ct1.RescheduleRandomly(403, company1_employee1.X_Auth_Token, &c2ID))

	tt.Describe("Company1 Employee1 unauthorized to reschedule appointment of employee2 in company1").
		Test(a_cy1_e2_ct1.RescheduleRandomly(403, company1_employee1.X_Auth_Token, nil))

	tt.Describe("Company1 Employee1 unauthorized to reschedule appointment of employee2 in company1 using company2 ID").
		Test(a_cy1_e2_ct1.RescheduleRandomly(403, company2_employee2.X_Auth_Token, &c2ID))

	tt.Describe("Company1 Employee1 unauthorized to cancel appointment of employee2 in company1").
		Test(a_cy1_e2_ct1.Cancel(403, company1_employee1.X_Auth_Token, nil))

	tt.Describe("Company1 Employee1 unauthorized to cancel appointment of employee2 in company2").
		Test(a_cy2_e2_ct1.Cancel(403, company1_employee1.X_Auth_Token, nil))

	tt.Describe("Company1 Employee1 unauthorized to cancel appointment of employee2 in company1 using company2 ID").
		Test(a_cy1_e2_ct1.Cancel(403, company2_employee2.X_Auth_Token, &c2ID))

	tt.Describe("Company1 Employee1 unauthorized to cancel appointment of employee2 in company2 using company1 ID").
		Test(a_cy2_e2_ct1.Cancel(403, company1_employee1.X_Auth_Token, &c1ID))

	tt.Describe("Company1 Employee1 unauthorized to cancel appointment of employee2 in company2 using company2 ID").
		Test(a_cy2_e2_ct1.Cancel(403, company1_employee1.X_Auth_Token, &c2ID))

	tt.Describe("Company1 Employee1 cancels their own appointment").
		Test(a_cy1_e1_ct1.Cancel(200, company1_employee1.X_Auth_Token, nil))

	tt.Describe("Company1 Employee1 creates appointment for client1 in company1").
		Test(a_cy1_e1_ct1.CreateRandomly(200, company1, client1, company1_employee1, company1_employee1.X_Auth_Token, c1ID))
}

func Test_Shutdown_Permissions_Instance(t *testing.T) {
	if permissions_test_instance != nil {
		permissions_test_instance.server.Shutdown()
		permissions_test_instance = nil
	}
}

// All tests working 100% :D