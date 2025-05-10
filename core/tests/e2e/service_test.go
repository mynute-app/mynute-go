package e2e_test

import (
	"agenda-kaki-go/core"
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/tests/handlers"
	"fmt"
	"testing"
)

type Service struct {
	created    DTO.Service
	auth_token string
	company    *Company
	employees  []*Employee
	branches   []*Branch
}

func Test_Service(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	client := &Client{}
	client.Create(t, 200)
	client.VerifyEmail(t, 200)
	client.Login(t, 200)
	company := &Company{}
	company.auth_token = client.auth_token
	company.Create(t, 200)
	service := &Service{}
	service.auth_token = client.auth_token
	service.company = company
	service.Create(t, 200)
	service.Update(t, 200, map[string]any{
		"name": lib.GenerateRandomName("Updated Service"),
	})
	service.GetById(t, 200)
	service.GetByName(t, 200)
	branch := &Branch{}
	branch.auth_token = client.auth_token
	branch.company = company
	branch.Create(t, 200)
	branch.AddService(t, 200, service, nil)
	service.Delete(t, 200)
	branch.Delete(t, 200)
}

func (s *Service) Create(t *testing.T, status int) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/service")
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Company, s.company.created.ID.String())
	http.Header(namespace.HeadersKey.Auth, s.auth_token)
	http.Send(DTO.CreateService{
		Name:        lib.GenerateRandomName("Service"),
		Description: lib.GenerateRandomName("Description"),
		CompanyID:   s.company.created.ID,
		Price:       int32(lib.GenerateRandomInt(3)),
		Duration:    60,
	})
	http.ParseResponse(&s.created)
	return http.ResBody
}

func (s *Service) Update(t *testing.T, status int, changes map[string]any) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("PATCH")
	http.Header(namespace.HeadersKey.Company, s.company.created.ID.String())
	http.URL("/service/" + fmt.Sprintf("%v", s.created.ID.String()))
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Auth, s.auth_token)
	http.Send(changes)
	s.GetById(t, 200)
}

func (s *Service) GetById(t *testing.T, status int) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.Header(namespace.HeadersKey.Company, s.company.created.ID.String())
	http.URL("/service/" + fmt.Sprintf("%v", s.created.ID.String()))
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Auth, s.auth_token)
	http.Send(nil)
	http.ParseResponse(&s.created)
	return http.ResBody
}

func (s *Service) GetByName(t *testing.T, status int) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.Header(namespace.HeadersKey.Company, s.company.created.ID.String())
	http.URL("/service/name/" + s.created.Name)
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Auth, s.auth_token)
	http.Send(nil)
	http.ParseResponse(&s.created)
	return http.ResBody
}

func (s *Service) Delete(t *testing.T, status int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("DELETE")
	http.Header(namespace.HeadersKey.Company, s.company.created.ID.String())
	http.URL("/service/" + fmt.Sprintf("%v", s.created.ID.String()))
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Auth, s.auth_token)
	http.Send(nil)
}
