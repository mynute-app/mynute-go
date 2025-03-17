package e2e_test

import (
	"agenda-kaki-go/core"
	DTO "agenda-kaki-go/core/config/api/dto"
	handler "agenda-kaki-go/core/tests/handlers"
	"fmt"
	"testing"
)

type Branch struct {
	created DTO.Branch
	auth_token string
	company *Company
	owner *Employee
	services *[]Service
	employees *[]Employee
}

func Test_Branch(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	company := &Company{}
	company.Create(t, 200)
	company.owner.VerifyEmail(t, 200)
	company.owner.Login(t, 200)
	company.auth_token = company.owner.auth_token
	branch := &Branch{}
	branch.auth_token = company.auth_token
	branch.company = company
	branch.Create(t, 200)
	branch.created.Name = "Updated Branch Name"
	branch.Update(t, 200)
	branch.GetById(t, 200)
	branch.GetByName(t, 200)
	service := &Service{}
	service.auth_token = company.auth_token
	service.company = company
	service.Create(t, 200)
	branch.AddService(t, 200, service)
	company.owner.AddBranch(t, 200, branch)
	company.GetById(t, 200)
	branch.Delete(t, 200)
}

func (b *Branch) Create(t *testing.T, status int) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/branch")
	http.ExpectStatus(status)
	http.Header("Authorization", b.auth_token)
	http.Send(DTO.CreateBranch{
		Name:         "Test Branch",
		CompanyID:    b.company.created.ID,
		Street:       "Test Street",
		Number:       "123",
		Neighborhood: "Test Neighborhood",
		ZipCode:      "12345678",
		City:         "Test City",
		State:        "Test State",
		Country:      "Test Country",
	})
	http.ParseResponse(&b.created)
	return http.ResBody
}

func (b *Branch) Update(t *testing.T, status int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("PATCH")
	http.URL("/branch/" + fmt.Sprintf("%v", b.created.ID))
	http.ExpectStatus(status)
	http.Header("Authorization", b.auth_token)
	http.Send(b.created)
}

func (b *Branch) GetByName(t *testing.T, status int) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(fmt.Sprintf("/branch/name/%s", b.created.Name))
	http.ExpectStatus(status)
	http.Send(nil)
	http.ParseResponse(&b.created)
	return http.ResBody
}

func (b *Branch) GetById(t *testing.T, status int) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(fmt.Sprintf("/branch/%d", b.created.ID))
	http.ExpectStatus(status)
	http.Send(nil)
	http.ParseResponse(&b.created)
	return http.ResBody
}

func (b *Branch) Delete(t *testing.T, status int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("DELETE")
	http.URL(fmt.Sprintf("/branch/%d", b.created.ID))
	http.ExpectStatus(status)
	http.Header("Authorization", b.auth_token)
	http.Send(nil)
}

func (b *Branch) AddService(t *testing.T, status int, service *Service) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL(fmt.Sprintf("/branch/%d/service/%d", b.created.ID, service.created.ID))
	http.ExpectStatus(status)
	http.Header("Authorization", b.auth_token)
	http.Send(nil)
	http.ParseResponse(&b.created)
}


