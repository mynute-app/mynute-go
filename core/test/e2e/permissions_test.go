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
	if permissions_test_instance != nil && permissions_test_instance.server != nil {
		return
	}
	Test_Setup_Permissions_Instance(t)
}

func Test_Setup_Permissions_Instance(t *testing.T) {
	if permissions_test_instance != nil {
		return
	}
	permissions_test_instance = &permissions_test{}
	tt := handlerT.NewTestErrorHandler(t)

	companies, err := utilsT.CreateCompaniesRandomly(2)
	if err != nil {
		t.Fatal("Failed to create companies:", err)
	}

	permissions_test_instance.server = core.NewServer().Run("parallel")

	if permissions_test_instance.server == nil {
		t.Fatal("Failed to start the server")
	}

	permissions_test_instance.company1 = companies[0]
	permissions_test_instance.company2 = companies[1]

	company_1_employee_number := 4
	company_2_employee_number := 3
	company_1_branch_number := 3
	company_2_branch_number := 2
	company_1_service_number := 22
	company_2_service_number := 6
	permissions_test_instance.company1 = &modelT.Company{}
	tt.Test(permissions_test_instance.company1.SetupRandomized(company_1_employee_number, company_1_branch_number, company_1_service_number)) // owner, 4 employees, 3 branches, 22 services
	permissions_test_instance.company2 = &modelT.Company{}
	tt.Test(permissions_test_instance.company2.SetupRandomized(company_2_employee_number, company_2_branch_number, company_2_service_number)) // owner, 3 employees, 2 branches, 6 services
	permissions_test_instance.client1 = &modelT.Client{}
	tt.Test(permissions_test_instance.client1.Set())
	permissions_test_instance.client2 = &modelT.Client{}
	tt.Test(permissions_test_instance.client2.Set())

	permissions_test_instance.company1_employee1 = permissions_test_instance.company1.Employees[1]
	permissions_test_instance.company1_employee2 = permissions_test_instance.company1.Employees[2]
	permissions_test_instance.company2_employee1 = permissions_test_instance.company2.Employees[1]
	permissions_test_instance.company2_employee2 = permissions_test_instance.company2.Employees[2]

	if permissions_test_instance.company1_employee1.Auth_token == "" {
		t.Fatal("Company 1 Employee 1 auth token is missing")
	} else if permissions_test_instance.company1_employee2.Auth_token == "" {
		t.Fatal("Company 1 Employee 2 auth token is missing")
	} else if permissions_test_instance.company2_employee1.Auth_token == "" {
		t.Fatal("Company 2 Employee 1 auth token is missing")
	} else if permissions_test_instance.company2_employee2.Auth_token == "" {
		t.Fatal("Company 2 Employee 2 auth token is missing")
	}
}

