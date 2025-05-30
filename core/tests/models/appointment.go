package models_test

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/tests/handlers"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

type Appointment struct {
	Created model.Appointment
}

func (a *Appointment) Create(t *testing.T, status int, auth_token string, startTime *string, b *Branch, e *Employee, s *Service, cy *Company, ct *Client) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/appointment")
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Company, cy.Created.ID.String())
	http.Header(namespace.HeadersKey.Auth, auth_token)
	if startTime == nil {
		tempStartTime := lib.GenerateDateRFC3339(2027, 10, 29)
		startTime = &tempStartTime
	}
	A := DTO.CreateAppointment{
		BranchID:   b.Created.ID,
		EmployeeID: e.Created.ID,
		ServiceID:  s.Created.ID,
		ClientID:   ct.Created.ID,
		CompanyID:  cy.Created.ID,
		StartTime:  *startTime,
	}
	http.Send(A)
	http.ParseResponse(&a.Created)
	b.GetById(t, 200)
	e.GetById(t, 200)
	s.GetById(t, 200, nil)
	cy.GetById(t, 200)
	ct.GetByEmail(t, 200)
	var ClientAppointment mJSON.ClientAppointment
	aCreatedByte, err := json.Marshal(a.Created)
	if err != nil {
		t.Fatalf("Failed to marshal appointment: %v", err)
	}
	err = json.Unmarshal(aCreatedByte, &ClientAppointment)
	if err != nil {
		t.Fatalf("Failed to unmarshal appointment: %v", err)
	}
	ct.Created.Appointments.Add(&ClientAppointment)
	e.Created.Appointments = append(e.Created.Appointments, a.Created)
	b.Created.Appointments = append(b.Created.Appointments, a.Created)
}

func (a *Appointment) GetById(t *testing.T, status int, token string) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(fmt.Sprintf("/appointment/%s", a.Created.ID.String()))
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Company, a.Created.CompanyID.String())
	http.Header(namespace.HeadersKey.Auth, token)
	http.Send(nil)
	http.ParseResponse(&a.Created)
	if a.Created.ID == uuid.Nil {
		t.Fatalf("Appointment ID is nil after GetById, something went wrong.")
	}
}
