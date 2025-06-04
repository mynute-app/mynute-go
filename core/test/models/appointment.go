package modelT

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/test/handlers"
	"encoding/json"
	"fmt"
)

type Appointment struct {
	Created model.Appointment
}

func (a *Appointment) Create(status int, x_auth_token string, x_company_id *string, startTime *string, b *Branch, e *Employee, s *Service, cy *Company, ct *Client) error {
	companyIDStr := cy.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	http := handler.NewHttpClient()
	http.Method("POST")
	http.URL("/appointment")
	http.ExpectedStatus(status)
	http.Header(namespace.HeadersKey.Company, cID)
	http.Header(namespace.HeadersKey.Auth, x_auth_token)
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
	if err := b.GetById(200, x_auth_token, x_company_id); err != nil {
		return err
	}
	if err := e.GetById(200, nil, x_company_id); err != nil {
		return err
	}
	if err := s.GetById(200, x_auth_token, x_company_id); err != nil {
		return err
	}
	if err := cy.GetById(200, x_auth_token, x_company_id); err != nil {
		return err
	}
	if err := ct.GetByEmail(200); err != nil {
		return err
	}
	var ClientAppointment mJSON.ClientAppointment
	aCreatedByte, err := json.Marshal(a.Created)
	if err != nil {
		return fmt.Errorf("failed to marshal appointment: %w", err)
	}
	err = json.Unmarshal(aCreatedByte, &ClientAppointment)
	if err != nil {
		return fmt.Errorf("failed to unmarshal appointment: %w", err)
	}
	ct.Created.Appointments.Add(&ClientAppointment)
	e.Created.Appointments = append(e.Created.Appointments, a.Created)
	b.Created.Appointments = append(b.Created.Appointments, a.Created)
	return nil
}

func (a *Appointment) GetById(s int, x_auth_token string, x_company_id *string) error {
	companyIDStr := a.Created.CompanyID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("GET").
		URL("/appointment/"+a.Created.ID.String()).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Header(namespace.HeadersKey.Company, cID).
		Send(nil).
		ParseResponse(&a.Created).Error; err != nil {
		return fmt.Errorf("failed to get appointment %s: %w", a.Created.ID.String(), err)
	}
	return nil
}

func (a *Appointment) Cancel(s int, x_auth_token string, x_company_id *string) error {
	companyIDStr := a.Created.CompanyID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("DELETE").
		URL("/appointment/"+a.Created.ID.String()).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Header(namespace.HeadersKey.Company, cID).
		Send(nil).Error; err != nil {
		return fmt.Errorf("failed to cancel appointment %s: %w", a.Created.ID.String(), err)
	}
	return nil
}