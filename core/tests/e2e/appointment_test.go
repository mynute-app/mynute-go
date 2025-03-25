package e2e_test

import (
	"agenda-kaki-go/core"
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/tests/handlers"
	"testing"
	"time"
)

type Appointment struct {
	created model.Appointment
}

func Test_Appointment(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	ct := &Client{}
	ct.Create(t, 200)
	ct.VerifyEmail(t, 200)
	ct.Login(t, 200)
	ct.Update(t, 200, map[string]any{"name": "Updated client Name"})
	ct.GetByEmail(t, 200)
	cy := &Company{}
	cy.Set(t)
	b := cy.branches[0]
	e := cy.employees[0]
	s := cy.services[0]
	a := []*Appointment{}
	a = append(a, &Appointment{})
	a[0].Create(t, 200, ct.auth_token, nil, b, e, s, cy, ct)
	a = append(a, &Appointment{})
	a1StartTime := lib.GenerateDateRFC3339(2027, 10, 28)
	a[1].Create(t, 200, ct.auth_token, &a1StartTime, b, e, s, cy, ct)
	a = append(a, &Appointment{})
	a2StartTime := lib.GenerateDateRFC3339(2027, 10, 27)
	a[2].Create(t, 200, cy.owner.auth_token, &a2StartTime, b, e, s, cy, ct)
	startTimeStr := ct.created.Appointments[0].StartTime.Format(time.RFC3339)
	a = append(a, &Appointment{})
	a[3].Create(t, 400, ct.auth_token, &startTimeStr, b, e, s, cy, ct)
}

func (a *Appointment) Create(t *testing.T, status int, auth_token string, startTime *string, b *Branch, e *Employee, s *Service, cy *Company, ct *Client) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/appointment")
	http.ExpectStatus(status)
	http.Header("Authorization", auth_token)
	if startTime == nil {
		tempStartTime := lib.GenerateDateRFC3339(2027, 10, 29)
		startTime = &tempStartTime
	}
	A := DTO.CreateAppointment{
		BranchID:   b.created.ID,
		EmployeeID: e.created.ID,
		ServiceID:  s.created.ID,
		ClientID:   ct.created.ID,
		CompanyID:  cy.created.ID,
		StartTime:  *startTime,
	}
	http.Send(A)
	http.ParseResponse(&a.created)
	b.GetById(t, 200)
	e.GetById(t, 200)
	s.GetById(t, 200)
	cy.GetById(t, 200)
	ct.GetByEmail(t, 200)
	ct.created.Appointments = append(ct.created.Appointments, a.created)
	e.created.Appointments = append(e.created.Appointments, a.created)
	b.created.Appointments = append(b.created.Appointments, a.created)
}