func Test_Owner_x_Appointments(t *testing.T) {
	load_permissions_test_instance(t)
	tt := handlerT.NewTestErrorHandler(t)

	client1 := permissions_test_instance.client1
	client2 := permissions_test_instance.client2
	company1 := permissions_test_instance.company1
	company2 := permissions_test_instance.company2
	company1_owner := permissions_test_instance.company1.Owner
	company2_owner := permissions_test_instance.company2.Owner
	company1_employee1 := permissions_test_instance.company1.Employees[1]
	company1_employee2 := permissions_test_instance.company1.Employees[2]
	company2_employee1 := permissions_test_instance.company2.Employees[1]
	company2_employee2 := permissions_test_instance.company2.Employees[2]

	t.Log("--- Testing Owner x Appointment Interactions ---")

	t.Log("---> Company1 Owner creating appointment for company1_employee1 and client1 at company1 and company1.Created.ID.String() : POST /appointment => HTTP 200")
	tt.Test(utilsT.CreateAppointmentRandomly(200, company1, client1, company1_employee1, company1_owner.Auth_token, company1.Created.ID.String(), nil))
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company1 Owner creating appointment for company1_employee1 and client2 at company1 and company1.Created.ID.String() : POST /appointment => HTTP 200")
	tt.Test(utilsT.CreateAppointmentRandomly(200, company1, client2, company1_employee1, company1_owner.Auth_token, company1.Created.ID.String(), nil))
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company1 Owner creating appointment for company1_employee2 and client1 at company1 and company1.Created.ID.String() : POST /appointment => HTTP 200")
	tt.Test(utilsT.CreateAppointmentRandomly(200, company1, client1, company1_employee2, company1_owner.Auth_token, company1.Created.ID.String(), nil))
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company1 Owner creating appointment for company1_employee2 and client2 at company1 and company1.Created.ID.String() : POST /appointment => HTTP 200")
	tt.Test(utilsT.CreateAppointmentRandomly(200, company1, client2, company1_employee2, company1_owner.Auth_token, company1.Created.ID.String(), nil))
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company2 Owner creating appointment for company2_employee1 and client1 at company2 and company2.Created.ID.String() : POST /appointment => HTTP 200")
	tt.Test(utilsT.CreateAppointmentRandomly(200, company2, client1, company2_employee1, company2_owner.Auth_token, company2.Created.ID.String(), nil))
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company2 Owner creating appointment for company2_employee1 and client2 at company2 and company2.Created.ID.String() : POST /appointment => HTTP 200")
	tt.Test(utilsT.CreateAppointmentRandomly(200, company2, client2, company2_employee1, company2_owner.Auth_token, company2.Created.ID.String(), nil))
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company2 Owner creating appointment for company2_employee2 and client1 at company2 and company2.Created.ID.String() : POST /appointment => HTTP 200")
	tt.Test(utilsT.CreateAppointmentRandomly(200, company2, client1, company2_employee2, company2_owner.Auth_token, company2.Created.ID.String(), nil))
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company2 Owner creating appointment for company2_employee2 and client2 at company2 and company2.Created.ID.String() : POST /appointment => HTTP 200")
	tt.Test(utilsT.CreateAppointmentRandomly(200, company2, client2, company2_employee2, company2_owner.Auth_token, company2.Created.ID.String(), nil))
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company1 Owner trying to get an appointment of Employee 1 from Company 1 : GET /appointment/{id} => HTTP 200")
	tt.Test(utilsT.GetAppointment(200, company1_employee1.Created.Appointments[0].ID.String(), company1.Created.ID.String(), company1_owner.Auth_token, nil))
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company1 Owner trying to get an appointment of Employee 1 from Company 2 : GET /appointment/{id} => HTTP 403")
	tt.Test(utilsT.GetAppointment(403, company2_employee1.Created.Appointments[0].ID.String(), company1.Created.ID.String(), company1_owner.Auth_token, nil))
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company1 Owner trying to reschedule an appointment of Employee 1 from Company 1 : PATCH /appointment/{id} => HTTP 200")
	tt.Test(utilsT.RescheduleAppointmentRandomly(200, company1_employee1, company1, company1_employee1.Created.Appointments[0].ID.String(), company1_owner.Auth_token))
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company1 Owner trying to reschedule an appointment of Employee 1 from Company 2 : PATCH /appointment/{id} => HTTP 403")
	tt.Test(utilsT.RescheduleAppointmentRandomly(403, company1_employee1, company1, company2_employee1.Created.Appointments[0].ID.String(), company1_owner.Auth_token))
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company1 Owner trying to cancel an appointment of Employee 1 from Company 1 : DELETE /appointment/{id} => HTTP 200")
	tt.Test(utilsT.CancelAppointment(200, company1_employee1.Created.Appointments[0].ID.String(), company1.Created.ID.String(), company1_owner.Auth_token))
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company1 Owner trying to cancel an appointment of Employee 1 from Company 2 : DELETE /appointment/{id} => HTTP 403")
	tt.Test(utilsT.CancelAppointment(403, company2_employee1.Created.Appointments[0].ID.String(), company1.Created.ID.String(), company1_owner.Auth_token))
	t.Log("---------------------- x ----------------------")
}

