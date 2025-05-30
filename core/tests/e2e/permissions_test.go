package e2e_test

import (
	"agenda-kaki-go/core"
	"agenda-kaki-go/core/config/namespace"
	handler "agenda-kaki-go/core/tests/handlers"
	"testing"
	"time"
)

type permissions_test struct {
	server             *core.Server
	company1           *Company
	company2           *Company
	client1            *Client
	client2            *Client
	company1_employee1 *Employee
	company1_employee2 *Employee
	company2_employee1 *Employee
	company2_employee2 *Employee
}

var permissions_test_instance *permissions_test

func Test_Start_Server(t *testing.T) {
	permissions_test_instance = &permissions_test{
		server: core.NewServer().Run("test"),
	}
	if permissions_test_instance.server == nil {
		t.Fatal("Failed to start the server")
	}
	t.Log("Server started successfully")
}

func Test_Setup_Environment(t *testing.T) {
	if permissions_test_instance == nil || permissions_test_instance.server == nil {
		Test_Start_Server(t)
	}
	company_1_employee_number := 4
	company_2_employee_number := 3
	company_1_branch_number := 3
	company_2_branch_number := 2
	company_1_service_number := 22
	company_2_service_number := 6
	permissions_test_instance.company1 = &Company{}
	permissions_test_instance.company1.SetupRandomized(t, company_1_employee_number, company_1_branch_number, company_1_service_number) // owner, 4 employees, 3 branches, 22 services
	permissions_test_instance.company2 = &Company{}
	permissions_test_instance.company2.SetupRandomized(t, company_2_employee_number, company_2_branch_number, company_2_service_number) // owner, 3 employees, 2 branches, 6 services
	permissions_test_instance.client1 = &Client{}
	permissions_test_instance.client1.Set(t)
	permissions_test_instance.client2 = &Client{}
	permissions_test_instance.client2.Set(t)

	if permissions_test_instance.company1.owner.auth_token == "" {
		t.Fatal("Company 1 Owner auth token is missing")
	} else if permissions_test_instance.company2.owner.auth_token == "" {
		t.Fatal("Company 2 Owner auth token is missing")
	} else if len(permissions_test_instance.company1.employees) != company_1_employee_number+1 { // +1 for the owner
		t.Fatalf("Company 1 does not have enough employees, expected %d, got %d", company_1_employee_number, len(permissions_test_instance.company1.employees))
	} else if len(permissions_test_instance.company2.employees) != company_2_employee_number+1 { // +1 for the owner
		t.Fatalf("Company 2 does not have enough employees, expected %d, got %d", company_2_employee_number, len(permissions_test_instance.company2.employees))
	} else if len(permissions_test_instance.company1.branches) != company_1_branch_number {
		t.Fatalf("Company 1 does not have enough branches, expected %d, got %d", company_1_branch_number, len(permissions_test_instance.company1.branches))
	} else if len(permissions_test_instance.company2.branches) != company_2_branch_number {
		t.Fatalf("Company 2 does not have enough branches, expected %d, got %d", company_2_branch_number, len(permissions_test_instance.company2.branches))
	} else if len(permissions_test_instance.company1.services) != company_1_service_number {
		t.Fatalf("Company 1 does not have enough services, expected %d, got %d", company_1_service_number, len(permissions_test_instance.company1.services))
	} else if len(permissions_test_instance.company2.services) != company_2_service_number {
		t.Fatalf("Company 2 does not have enough services, expected %d, got %d", company_2_service_number, len(permissions_test_instance.company2.services))
	}

	permissions_test_instance.company1_employee1 = permissions_test_instance.company1.employees[1]
	permissions_test_instance.company1_employee2 = permissions_test_instance.company1.employees[2]
	permissions_test_instance.company2_employee1 = permissions_test_instance.company2.employees[1]
	permissions_test_instance.company2_employee2 = permissions_test_instance.company2.employees[2]

	if permissions_test_instance.company1_employee1.auth_token == "" {
		t.Fatal("Company 1 Employee 1 auth token is missing")
	} else if permissions_test_instance.company1_employee2.auth_token == "" {
		t.Fatal("Company 1 Employee 2 auth token is missing")
	} else if permissions_test_instance.company2_employee1.auth_token == "" {
		t.Fatal("Company 2 Employee 1 auth token is missing")
	} else if permissions_test_instance.company2_employee2.auth_token == "" {
		t.Fatal("Company 2 Employee 2 auth token is missing")
	}
}

