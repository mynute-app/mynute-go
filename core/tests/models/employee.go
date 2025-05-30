package models_test

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/tests/handlers"
	"fmt"
	"testing"
)

type Employee struct {
	Auth_token string
	Company    *Company
	Created    model.Employee
	Services   []*Service
	Branches   []*Branch
}


func (e *Employee) Create(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/employee")
	http.ExpectStatus(s)
	http.Header(namespace.HeadersKey.Company, e.Company.Created.ID.String())
	http.Header(namespace.HeadersKey.Auth, e.Company.Auth_token)
	pswd := "1SecurePswd!"
	CreateEmployeeBody := &DTO.CreateEmployee{
		CompanyID: e.Company.Created.ID,
		Name:      lib.GenerateRandomName("Employee Name"),
		Surname:   lib.GenerateRandomName("Employee Surname"),
		Email:     lib.GenerateRandomEmail("employee"),
		Phone:     lib.GenerateRandomPhoneNumber(),
		Password:  pswd,
	}
	http.Send(CreateEmployeeBody)
	http.ParseResponse(&e.Created)
	e.Created.Password = pswd
}

func (e *Employee) CreateBranch(t *testing.T, s int) {
	Branch := &Branch{}
	Branch.Auth_token = e.Auth_token
	Branch.Company = e.Company
	Branch.Create(t, s)
	e.Company.Branches = append(e.Company.Branches, Branch)
}

func (e *Employee) CreateService(t *testing.T, s int) {
	Service := &Service{}
	Service.Auth_token = e.Auth_token
	Service.Company = e.Company
	Service.Create(t, s)
	e.Company.Services = append(e.Company.Services, Service)
}

func (e *Employee) Update(t *testing.T, s int, changes map[string]any) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("PATCH")
	http.Header(namespace.HeadersKey.Company, e.Company.Created.ID.String())
	http.URL(fmt.Sprintf("/employee/%s", e.Created.ID.String()))
	http.ExpectStatus(s)
	http.Header(namespace.HeadersKey.Auth, e.Company.Auth_token)
	http.Send(changes)
	http.ParseResponse(&e.Created)
}

func (e *Employee) GetById(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.Header(namespace.HeadersKey.Company, e.Company.Created.ID.String())
	http.URL(fmt.Sprintf("/employee/%s", e.Created.ID.String()))
	http.ExpectStatus(s)
	http.Header(namespace.HeadersKey.Auth, e.Company.Auth_token)
	http.Send(nil)
	http.ParseResponse(&e.Created)
}

func (e *Employee) GetByEmail(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.Header(namespace.HeadersKey.Company, e.Company.Created.ID.String())
	http.URL(fmt.Sprintf("/employee/email/%s", e.Created.Email))
	http.ExpectStatus(s)
	http.Header(namespace.HeadersKey.Auth, e.Company.Auth_token)
	http.Send(nil)
	http.ParseResponse(&e.Created)
}

func (e *Employee) Delete(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("DELETE")
	http.Header(namespace.HeadersKey.Company, e.Company.Created.ID.String())
	http.URL(fmt.Sprintf("/employee/%s", e.Created.ID.String()))
	http.ExpectStatus(s)
	http.Header(namespace.HeadersKey.Company, e.Company.Created.ID.String())
	http.Header(namespace.HeadersKey.Auth, e.Company.Auth_token)
	http.Send(nil)
}

func (e *Employee) Login(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/employee/login")
	http.ExpectStatus(s)
	http.Header(namespace.HeadersKey.Company, e.Created.CompanyID.String())
	http.Send(model.Employee{
		Email:    e.Created.Email,
		Password: e.Created.Password,
	})
	auth := http.ResHeaders[namespace.HeadersKey.Auth]
	if len(auth) == 0 {
		t.Errorf("Authorization header not found")
	}
	e.Auth_token = auth[0]
}

func (e *Employee) VerifyEmail(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.Header(namespace.HeadersKey.Company, e.Created.CompanyID.String())
	http.URL(fmt.Sprintf("/employee/verify-email/%v/%s", e.Created.Email, "12345"))
	http.ExpectStatus(s)
	http.Send(nil)
}

func (e *Employee) AddBranch(t *testing.T, s int, branch *Branch, token *string) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL(fmt.Sprintf("/employee/%s/branch/%s", e.Created.ID.String(), branch.Created.ID.String()))
	http.ExpectStatus(s)
	if token != nil {
		http.Header(namespace.HeadersKey.Auth, *token)
	} else {
		http.Header(namespace.HeadersKey.Auth, e.Auth_token)
	}
	http.Header(namespace.HeadersKey.Company, e.Company.Created.ID.String())
	http.Send(nil)
	http.ParseResponse(&e.Created)
	branch.GetById(t, 200)
	branch.Employees = append(branch.Employees, e)
	e.Branches = append(e.Branches, branch)
}

func (e *Employee) AddService(t *testing.T, s int, service *Service) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL(fmt.Sprintf("/employee/%s/service/%s", e.Created.ID.String(), service.Created.ID.String()))
	http.ExpectStatus(s)
	http.Header(namespace.HeadersKey.Company, e.Company.Created.ID.String())
	http.Header(namespace.HeadersKey.Auth, e.Auth_token)
	http.Send(nil)
	http.ParseResponse(&e.Created)
	service.GetById(t, 200, nil)
	service.Employees = append(service.Employees, e)
	e.Services = append(e.Services, service)
}

func (e *Employee) AddRole(t *testing.T, s int, role *Role) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL(fmt.Sprintf("/employee/%s/role/%s", e.Created.ID.String(), role.Created.ID.String()))
	http.ExpectStatus(s)
	http.Header(namespace.HeadersKey.Auth, e.Auth_token)
	http.Send(nil)
	http.ParseResponse(&e.Created)
	role.Employees = append(role.Employees, e)
}