func Test_Employee_x_Appointments(t *testing.T) {
	load_permissions_test_instance(t)

	client1 := permissions_test_instance.client1
	client2 := permissions_test_instance.client2
	company1 := permissions_test_instance.company1
	company2 := permissions_test_instance.company2
	company1_employee1 := permissions_test_instance.company1.Employees[1]
	company1_employee2 := permissions_test_instance.company1.Employees[2]
	company2_employee1 := permissions_test_instance.company2.Employees[1]
	company2_employee2 := permissions_test_instance.company2.Employees[2]

	t.Log("--- Testing Employee x Appointment Interactions ---")

	t.Log("---> Company1 Employee 1 creating appointment for client1 at company1 and company1.Created.ID.String() : POST /appointment => HTTP 200")
	utilsT.CreateAppointmentRandomly(200, company1, client1, company1_employee1, company1_employee1.Auth_token, company1.Created.ID.String(), nil)
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company1 Employee 1 creating appointment for client2 at company1 and company1.Created.ID.String() : POST /appointment => HTTP 200")
	utilsT.CreateAppointmentRandomly(200, company1, client2, company1_employee1, company1_employee1.Auth_token, company1.Created.ID.String(), nil)
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company2 Employee 2 creating appointment for client1 at company2 and company2.Created.ID.String() : POST /appointment => HTTP 200")
	utilsT.CreateAppointmentRandomly(200, company2, client1, company2_employee2, company2_employee2.Auth_token, company2.Created.ID.String(), nil)
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company2 Employee 2 creating appointment for client2 at company2 and company2.Created.ID.String() : POST /appointment => HTTP 200")
	utilsT.CreateAppointmentRandomly(200, company2, client2, company2_employee2, company2_employee2.Auth_token, company2.Created.ID.String(), nil)
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company 1 Employee 1 trying to create an appointment for Employee 2 from Company 1 : POST /appointment => HTTP 403")
	utilsT.CreateAppointmentRandomly(403, company1, client1, company1_employee2, company1_employee1.Auth_token, company1.Created.ID.String(), nil)
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company 1 Employee 1 trying to create an appointment for Employee 2 from Company 2 : POST /appointment => HTTP 403")
	utilsT.CreateAppointmentRandomly(403, company2, client1, company2_employee2, company1_employee1.Auth_token, company2.Created.ID.String(), nil)
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company 1 Employee 1 trying to get an appointment of Employee 2 from Company 2 : GET /appointment/{id} => HTTP 403")
	utilsT.GetAppointment(403, company2_employee2.Created.Appointments[0].ID.String(), company2.Created.ID.String(), company1_employee1.Auth_token, nil)
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company 2 Employee 1 trying to get an appointment of Employee 2 from Company 2 : GET /appointment/{id} => HTTP 403")
	utilsT.GetAppointment(403, company2_employee2.Created.Appointments[0].ID.String(), company1.Created.ID.String(), company2_employee1.Auth_token, nil)
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company 1 Employee 1 trying to get an appointment of Employee 1 from Company 1 : GET /appointment/{id} => HTTP 200")
	utilsT.GetAppointment(200, company1_employee1.Created.Appointments[0].ID.String(), company1.Created.ID.String(), company1_employee1.Auth_token, nil)
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company 1 Employee 1 trying to reschedule an appointment of Employee 1 from Company 1 : PATCH /appointment/{id} => HTTP 200")
	utilsT.RescheduleAppointmentRandomly(200, company1_employee1, company1, company1_employee1.Created.Appointments[0].ID.String(), company1_employee1.Auth_token)
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company 1 Employee 1 trying to reschedule an appointment of Employee 2 from Company 1 : PATCH /appointment/{id} => HTTP 403")
	utilsT.RescheduleAppointmentRandomly(403, company1_employee2, company1, company1_employee2.Created.Appointments[0].ID.String(), company1_employee1.Auth_token)
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company 1 Employee 1 trying to reschedule an appointment of Employee 2 from Company 2 : PATCH /appointment/{id} => HTTP 403")
	utilsT.RescheduleAppointmentRandomly(403, company2_employee2, company1, company2_employee2.Created.Appointments[0].ID.String(), company1_employee1.Auth_token)
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company 1 Employee 1 trying to cancel an appointment of Employee 1 from Company 1 : DELETE /appointment/{id} => HTTP 200")
	utilsT.CancelAppointment(200, company1_employee1.Created.Appointments[0].ID.String(), company1.Created.ID.String(), company1_employee1.Auth_token)
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company 1 Employee 1 trying to cancel an appointment of Employee 2 from Company 1 : DELETE /appointment/{id} => HTTP 403")
	utilsT.CancelAppointment(403, company1_employee2.Created.Appointments[0].ID.String(), company1.Created.ID.String(), company1_employee1.Auth_token)
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company 1 Employee 1 trying to cancel an appointment of Employee 2 from Company 2 : DELETE /appointment/{id} => HTTP 403")
	utilsT.CancelAppointment(403, company2_employee2.Created.Appointments[0].ID.String(), company1.Created.ID.String(), company1_employee1.Auth_token)
	t.Log("---------------------- x ----------------------")
}

func Test_Cleanup_Environment(t *testing.T) {
	if permissions_test_instance == nil || permissions_test_instance.server == nil {
		Test_Setup_Permissions_Instance(t)
	}
	t.Log("--- Cleaning up environment ---")
	permissions_test_instance.server.Shutdown()
	t.Log("Environment cleaned up successfully")
}