func Test_Owner_x_Appointments(t *testing.T) {
	if permissions_test_instance == nil || permissions_test_instance.server == nil {
		Test_Start_Server(t)
		Test_Setup_Environment(t)
	}

	client1 := permissions_test_instance.client1
	client2 := permissions_test_instance.client2
	company1 := permissions_test_instance.company1
	company2 := permissions_test_instance.company2
	company1_owner := permissions_test_instance.company1.owner
	company2_owner := permissions_test_instance.company2.owner
	company1_employee1 := permissions_test_instance.company1.employees[1]
	company1_employee2 := permissions_test_instance.company1.employees[2]
	company2_employee1 := permissions_test_instance.company2.employees[1]
	company2_employee2 := permissions_test_instance.company2.employees[2]

	t.Log("--- Testing Owner x Appointment Interactions ---")

	t.Log("---> Company1 Owner creating appointment for company1_employee1 and client1 at company1 and company1.created.ID.String() : POST /appointment => HTTP 200")
	permissions_test_instance.CreateAppointment(t, 200, company1, client1, company1_employee1, company1_owner.auth_token, company1.created.ID.String())
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company1 Owner creating appointment for company1_employee1 and client2 at company1 and company1.created.ID.String() : POST /appointment => HTTP 200")
	permissions_test_instance.CreateAppointment(t, 200, company1, client2, company1_employee1, company1_owner.auth_token, company1.created.ID.String())
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company1 Owner creating appointment for company1_employee2 and client1 at company1 and company1.created.ID.String() : POST /appointment => HTTP 200")
	permissions_test_instance.CreateAppointment(t, 200, company1, client1, company1_employee2, company1_owner.auth_token, company1.created.ID.String())
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company1 Owner creating appointment for company1_employee2 and client2 at company1 and company1.created.ID.String() : POST /appointment => HTTP 200")
	permissions_test_instance.CreateAppointment(t, 200, company1, client2, company1_employee2, company1_owner.auth_token, company1.created.ID.String())
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company2 Owner creating appointment for company2_employee1 and client1 at company2 and company2.created.ID.String() : POST /appointment => HTTP 200")
	permissions_test_instance.CreateAppointment(t, 200, company2, client1, company2_employee1, company2_owner.auth_token, company2.created.ID.String())
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company2 Owner creating appointment for company2_employee1 and client2 at company2 and company2.created.ID.String() : POST /appointment => HTTP 200")
	permissions_test_instance.CreateAppointment(t, 200, company2, client2, company2_employee1, company2_owner.auth_token, company2.created.ID.String())
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company2 Owner creating appointment for company2_employee2 and client1 at company2 and company2.created.ID.String() : POST /appointment => HTTP 200")
	permissions_test_instance.CreateAppointment(t, 200, company2, client1, company2_employee2, company2_owner.auth_token, company2.created.ID.String())
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company2 Owner creating appointment for company2_employee2 and client2 at company2 and company2.created.ID.String() : POST /appointment => HTTP 200")
	permissions_test_instance.CreateAppointment(t, 200, company2, client2, company2_employee2, company2_owner.auth_token, company2.created.ID.String())
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company1 Owner trying to get an appointmentm of Employee 1 from Company 1 : GET /appointment/{id} => HTTP 200")
	permissions_test_instance.GetAppointment(t, 200, company1_employee1.created.Appointments[0].ID.String(), company1.created.ID.String(), company1_owner.auth_token)
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company1 Owner trying to get an appointment of Employee 1 from Company 2 : GET /appointment/{id} => HTTP 403")
	permissions_test_instance.GetAppointment(t, 403, company2_employee1.created.Appointments[0].ID.String(), company1.created.ID.String(), company1_owner.auth_token)
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company1 Owner trying to reschedule an appointment of Employee 1 from Company 1 : PATCH /appointment/{id} => HTTP 200")
	permissions_test_instance.RescheduleAppointment(t, 200, company1_employee1, company1, company1_employee1.created.Appointments[0].ID.String(), company1_owner.auth_token)
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company1 Owner trying to reschedule an appointment of Employee 1 from Company 2 : PATCH /appointment/{id} => HTTP 403")
	permissions_test_instance.RescheduleAppointment(t, 403, company1_employee1, company1, company2_employee1.created.Appointments[0].ID.String(), company1_owner.auth_token)
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company1 Owner trying to cancel an appointment of Employee 1 from Company 1 : DELETE /appointment/{id} => HTTP 200")
	permissions_test_instance.CancelAppointment(t, 200, company1_employee1.created.Appointments[0].ID.String(), company1.created.ID.String(), company1_owner.auth_token)
	t.Log("---------------------- x ----------------------")

	t.Log("---> Company1 Owner trying to cancel an appointment of Employee 1 from Company 2 : DELETE /appointment/{id} => HTTP 403")
	permissions_test_instance.CancelAppointment(t, 403, company2_employee1.created.Appointments[0].ID.String(), company1.created.ID.String(), company1_owner.auth_token)
	t.Log("---------------------- x ----------------------")
}

