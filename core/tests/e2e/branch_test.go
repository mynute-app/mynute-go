package e2e_test

import (
	"agenda-kaki-go/core"
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/tests/handlers"
	"fmt"
	"testing"
)

type Branch struct {
	created    DTO.Branch
	auth_token string
	company    *Company
	services   []*Service
	employees  []*Employee
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
	branch.Update(t, 200, map[string]any{
		"name": branch.created.Name,
	})
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
		Name:         lib.GenerateRandomName("Branch Name"),
		CompanyID:    b.company.created.ID,
		Street:       lib.GenerateRandomName("Street"),
		Number:       lib.GenerateRandomStrNumber(3),
		Neighborhood: lib.GenerateRandomName("Neighborhood"),
		ZipCode:      lib.GenerateRandomStrNumber(5),
		City:         lib.GenerateRandomName("City"),
		State:        lib.GenerateRandomName("State"),
		Country:      lib.GenerateRandomName("Country"),
	})
	http.ParseResponse(&b.created)
	return http.ResBody
}

func (b *Branch) Update(t *testing.T, status int, changes map[string]any) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("PATCH")
	http.URL("/branch/" + fmt.Sprintf("%v", b.created.ID))
	http.ExpectStatus(status)
	http.Header("Authorization", b.auth_token)
	http.Send(changes)
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
	service.branches = append(service.branches, b)
	b.services = append(b.services, service)
}
