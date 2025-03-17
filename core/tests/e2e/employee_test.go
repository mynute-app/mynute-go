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
}

func Test_Employee(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	company := &Company{}
	company.Set(t)
	employee := &Employee{}
	employee.company = company
	employee.Create(t, 200)
	employee.VerifyEmail(t, 200)
	employee.Login(t, 200)
	employee.GetById(t, 200)
	employee.GetByEmail(t, 200)
	employee.created.Name = "Updated Employee Name"
	employee.Update(t, 200)
	branch := &Branch{}
	branch.auth_token = company.owner.auth_token
	branch.company = company
	branch.Create(t, 200)
	employee.AddBranch(t, 200, branch)
	service := &Service{}
	service.auth_token = company.owner.auth_token
	service.company = company
	service.Create(t, 200)
	employee.AddService(t, 200, service)
	branch.AddService(t, 200, service)
	company.GetById(t, 200)
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
		Email:     lib.GenerateRandomEmail(),
		Phone:     lib.GenerateRandomPhoneNumber(),
		Password:  pswd,
	})
	http.ParseResponse(&e.created)
	e.created.Password = pswd
}

func (e *Employee) Update(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("PUT")
	http.URL(fmt.Sprintf("/employee/%d", e.created.ID))
	http.ExpectStatus(s)
	http.Header("Authorization", e.company.auth_token)
	http.Send(e.created)
	http.ParseResponse(&e.created)
}

func (e *Employee) GetById(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(fmt.Sprintf("/employee/%d", e.created.ID))
	http.ExpectStatus(s)
	http.Header("Authorization", e.company.auth_token)
	http.ParseResponse(&e.created)
}

func (e *Employee) GetByEmail(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(fmt.Sprintf("/employee/email/%s", e.created.Email))
	http.ExpectStatus(s)
	http.Header("Authorization", e.company.auth_token)
	http.ParseResponse(&e.created)
}

func (e *Employee) Delete(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("DELETE")
	http.URL(fmt.Sprintf("/employee/%d", e.created.ID))
	http.ExpectStatus(s)
	http.Header("Authorization", e.company.auth_token)
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
}
