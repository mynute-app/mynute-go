package e2e_test

import (
	"agenda-kaki-go/core"
	handler "agenda-kaki-go/core/tests/handlers"
	"testing"

)

func Test_Permissions(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	company1 := &Company{}
	company1.SetupRandomized(t, 5, 3, 24)
	company2 := &Company{}
	company2.SetupRandomized(t, 3, 2, 6)
	client1 := &Client{}
	client1.Set(t)
	client2 := &Client{}
	client2.Set(t)
	http := (&handler.HttpClient{}).SetTest(t)

	// --- Client x Appointment --- Interactions ---
	// Client tries to create his appointment : POST /appointment => 200
	var employee_0_start_time string
	if len(company1.employees[0].created.WorkSchedule.Monday) > 0 {
		employee_0_start_time = company1.employees[0].created.WorkSchedule.Monday[0].Start
	} else if len(company1.employees[0].created.WorkSchedule.Tuesday) > 0 {
		employee_0_start_time = company1.employees[0].created.WorkSchedule.Tuesday[0].Start
	} else if len(company1.employees[0].created.WorkSchedule.Wednesday) > 0 {
		employee_0_start_time = company1.employees[0].created.WorkSchedule.Wednesday[0].Start
	} else if len(company1.employees[0].created.WorkSchedule.Thursday) > 0 {
		employee_0_start_time = company1.employees[0].created.WorkSchedule.Thursday[0].Start
	} else if len(company1.employees[0].created.WorkSchedule.Friday) > 0 {
		employee_0_start_time = company1.employees[0].created.WorkSchedule.Friday[0].Start
	} else if len(company1.employees[0].created.WorkSchedule.Saturday) > 0 {
		employee_0_start_time = company1.employees[0].created.WorkSchedule.Saturday[0].Start
	} else if len(company1.employees[0].created.WorkSchedule.Sunday) > 0 {
		employee_0_start_time = company1.employees[0].created.WorkSchedule.Sunday[0].Start
	} else {
		t.Fatal("No work schedule found for employee 0")
	}
	http.
		Method("POST").
		URL("/appointment").
		ExpectStatus(200).
		Header("Authorization", client1.auth_token).
		Send(map[string]any{
			"branch_id":   company1.employees[0].branches[0].created.ID.String(),
			"service_id":  company1.employees[0].services[0].created.ID.String(),
			"employee_id": company1.employees[0].created.ID.String(),
			"company_id":  company1.created.ID.String(),
			"client_id":   client1.created.ID.String(),
			"start_time":  employee_0_start_time,
		})
	// Client tries to cancel his ongoing appointment : DELETE /appointment/{id} => 200
	http.
		Method("DELETE").
		URL("/appointment/"+client1.created.Appointments[0].ID.String()).
		ExpectStatus(200).
		Header("Authorization", client1.auth_token).
		Send(nil)
	// Client tries to create someone else's appointment : POST /appointment => 403
	http.
		Method("/POST").
		URL("/appointment").
		ExpectStatus(403).
		Header("Authorization", client1.auth_token).
		Send(map[string]any{
			"branch_id":   company1.employees[0].branches[0].created.ID.String(),
			"service_id":  company1.employees[0].services[0].created.ID.String(),
			"employee_id": company1.employees[0].created.ID.String(),
			"company_id":  company1.created.ID.String(),
			"client_id":   client2.created.ID.String(),
			"start_time":  employee_0_start_time,
		})
	// Client tries to get his appointment : GET /appointment/{id} => 200
	// Client tries to get someone else's appointment : GET /appointment/{id} => 403
	// Client tries to reschedule his appointment : PATCH /appointment/{id} => 200
	// Client tries to reschedule someone else's appointment : PATCH /appointment/{id} => 403

	// Client tries to cancel someone else's appointment : DELETE /appointment/{id} => 403

	// --- Client x Branch --- Interactions ---
	// Client tries to get a branch : GET /branch/{id} => 200
	// Client tries to create a branch : POST /branch => 403
	// Client tries to edit a branch : PATCH /branch/{id} => 403
	// Client tries to delete a branch : DELETE /branch/{id} => 403

	// --- Client x Client --- Interactions ---
	// Client tries to get a client : GET /client/{id} => 403
	http.
		Method("GET").
		URL("/client/"+client2.created.ID.String()).
		ExpectStatus(403).
		Header("Authorization", client1.auth_token).
		Send(nil)
	// Client tries to create a client : POST /client => 403
	// Client tries to edit a client : PATCH /client/{id} => 403
	http.
		Method("PATCH").
		URL("/client/"+client2.created.ID.String()).
		ExpectStatus(403).
		Header("Authorization", client1.auth_token).
		Send(map[string]any{
			"name": "New Client Name",
		})
	// Client tries to delete a client : DELETE /client/{id} => 403
	// Client tries to change something on himself : PATCH /client/{id} => 200
	http.
		Method("PATCH").
		URL("/client/"+client1.created.ID.String()).
		ExpectStatus(200).
		Header("Authorization", client1.auth_token).
		Send(map[string]any{
			"name": "New Client Name",
		})
	// Client tries to delete himself : DELETE /client/{id} => 200

	// --- Client x Company --- Interactions ---
	// Client tries to get a company : GET /company/{id} => 200
	// Client tries to get all companies : GET /company => 403
	// Client tries to change something in a company : PATCH /company/{id} => 403
	http.
		Method("PATCH").
		URL("/company/"+company1.created.ID.String()).
		ExpectStatus(403).
		Header("Authorization", client1.auth_token).
		Send(map[string]any{
			"name": "New Company Name",
		})

	// --- Client x Employee --- Interactions ---
	// Client tries to get an employee : GET /employee/{id} => 200
	// Client tries to create an employee : POST /employee => 403
	// Client tries to edit an employee : PATCH /employee/{id} => 403
	// Client tries to delete an employee : DELETE /employee/{id} => 403

	// --- Client x Role --- Interactions ---
	// Client tries to get a role : GET /role/{id} => 200
	// Client tries to get all roles : GET /role => 404
	// Client tries to create a role : POST /role => 403
	// Client tries to edit a role : PATCH /role/{id} => 403
	// Client tries to delete a role : DELETE /role/{id} => 403

	// --- Client x Sector --- Interactions ---
	// Client tries to get a sector : GET /sector/{id} => 200
	// Client tries to get all sectors : GET /sector => 200
	// Client tries to create a sector : POST /sector => 403
	// Client tries to edit a sector : PATCH /sector/{id} => 403
	// Client tries to delete a sector : DELETE /sector/{id} => 403

	// --- Client x Service --- Interactions ---
	// Client tries to get a service : GET /service/{id} => 200
	// Client tries to create a service : POST /service => 403
	// Client tries to edit a service : PATCH /service/{id} => 403
	// Client tries to delete a service : DELETE /service/{id} => 403

	// --- Employee x Appointments --- Interactions ---
	// Employee tries to get his appointment : GET /appointment/{id} => 200
	// Employee tries to get someone else's appointment : GET /appointment/{id} => 403
	// Employee tries to create an appointment : POST /appointment => 200
	// Employee tries to edit his appointment : PATCH /appointment/{id} => 404
	// Employee tries to delete his appointment : DELETE /appointment/{id} => 404
	// Employee tries to get someone else's appointment : GET /appointment/{id} => 403
	// Employee tries to edit someone else's appointment : PATCH /appointment/{id} => 403
	// Employee tries to delete someone else's appointment : DELETE /appointment/{id} => 403

	// --- Employee x Company --- Interactions ---
	// Employee tries to get a company : GET /company/{id} => 403
	// Employee tries to get all companies : GET /company => 403
	// Employee tries to change something in a company : PATCH /company/{id} => 403
	// Employee tries to delete a company : DELETE /company/{id} => 403

	// --- Employee x Branch --- Interactions ---
	// Employee tries to get a branch : GET /branch/{id} => 403
	// Employee tries to create a branch : POST /branch => 403
	// Employee tries to edit a branch : PATCH /branch/{id} => 403
	// Employee tries to delete a branch : DELETE /branch/{id} => 403

	// --- Employee x Service --- Interactions ---
	// Employee tries to get a service : GET /service/{id} => 200
	// Employee tries to create a service : POST /service => 403
	// Employee tries to edit a service : PATCH /service/{id} => 403
	// Employee tries to delete a service : DELETE /service/{id} => 403
	// Employee tries to add a service to himself : POST /employee/{id}/service/{id} => 200
	http.
		Method("POST").
		URL("/employee/"+company1.employees[0].created.ID.String()+"/service/"+company1.services[0].created.ID.String()).
		ExpectStatus(200).
		Header("Authorization", company1.owner.auth_token).
		Send(nil)
	// Employee tries to remove a service from himself : DELETE /employee/{id}/service/{id} => 200

}