func (permissions_test) CreateAppointment(t *testing.T, s int, company *Company, client *Client, employee *Employee, token, company_id string) {
	preferredLocation := time.UTC // Choose your timezone (e.g., UTC)
	appointmentSlot, found := findValidAppointmentSlot(t, employee, company, preferredLocation)
	if !found {
		t.Logf("No valid appointment slot found for employee %s in company %s", employee.created.ID.String(), company.created.ID.String())
		t.Logf("Employee Work Schedule: %+v", employee.created.WorkSchedule)
		t.Fatal("Setup failed: Could not find a valid appointment slot for initial booking.")
	}
	http := handler.NewHttpClient(t)
	http.
		Method("POST").
		URL("/appointment").
		ExpectStatus(s).
		Header(namespace.HeadersKey.Auth, token).
		Header(namespace.HeadersKey.Company, company_id).
		Send(map[string]any{
			"branch_id":   appointmentSlot.BranchID,
			"service_id":  appointmentSlot.ServiceID,
			"employee_id": employee.created.ID.String(),
			"company_id":  company.created.ID.String(),
			"client_id":   client.created.ID.String(),
			"start_time":  appointmentSlot.StartTimeRFC3339, // Use found start time
		})
	_, ok := http.ResBody["id"].(string)
	if !ok {
		t.Fatal("Failed to get appointment id from response for client1")
	}
	company.GetById(t, 200)
	client.GetByEmail(t, 200)
	employee.GetById(t, 200)
}

func (permissions_test) RescheduleAppointment(t *testing.T, s int, employee *Employee, company *Company, appointment_id, token string) {
	preferredLocation := time.UTC // Choose your timezone (e.g., UTC)
	appointmentSlot, found := findValidAppointmentSlot(t, employee, company, preferredLocation)
	if !found {
		t.Logf("No valid appointment slot found for employee %s in company %s", employee.created.ID.String(), company.created.ID.String())
		t.Logf("Employee Work Schedule: %+v", employee.created.WorkSchedule)
		t.Fatal("Setup failed: Could not find a valid appointment slot for initial booking.")
	}
	http := handler.NewHttpClient(t)
	http.
		Method("PATCH").
		URL("/appointment/"+appointment_id).
		ExpectStatus(s).
		Header(namespace.HeadersKey.Auth, token).
		Header(namespace.HeadersKey.Company, company.created.ID.String()).
		Send(map[string]any{
			"branch_id":  appointmentSlot.BranchID,
			"service_id": appointmentSlot.ServiceID,
			"start_time": appointmentSlot.StartTimeRFC3339,
		})
	employee.GetById(t, 200)
	company.GetById(t, 200)
}

func (permissions_test) GetAppointment(t *testing.T, s int, appointment_id, company_id, token string) {
	http := handler.NewHttpClient(t)
	http.
		Method("GET").
		URL("/appointment/"+appointment_id).
		ExpectStatus(s).
		Header(namespace.HeadersKey.Auth, token).
		Header(namespace.HeadersKey.Company, company_id).
		Send(nil)
}

func (permissions_test) CancelAppointment(t *testing.T, s int, appointment_id, company_id, token string) {
	http := handler.NewHttpClient(t)
	http.
		Method("DELETE").
		URL("/appointment/"+appointment_id).
		ExpectStatus(s).
		Header(namespace.HeadersKey.Auth, token).
		Header(namespace.HeadersKey.Company, company_id).
		Send(nil)
}
