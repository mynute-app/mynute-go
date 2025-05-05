package e2e_test

import (
	"agenda-kaki-go/core"
	mJSON "agenda-kaki-go/core/config/db/model/json"
	handler "agenda-kaki-go/core/tests/handlers"
	"fmt"
	"testing"
	"time"
)

func Test_Permissions(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	company1 := &Company{}
	company1.SetupRandomized(t, 2, 3, 2) // owner, 3 employees, 2 branches, 4 services
	company2 := &Company{}
	company2.SetupRandomized(t, 3, 2, 6) // owner, 2 employees, 1 branch, 3 services
	client1 := &Client{}
	client1.Set(t)
	client2 := &Client{}
	client2.Set(t)
	http := (&handler.HttpClient{}).SetTest(t)

	// Ensure employees have auth tokens (assuming SetupRandomized/Set provides them)
	if company1.owner.auth_token == "" {
		t.Fatal("Company 1 Owner auth token is missing")
	}
	if len(company1.employees) == 0 || company1.employees[0].auth_token == "" {
		t.Fatalf("Company 1 Employee[0] auth token is missing (or no employees)")
	}
	if client1.auth_token == "" || client2.auth_token == "" {
		t.Fatal("Client auth tokens are missing")
	}
	// Helper variable for employee 0 in company 1
	employee0 := company1.employees[0]
	employee0ID := employee0.created.ID.String()
	company1ID := company1.created.ID.String()
	client1ID := client1.created.ID.String()
	client2ID := client2.created.ID.String()

	// --- Client x Appointment --- Interactions ---
	t.Log("--- Testing Client x Appointment Interactions ---")
	preferredLocation := time.UTC // Choose your timezone (e.g., UTC)
	appointmentSlot, found := findValidAppointmentSlot(t, employee0, company1, preferredLocation)
	if !found {
		t.Fatal("Setup failed: Could not find a valid appointment slot for initial booking.")
	}
	http.
		Method("POST").
		URL("/appointment").
		ExpectStatus(200).
		Header("Authorization", client1.auth_token).
		Send(map[string]any{
			"branch_id":   appointmentSlot.BranchID,
			"service_id":  appointmentSlot.ServiceID,
			"employee_id": employee0ID,
			"company_id":  company1ID,
			"client_id":   client1ID,
			"start_time":  appointmentSlot.StartTimeRFC3339, // Use found start time
		})
	appointment_id_client1, ok := http.ResBody["id"].(string)
	if !ok {
		t.Fatal("Failed to get appointment id from response for client1")
	}
	t.Logf("Client 1 created appointment %s", appointment_id_client1)

	// Client tries to get his appointment : GET /appointment/{id} => 200
	http.
		Method("GET").
		URL("/appointment/"+appointment_id_client1).
		ExpectStatus(200).
		Header("Authorization", client1.auth_token).
		Send(nil)

	// Client tries to get someone else's appointment : GET /appointment/{id} => 403
	http.
		Method("GET").
		URL("/appointment/"+appointment_id_client1).
		ExpectStatus(403).
		Header("Authorization", client2.auth_token). // Client 2 trying to access Client 1's appt
		Send(nil)

	// Client tries to reschedule someone else's appointment : PATCH /appointment/{id} => 403
	next_slot_attempt := findNextAvailableSlotRFC3339(t, employee0, appointmentSlot.StartTimeRFC3339)
	http.
		Method("PATCH").
		URL("/appointment/"+appointment_id_client1).
		ExpectStatus(403).
		Header("Authorization", client2.auth_token). // Client 2 trying to change Client 1's appt
		Send(map[string]any{
			"start_time": next_slot_attempt, // Provide a body even though it should fail
		})

	// Client tries to reschedule his appointment : PATCH /appointment/{id} => 200
	// Note: findNextAvailableSlot is simplified; rescheduling might fail if the slot isn't truly available
	http.
		Method("PATCH").
		URL("/appointment/"+appointment_id_client1).
		ExpectStatus(200).
		Header("Authorization", client1.auth_token). // Client 1 changing his own appt
		Send(map[string]any{
			"start_time": next_slot_attempt, // Use the calculated next slot
		})
	t.Logf("Client 1 rescheduled appointment %s to %s (attempted)", appointment_id_client1, next_slot_attempt)

	// Client tries to create someone else's appointment : POST /appointment => 403
	http.
		Method("POST"). // Fixed: Method was /POST
		URL("/appointment").
		ExpectStatus(403).
		Header("Authorization", client1.auth_token). // Client 1 trying to create for Client 2
		Send(map[string]any{
			"branch_id":   appointmentSlot.BranchID,
			"service_id":  appointmentSlot.ServiceID,
			"employee_id": employee0ID,
			"company_id":  company1ID,
			"client_id":   client2ID, // Different client ID
			"start_time":  appointmentSlot.StartTimeRFC3339,
		})

	// Client tries to cancel someone else's appointment : DELETE /appointment/{id} => 403
	// Client 2 trying to delete client 1's appointment
	http.
		Method("DELETE").
		URL("/appointment/"+appointment_id_client1).
		ExpectStatus(403).
		Header("Authorization", client2.auth_token).
		Send(nil)

	// Client tries to cancel his ongoing appointment : DELETE /appointment/{id} => 200
	// Client 1 cancels his own appointment
	http.
		Method("DELETE").
		URL("/appointment/"+appointment_id_client1).
		ExpectStatus(200).
		Header("Authorization", client1.auth_token).
		Send(nil)
	t.Logf("Client 1 deleted appointment %s", appointment_id_client1)

	// --- Client x Branch --- Interactions ---
	t.Log("--- Testing Client x Branch Interactions ---")
	branch0ID := company1.branches[0].created.ID.String()
	// Client tries to get a branch : GET /branch/{id} => 200
	http.
		Method("GET").
		URL("/branch/"+branch0ID).
		ExpectStatus(200).
		Header("Authorization", client1.auth_token). // Any logged-in user can view?
		Send(nil)

	// Client tries to create a branch : POST /branch => 403
	http.
		Method("POST").
		URL("/branch").
		ExpectStatus(403).
		Header("Authorization", client1.auth_token).
		Send(map[string]any{
			"company_id": company1ID,
			"name":       "Client Branch Test",
			"address":    "123 Client St",
		})

	// Client tries to edit a branch : PATCH /branch/{id} => 403
	http.
		Method("PATCH").
		URL("/branch/"+branch0ID).
		ExpectStatus(403).
		Header("Authorization", client1.auth_token).
		Send(map[string]any{
			"name": "Client Edited Branch Name",
		})

	// Client tries to delete a branch : DELETE /branch/{id} => 403
	http.
		Method("DELETE").
		URL("/branch/"+branch0ID).
		ExpectStatus(403).
		Header("Authorization", client1.auth_token).
		Send(nil)

	// --- Client x Client --- Interactions ---
	t.Log("--- Testing Client x Client Interactions ---")
	// Client tries to get a client : GET /client/{id} => 403 (Cannot get other client's details)
	http.
		Method("GET").
		URL("/client/"+client2ID).
		ExpectStatus(403).
		Header("Authorization", client1.auth_token).
		Send(nil)

	// Client tries to create a client : POST /client => 403 (Client creation likely restricted / signup flow)
	http.
		Method("POST").
		URL("/client").
		ExpectStatus(403).
		Header("Authorization", client1.auth_token).
		Send(map[string]any{
			"name":  "New Client By Client",
			"email": fmt.Sprintf("clientbyclient%d@test.com", time.Now().UnixNano()), // Unique email
			"phone": "111222333",
		})

	// Client tries to edit a client : PATCH /client/{id} => 403 (Cannot edit other client)
	http.
		Method("PATCH").
		URL("/client/"+client2ID).
		ExpectStatus(403).
		Header("Authorization", client1.auth_token).
		Send(map[string]any{
			"name": "Client Edited Other Client Name",
		})

	// Client tries to delete a client : DELETE /client/{id} => 403 (Cannot delete other client)
	http.
		Method("DELETE").
		URL("/client/"+client2ID).
		ExpectStatus(403).
		Header("Authorization", client1.auth_token).
		Send(nil)

	// Client tries to change something on himself : PATCH /client/{id} => 200
	newClient1Name := "Client 1 New Name"
	http.
		Method("PATCH").
		URL("/client/"+client1ID).
		ExpectStatus(200).
		Header("Authorization", client1.auth_token).
		Send(map[string]any{
			"name": newClient1Name,
		})
	// Optional: Verify change if GET self is allowed (assume GET /client/me or similar)
	// For now, just check status 200

	// Client tries to delete himself : DELETE /client/{id} => 200
	// Note: This will invalidate client1.auth_token for subsequent requests
	http.
		Method("DELETE").
		URL("/client/"+client1ID).
		ExpectStatus(200).
		Header("Authorization", client1.auth_token). // Use token before it's invalidated
		Send(nil)
	t.Logf("Client 1 deleted himself (%s)", client1ID)
	// client1.auth_token = "" // Mark token as invalid locally

	// --- Client x Company --- Interactions ---
	t.Log("--- Testing Client x Company Interactions ---")
	// Use client 2 now as client 1 deleted himself
	if client2.auth_token == "" {
		t.Fatal("Client 2 auth token missing for subsequent tests")
	}
	// Client tries to get a company : GET /company/{id} => 200 (Public info)
	http.
		Method("GET").
		URL("/company/"+company1ID).
		ExpectStatus(200).
		Header("Authorization", client2.auth_token).
		Send(nil)

	// Client tries to get all companies : GET /company => 403 (Listing all companies usually restricted)
	http.
		Method("GET").
		URL("/company").
		ExpectStatus(403).
		Header("Authorization", client2.auth_token).
		Send(nil)

	// Client tries to change something in a company : PATCH /company/{id} => 403
	http.
		Method("PATCH").
		URL("/company/"+company1ID).
		ExpectStatus(403).
		Header("Authorization", client2.auth_token).
		Send(map[string]any{
			"name": "Client Edited Company Name",
		})

	// Client tries to delete a company : DELETE /company/{id} => 403 (Not implemented or forbidden)
	http.
		Method("DELETE").
		URL("/company/"+company1ID).
		ExpectStatus(403). // Assuming 403 for permission denied
		Header("Authorization", client2.auth_token).
		Send(nil)

	// --- Client x Employee --- Interactions ---
	t.Log("--- Testing Client x Employee Interactions ---")
	// Client tries to get an employee : GET /employee/{id} => 200 (Needed for booking)
	http.
		Method("GET").
		URL("/employee/"+employee0ID).
		ExpectStatus(200).
		Header("Authorization", client2.auth_token).
		Send(nil)

	// Client tries to create an employee : POST /employee => 403
	http.
		Method("POST").
		URL("/employee").
		ExpectStatus(403).
		Header("Authorization", client2.auth_token).
		Send(map[string]any{
			"company_id": company1ID,
			"name":       "Client Created Employee",
			"email":      fmt.Sprintf("clientemp%d@test.com", time.Now().UnixNano()),
			"phone":      "222333444",
		})

	// Client tries to edit an employee : PATCH /employee/{id} => 403
	http.
		Method("PATCH").
		URL("/employee/"+employee0ID).
		ExpectStatus(403).
		Header("Authorization", client2.auth_token).
		Send(map[string]any{
			"name": "Client Edited Employee Name",
		})

	// Client tries to delete an employee : DELETE /employee/{id} => 403
	http.
		Method("DELETE").
		URL("/employee/"+employee0ID).
		ExpectStatus(403).
		Header("Authorization", client2.auth_token).
		Send(nil)

	// --- Client x Role --- Interactions ---
	t.Log("--- Testing Client x Role Interactions ---")
	roleID := employee0.created.Roles[0].ID.String() // Assuming employee has at least one role
	// Client tries to get a role : GET /role/{id} => 200 (Possibly public info)
	http.
		Method("GET").
		URL("/role/"+roleID).
		ExpectStatus(200).
		Header("Authorization", client2.auth_token).
		Send(nil)

	// Client tries to get all roles : GET /role => 404 (Endpoint might not exist or be public)
	http.
		Method("GET").
		URL("/role").
		ExpectStatus(404). // As per original comment
		Header("Authorization", client2.auth_token).
		Send(nil)

	// Client tries to create a role : POST /role => 403
	http.
		Method("POST").
		URL("/role").
		ExpectStatus(403).
		Header("Authorization", client2.auth_token).
		Send(map[string]any{
			"company_id": company1ID,
			"name":       "Client Role",
		})

	// Client tries to edit a role : PATCH /role/{id} => 403
	http.
		Method("PATCH").
		URL("/role/"+roleID).
		ExpectStatus(403).
		Header("Authorization", client2.auth_token).
		Send(map[string]any{
			"name": "Client Edited Role",
		})

	// Client tries to delete a role : DELETE /role/{id} => 403
	http.
		Method("DELETE").
		URL("/role/"+roleID).
		ExpectStatus(403).
		Header("Authorization", client2.auth_token).
		Send(nil)

	// --- Client x Sector --- Interactions ---
	t.Log("--- Testing Client x Sector Interactions ---")
	sectorID := company1.created.Sectors[0].ID.String()
	// Client tries to get a sector : GET /sector/{id} => 200 (Public classification)
	http.
		Method("GET").
		URL("/sector/"+sectorID).
		ExpectStatus(200).
		Header("Authorization", client2.auth_token).
		Send(nil)

	// Client tries to get all sectors : GET /sector => 200 (Public listing)
	http.
		Method("GET").
		URL("/sector").
		ExpectStatus(200).
		Header("Authorization", client2.auth_token).
		Send(nil)

	// Client tries to create a sector : POST /sector => 403
	http.
		Method("POST").
		URL("/sector").
		ExpectStatus(403).
		Header("Authorization", client2.auth_token).
		Send(map[string]any{
			"name": "Client Sector",
		})

	// Client tries to edit a sector : PATCH /sector/{id} => 403
	http.
		Method("PATCH").
		URL("/sector/"+sectorID).
		ExpectStatus(403).
		Header("Authorization", client2.auth_token).
		Send(map[string]any{
			"name": "Client Edited Sector",
		})

	// Client tries to delete a sector : DELETE /sector/{id} => 403
	http.
		Method("DELETE").
		URL("/sector/"+sectorID).
		ExpectStatus(403).
		Header("Authorization", client2.auth_token).
		Send(nil)

	// --- Client x Service --- Interactions ---
	t.Log("--- Testing Client x Service Interactions ---")
	service0ID := company1.services[0].created.ID.String()
	// Client tries to get a service : GET /service/{id} => 200 (Needed for booking)
	http.
		Method("GET").
		URL("/service/"+service0ID).
		ExpectStatus(200).
		Header("Authorization", client2.auth_token).
		Send(nil)

	// Client tries to create a service : POST /service => 403
	http.
		Method("POST").
		URL("/service").
		ExpectStatus(403).
		Header("Authorization", client2.auth_token).
		Send(map[string]any{
			"company_id": company1ID,
			"name":       "Client Service",
			"price":      10.50,
			"duration":   30,
		})

	// Client tries to edit a service : PATCH /service/{id} => 403
	http.
		Method("PATCH").
		URL("/service/"+service0ID).
		ExpectStatus(403).
		Header("Authorization", client2.auth_token).
		Send(map[string]any{
			"name": "Client Edited Service",
		})

	// Client tries to delete a service : DELETE /service/{id} => 403
	http.
		Method("DELETE").
		URL("/service/"+service0ID).
		ExpectStatus(403).
		Header("Authorization", client2.auth_token).
		Send(nil)

	// --- Setup for Employee Tests ---
	t.Log("--- Setting up for Employee Interactions ---")
	// Create an appointment for employee0 booked by client2 (owner action likely needed if client can't)
	var employee0AppointmentID string
	http.
		Method("POST").
		URL("/appointment").
		ExpectStatus(200).
		Header("Authorization", company1.owner.auth_token). // Owner creates appt for client2 with employee0
		Send(map[string]any{
			"branch_id":   appointmentSlot.BranchID,
			"service_id":  appointmentSlot.ServiceID,
			"employee_id": employee0ID,
			"company_id":  company1ID,
			"client_id":   client2ID,
			"start_time":  appointmentSlot.StartTimeRFC3339,
		})
	employee0AppointmentID, ok = http.ResBody["id"].(string)
	if !ok {
		t.Fatal("Failed to get appointment id created by owner for employee0")
	}
	t.Logf("Owner created appointment %s for employee %s with client %s", employee0AppointmentID, employee0ID, client2ID)

	// Create an appointment for employee1 (if exists) booked by client2 for permission testing
	var otherEmployeeAppointmentID string
	if len(company1.employees) > 1 {
		employee1 := company1.employees[1]
		employee1ID := employee1.created.ID.String()
		employee1Branch0 := employee1.branches[0].created.ID.String()  // Assume branch access
		employee1Service0 := employee1.services[0].created.ID.String() // Assume service access
		var employee1StartTime string
		// Find start time for employee1
		schedule1 := employee1.created.WorkSchedule
		days1 := [][]mJSON.WorkRange{schedule1.Monday, schedule1.Tuesday, schedule1.Wednesday, schedule1.Thursday, schedule1.Friday, schedule1.Saturday, schedule1.Sunday}
		foundTime1 := false
		for _, day := range days1 {
			if len(day) > 0 {
				employee1StartTime = day[0].Start
				foundTime1 = true
				break
			}
		}
		if foundTime1 {
			http.
				Method("POST").
				URL("/appointment").
				ExpectStatus(200).
				Header("Authorization", company1.owner.auth_token). // Owner creates
				Send(map[string]any{
					"branch_id":   employee1Branch0,
					"service_id":  employee1Service0,
					"employee_id": employee1ID,
					"company_id":  company1ID,
					"client_id":   client2ID, // Use client2 again
					"start_time":  employee1StartTime,
				})
			otherEmployeeAppointmentID, ok = http.ResBody["id"].(string)
			if !ok {
				t.Logf("Warning: Failed to get appointment id for employee1 %s", employee1ID)
				otherEmployeeAppointmentID = "" // Mark as unavailable
			} else {
				t.Logf("Owner created appointment %s for employee %s with client %s", otherEmployeeAppointmentID, employee1ID, client2ID)
			}
		} else {
			t.Logf("Warning: No work schedule found for employee %s, skipping other employee appt test", employee1ID)
		}
	} else {
		t.Log("Warning: Only one employee available, skipping 'other employee' appointment tests")
	}

	// --- Employee x Appointments --- Interactions ---
	t.Log("--- Testing Employee x Appointment Interactions ---")
	employee0AuthToken := employee0.auth_token
	if employee0AuthToken == "" {
		t.Fatal("Employee 0 auth token is missing for tests")
	}

	// Employee tries to get his appointment : GET /appointment/{id} => 200 (Appointment assigned to him)
	http.
		Method("GET").
		URL("/appointment/"+employee0AppointmentID).
		ExpectStatus(200).
		Header("Authorization", employee0AuthToken).
		Send(nil)

	// Employee tries to get someone else's appointment : GET /appointment/{id} => 403
	if otherEmployeeAppointmentID != "" {
		http.
			Method("GET").
			URL("/appointment/"+otherEmployeeAppointmentID).
			ExpectStatus(403).
			Header("Authorization", employee0AuthToken).
			Send(nil)
	} else {
		t.Log("Skipping test: get someone else's appointment (no other employee appointment available)")
	}

	// Employee tries to create an appointment : POST /appointment => 200 (Booking for a client)
	// Use client2 again for this test booking
	var employeeCreatedApptID string
	next_slot_for_employee_booking := findNextAvailableSlotRFC3339(t, employee0, appointmentSlot.StartTimeRFC3339)
	http.
		Method("POST").
		URL("/appointment").
		ExpectStatus(200).
		Header("Authorization", employee0AuthToken). // Employee making the booking
		Send(map[string]any{
			"branch_id":   appointmentSlot.BranchID,
			"service_id":  appointmentSlot.ServiceID,
			"employee_id": employee0ID, // Booking is for himself
			"company_id":  company1ID,
			"client_id":   client2ID, // Booking *for* client 2
			"start_time":  next_slot_for_employee_booking,
		})
	employeeCreatedApptID, ok = http.ResBody["id"].(string)
	if !ok {
		t.Log("Warning: Failed to get ID for appointment created by employee")
	} else {
		t.Logf("Employee %s created appointment %s for client %s", employee0ID, employeeCreatedApptID, client2ID)
	}

	// Employee tries to edit his appointment : PATCH /appointment/{id} => 404 (Endpoint/action not allowed/found for employee?)
	// Using the appointment originally created by the owner for this employee
	http.
		Method("PATCH").
		URL("/appointment/"+employee0AppointmentID).
		ExpectStatus(404). // Following comment expectation
		Header("Authorization", employee0AuthToken).
		Send(map[string]any{
			"start_time": next_slot_for_employee_booking, // Attempt change
		})

	// Employee tries to delete his appointment : DELETE /appointment/{id} => 404 (Endpoint/action not allowed/found for employee?)
	// Using the appointment originally created by the owner for this employee
	http.
		Method("DELETE").
		URL("/appointment/"+employee0AppointmentID).
		ExpectStatus(404). // Following comment expectation
		Header("Authorization", employee0AuthToken).
		Send(nil)

	// Re-test: Employee tries to get someone else's appointment : GET /appointment/{id} => 403 (Duplicate check from above)
	if otherEmployeeAppointmentID != "" {
		http.
			Method("GET").
			URL("/appointment/"+otherEmployeeAppointmentID).
			ExpectStatus(403).
			Header("Authorization", employee0AuthToken).
			Send(nil)
	} // Skipped if no other appt

	// Employee tries to edit someone else's appointment : PATCH /appointment/{id} => 403 (Permission denied)
	if otherEmployeeAppointmentID != "" {
		http.
			Method("PATCH").
			URL("/appointment/"+otherEmployeeAppointmentID).
			ExpectStatus(403). // Should be forbidden (or 404 if PATCH route is generally unavailable to employees)
			Header("Authorization", employee0AuthToken).
			Send(map[string]any{
				"start_time": next_slot_for_employee_booking,
			})
	} else {
		t.Log("Skipping test: edit someone else's appointment (no other employee appointment available)")
	}

	// Employee tries to delete someone else's appointment : DELETE /appointment/{id} => 403 (Permission denied)
	if otherEmployeeAppointmentID != "" {
		http.
			Method("DELETE").
			URL("/appointment/"+otherEmployeeAppointmentID).
			ExpectStatus(403). // Should be forbidden (or 404 if DELETE route unavailable)
			Header("Authorization", employee0AuthToken).
			Send(nil)
	} else {
		t.Log("Skipping test: delete someone else's appointment (no other employee appointment available)")
	}
	// Clean up employee-created appointment if ID was captured
	if employeeCreatedApptID != "" {
		// Cancellation might need owner/client permission or specific endpoint, using owner for cleanup
		http.
			Method("DELETE").
			URL("/appointment/"+employeeCreatedApptID).
			ExpectStatus(200). // Assume owner can delete any appt in their company
			Header("Authorization", company1.owner.auth_token).
			Send(nil)
		t.Logf("Cleaned up employee-created appointment %s using owner token", employeeCreatedApptID)
	}

	// --- Employee x Company --- Interactions ---
	t.Log("--- Testing Employee x Company Interactions ---")
	// Employee tries to get a company : GET /company/{id} => 403 (Assume restricted internal info)
	http.
		Method("GET").
		URL("/company/"+company1ID).
		ExpectStatus(403).
		Header("Authorization", employee0AuthToken).
		Send(nil)

	// Employee tries to get all companies : GET /company => 403 (Definitely restricted)
	http.
		Method("GET").
		URL("/company").
		ExpectStatus(403).
		Header("Authorization", employee0AuthToken).
		Send(nil)

	// Employee tries to change something in a company : PATCH /company/{id} => 403
	http.
		Method("PATCH").
		URL("/company/"+company1ID).
		ExpectStatus(403).
		Header("Authorization", employee0AuthToken).
		Send(map[string]any{
			"name": "Employee Edited Company Name",
		})

	// Employee tries to delete a company : DELETE /company/{id} => 403
	http.
		Method("DELETE").
		URL("/company/"+company1ID).
		ExpectStatus(403).
		Header("Authorization", employee0AuthToken).
		Send(nil)

	// --- Employee x Branch --- Interactions ---
	t.Log("--- Testing Employee x Branch Interactions ---")
	// Employee tries to get a branch : GET /branch/{id} => 403 (Assume restricted)
	http.
		Method("GET").
		URL("/branch/"+branch0ID).
		ExpectStatus(403).
		Header("Authorization", employee0AuthToken).
		Send(nil)

	// Employee tries to create a branch : POST /branch => 403
	http.
		Method("POST").
		URL("/branch").
		ExpectStatus(403).
		Header("Authorization", employee0AuthToken).
		Send(map[string]any{
			"company_id": company1ID,
			"name":       "Employee Branch",
			"address":    "456 Employee Ave",
		})

	// Employee tries to edit a branch : PATCH /branch/{id} => 403
	http.
		Method("PATCH").
		URL("/branch/"+branch0ID).
		ExpectStatus(403).
		Header("Authorization", employee0AuthToken).
		Send(map[string]any{
			"name": "Employee Edited Branch Name",
		})

	// Employee tries to delete a branch : DELETE /branch/{id} => 403
	http.
		Method("DELETE").
		URL("/branch/"+branch0ID).
		ExpectStatus(403).
		Header("Authorization", employee0AuthToken).
		Send(nil)

	// --- Employee x Service --- Interactions ---
	t.Log("--- Testing Employee x Service Interactions ---")
	// Employee tries to get a service : GET /service/{id} => 200 (Needed for work)
	http.
		Method("GET").
		URL("/service/"+service0ID).
		ExpectStatus(200).
		Header("Authorization", employee0AuthToken).
		Send(nil)

	// Employee tries to create a service : POST /service => 403
	http.
		Method("POST").
		URL("/service").
		ExpectStatus(403).
		Header("Authorization", employee0AuthToken).
		Send(map[string]any{
			"company_id": company1ID,
			"name":       "Employee Service",
			"price":      20.00,
			"duration":   60,
		})

	// Employee tries to edit a service : PATCH /service/{id} => 403
	http.
		Method("PATCH").
		URL("/service/"+service0ID).
		ExpectStatus(403).
		Header("Authorization", employee0AuthToken).
		Send(map[string]any{
			"name": "Employee Edited Service",
		})

	// Employee tries to delete a service : DELETE /service/{id} => 403
	http.
		Method("DELETE").
		URL("/service/"+service0ID).
		ExpectStatus(403).
		Header("Authorization", employee0AuthToken).
		Send(nil)

	// Employee tries to add a service to himself : POST /employee/{id}/service/{id} => 403 (Assume only Owner/Admin can)
	// Let's use service[1] from company1 if available
	service1ID := ""
	if len(company1.services) > 1 {
		service1ID = company1.services[1].created.ID.String()
		http.
			Method("POST").
			URL("/employee/"+employee0ID+"/service/"+service1ID).
			ExpectStatus(403).
			Header("Authorization", employee0AuthToken). // Employee trying to add to himself
			Send(nil)
	} else {
		t.Log("Skipping test: employee add service to self (only one service defined)")
	}

	// Owner adds service[0] to employee[0] (this test already existed, added log)
	http.
		Method("POST").
		URL("/employee/"+employee0ID+"/service/"+service0ID).
		ExpectStatus(200).
		Header("Authorization", company1.owner.auth_token). // Owner performs action
		Send(nil)
	t.Logf("Owner added service %s to employee %s", service0ID, employee0ID)

	// Employee tries to remove a service from himself : DELETE /employee/{id}/service/{id} => 403 (Assume only Owner/Admin can)
	http.
		Method("DELETE").
		URL("/employee/"+employee0ID+"/service/"+service0ID). // Try to remove the service added by owner
		ExpectStatus(403).
		Header("Authorization", employee0AuthToken). // Employee trying action
		Send(nil)

	// Owner removes service from employee: DELETE /employee/{id}/service/{id} => 200 (Cleanup / successful case)
	http.
		Method("DELETE").
		URL("/employee/"+employee0ID+"/service/"+service0ID).
		ExpectStatus(200).
		Header("Authorization", company1.owner.auth_token). // Owner performs action
		Send(nil)
	t.Logf("Owner removed service %s from employee %s", service0ID, employee0ID)

}
