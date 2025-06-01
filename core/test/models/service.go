package models_test

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/test/handlers"
	"fmt"
	"testing"
)

type Service struct {
	Created    DTO.Service
	Auth_token string
	Company    *Company
	Employees  []*Employee
	Branches   []*Branch
}

func (s *Service) Create(t *testing.T, status int) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/service")
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Company, s.Company.Created.ID.String())
	http.Header(namespace.HeadersKey.Auth, s.Auth_token)
	http.Send(DTO.CreateService{
		Name:        lib.GenerateRandomName("Service"),
		Description: lib.GenerateRandomName("Description"),
		CompanyID:   s.Company.Created.ID,
		Price:       int32(lib.GenerateRandomInt(3)),
		Duration:    60,
	})
	http.ParseResponse(&s.Created)
	return http.ResBody
}

func (s *Service) Update(t *testing.T, status int, changes map[string]any) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("PATCH")
	http.Header(namespace.HeadersKey.Company, s.Company.Created.ID.String())
	http.URL("/service/" + fmt.Sprintf("%v", s.Created.ID.String()))
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Auth, s.Auth_token)
	http.Send(changes)
	s.GetById(t, 200, nil)
}

func (s *Service) GetById(t *testing.T, status int, token *string) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.Header(namespace.HeadersKey.Company, s.Company.Created.ID.String())
	http.URL("/service/" + fmt.Sprintf("%v", s.Created.ID.String()))
	http.ExpectStatus(status)
	if token != nil {
		http.Header(namespace.HeadersKey.Auth, *token)
	} else {
		http.Header(namespace.HeadersKey.Auth, s.Auth_token)
	}
	http.Send(nil)
	http.ParseResponse(&s.Created)
	return http.ResBody
}

func (s *Service) GetByName(t *testing.T, status int) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.Header(namespace.HeadersKey.Company, s.Company.Created.ID.String())
	http.URL("/service/name/" + s.Created.Name)
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Auth, s.Auth_token)
	http.Send(nil)
	http.ParseResponse(&s.Created)
	return http.ResBody
}

func (s *Service) Delete(t *testing.T, status int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("DELETE")
	http.Header(namespace.HeadersKey.Company, s.Company.Created.ID.String())
	http.URL("/service/" + fmt.Sprintf("%v", s.Created.ID.String()))
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Auth, s.Auth_token)
	http.Send(nil)
}
