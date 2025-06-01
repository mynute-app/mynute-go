package models_test

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/test/handlers"
	"fmt"
	"testing"
)

type Branch struct {
	Created    model.Branch
	Auth_token string
	Company    *Company
	Services   []*Service
	Employees  []*Employee
}

func (b *Branch) Create(t *testing.T, status int) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/branch")
	http.Header(namespace.HeadersKey.Company, b.Company.Created.ID.String())
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Company, b.Company.Created.ID.String())
	http.Header(namespace.HeadersKey.Auth, b.Auth_token)
	http.Send(DTO.CreateBranch{
		Name:         lib.GenerateRandomName("Branch Name"),
		CompanyID:    b.Company.Created.ID,
		Street:       lib.GenerateRandomName("Street"),
		Number:       lib.GenerateRandomStrNumber(3),
		Neighborhood: lib.GenerateRandomName("Neighborhood"),
		ZipCode:      lib.GenerateRandomStrNumber(5),
		City:         lib.GenerateRandomName("City"),
		State:        lib.GenerateRandomName("State"),
		Country:      lib.GenerateRandomName("Country"),
	})
	http.ParseResponse(&b.Created)
	return http.ResBody
}

func (b *Branch) Update(t *testing.T, status int, changes map[string]any) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("PATCH")
	http.URL("/branch/" + fmt.Sprintf("%v", b.Created.ID.String()))
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Company, b.Company.Created.ID.String())
	http.Header(namespace.HeadersKey.Auth, b.Auth_token)
	http.Send(changes)
	http.ParseResponse(&b.Created)
}

func (b *Branch) GetByName(t *testing.T, status int) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(fmt.Sprintf("/branch/name/%s", b.Created.Name))
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Company, b.Company.Created.ID.String())
	http.Header(namespace.HeadersKey.Auth, b.Auth_token)
	http.Send(nil)
	http.ParseResponse(&b.Created)
	return http.ResBody
}

func (b *Branch) GetById(t *testing.T, status int) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(fmt.Sprintf("/branch/%s", b.Created.ID.String()))
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Company, b.Company.Created.ID.String())
	http.Header(namespace.HeadersKey.Auth, b.Auth_token)
	http.Send(nil)
	http.ParseResponse(&b.Created)
	return http.ResBody
}

func (b *Branch) Delete(t *testing.T, status int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("DELETE")
	http.URL(fmt.Sprintf("/branch/%s", b.Created.ID.String()))
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Company, b.Company.Created.ID.String())
	http.Header(namespace.HeadersKey.Auth, b.Auth_token)
	http.Send(nil)
}

func (b *Branch) AddService(t *testing.T, status int, service *Service, token *string) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL(fmt.Sprintf("/branch/%s/service/%s", b.Created.ID.String(), service.Created.ID.String()))
	http.ExpectStatus(status)
	if token != nil {
		http.Header(namespace.HeadersKey.Auth, *token)
	} else {
		http.Header(namespace.HeadersKey.Auth, b.Auth_token)
	}
	http.Header(namespace.HeadersKey.Company, b.Company.Created.ID.String())
	http.Send(nil)
	http.ParseResponse(&b.Created)
	service.Branches = append(service.Branches, b)
	b.Services = append(b.Services, service)
}
