package e2e_test

import (
	"agenda-kaki-go/core"
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/tests/handlers"
	"fmt"
	"testing"
)

// import (
// 	"agenda-kaki-go/core"
// 	"agenda-kaki-go/core/config/db/model"
// 	handler "agenda-kaki-go/core/tests/handlers"
// 	"fmt"
// 	"testing"

// 	"github.com/prometheus/common/server"
// )

type Employee struct {
	auth_token string
	company    *Company
	created    model.Employee
	services   []*Service
	branches   []*Branch
}

func Test_Employee(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	company := &Company{}
	company.Set(t)
	employee := company.employees[0]
	employee.GetById(t, 200)
	employee.GetByEmail(t, 200)
	employee.Update(t, 200, map[string]any{"name": "Updated Employee Name xD"})
	employee.Delete(t, 200)
}

func (e *Employee) Create(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/employee")
	http.ExpectStatus(s)
	http.Header("Authorization", e.company.auth_token)
	pswd := "1SecurePswd!"
	http.Send(DTO.CreateEmployee{
		CompanyID: e.company.created.ID,
		Name:      lib.GenerateRandomName("Employee Name"),
		Surname:   lib.GenerateRandomName("Employee Surname"),
		Email:     lib.GenerateRandomEmail("employee"),
		Phone:     lib.GenerateRandomPhoneNumber(),
		Password:  pswd,
	})
	http.ParseResponse(&e.created)
	e.created.Password = pswd
}

func (e *Employee) CreateBranch(t *testing.T, s int) {
	Branch := &Branch{}
	Branch.auth_token = e.auth_token
	Branch.company = e.company
	Branch.Create(t, s)
	e.company.branches = append(e.company.branches, Branch)
}

func (e *Employee) CreateService(t *testing.T, s int) {
	Service := &Service{}
	Service.auth_token = e.auth_token
	Service.company = e.company
	Service.Create(t, s)
	e.company.services = append(e.company.services, Service)
}

func (e *Employee) Update(t *testing.T, s int, changes map[string]any) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("PATCH")
	http.URL(fmt.Sprintf("/employee/%d", e.created.ID))
	http.ExpectStatus(s)
	http.Header("Authorization", e.company.auth_token)
	http.Send(changes)
	http.ParseResponse(&e.created)
}

func (e *Employee) GetById(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(fmt.Sprintf("/employee/%d", e.created.ID))
	http.ExpectStatus(s)
	http.Header("Authorization", e.company.auth_token)
	http.Send(nil)
	http.ParseResponse(&e.created)
}

func (e *Employee) GetByEmail(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(fmt.Sprintf("/employee/email/%s", e.created.Email))
	http.ExpectStatus(s)
	http.Header("Authorization", e.company.auth_token)
	http.Send(nil)
	http.ParseResponse(&e.created)
}

func (e *Employee) Delete(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("DELETE")
	http.URL(fmt.Sprintf("/employee/%d", e.created.ID))
	http.ExpectStatus(s)
	http.Header("Authorization", e.company.auth_token)
	http.Send(nil)
	http.ParseResponse(&e.created)
}

func (e *Employee) Login(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/employee/login")
	http.ExpectStatus(s)
	http.Send(model.Employee{
		Email:    e.created.Email,
		Password: e.created.Password,
	})
	auth := http.ResHeaders["Authorization"]
	if len(auth) == 0 {
		t.Errorf("Authorization header not found")
	}
	e.auth_token = auth[0]
}

func (e *Employee) VerifyEmail(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL(fmt.Sprintf("/employee/verify-email/%v/%s", e.created.Email, "12345"))
	http.ExpectStatus(s)
	http.Send(nil)
}

func (e *Employee) AddBranch(t *testing.T, s int, branch *Branch) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL(fmt.Sprintf("/employee/%d/branch/%d", e.created.ID, branch.created.ID))
	http.ExpectStatus(s)
	http.Header("Authorization", e.auth_token)
	http.Send(nil)
	http.ParseResponse(&e.created)
	branch.GetById(t, 200)
	branch.employees = append(branch.employees, e)
	e.branches = append(e.branches, branch)
}

func (e *Employee) AddService(t *testing.T, s int, service *Service) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL(fmt.Sprintf("/employee/%d/service/%d", e.created.ID, service.created.ID))
	http.ExpectStatus(s)
	http.Header("Authorization", e.auth_token)
	http.Send(nil)
	http.ParseResponse(&e.created)
	service.GetById(t, 200)
	service.employees = append(service.employees, e)
	e.services = append(e.services, service)
}
