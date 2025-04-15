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
	company1.SetupRandomized(t, 33, 3, 24)
	company2 := &Company{}
	company2.SetupRandomized(t, 22, 6, 6)
	client := &Client{}
	client.Set(t)
	http := (&handler.HttpClient{}).SetTest(t)
	// --- Client x Appointment --- Interactions ---
	// Client tries to get his appointment : GET /appointment/{id} => 200
	// Client tries to get someone else's appointment : GET /appointment/{id} => 403
	// Client tries to create an appointment : POST /appointment => 200
	// Client tries to edit his appointment : PATCH /appointment/{id} => 404
	// Client tries to delete his appointment : DELETE /appointment/{id} => 404
	// Client tries to get someone else's appointment : GET /appointment/{id} => 403
	// Client tries to edit someone else's appointment : PATCH /appointment/{id} => 403
	// Client tries to delete someone else's appointment : DELETE /appointment/{id} => 403
	// --- Client x Branch --- Interactions ---
	// Client tries to get a branch : GET /branch/{id} => 200
	// Client tries to create a branch : POST /branch => 403
	// Client tries to edit a branch : PATCH /branch/{id} => 403
	// Client tries to delete a branch : DELETE /branch/{id} => 403
	// --- Client x Client --- Interactions ---
	// Client tries to create a client : POST /client => 403
	// Client tries to edit a client : PATCH /client/{id} => 403
	// Client tries to delete a client : DELETE /client/{id} => 403
	// Client tries to change something on himself : PATCH /client/{id} => 200
	http.
		Method("PATCH").
		URL("/client/"+client.created.ID.String()).
		ExpectStatus(200).
		Header("Authorization", client.auth_token).
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
		Header("Authorization", client.auth_token).
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
