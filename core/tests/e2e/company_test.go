package e2e_test

import (
	"agenda-kaki-go/core"
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	handler "agenda-kaki-go/core/tests/handlers"
	"fmt"
	"testing"
)

type Company struct {
	created    model.Company
	owner      *Employee
	employees  *[]Employee
	auth_token string
}

func Test_Company(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	company := &Company{}
	company.Create(t, 200)
	company.owner.Login(t, 200)
	company.created.Name = "Updated Company Name"
	company.auth_token = company.owner.auth_token
	company.Update(t, 200)
	company.GetById(t, 200)
	company.GetByName(t, 200)
	company.Delete(t, 200)
}

func (c *Company) Create(t *testing.T, status int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/company")
	http.ExpectStatus(status)
	http.Header("Authorization", c.auth_token)
	ownerEmail := "owner.email@gmail.com"
	ownerPswd := "1SecurePswd!"
	http.Send(DTO.CreateCompany{
		Name:          "Test Company",
		TaxID:         "41915230000168",
		OwnerName:     "Test Owner Name",
		OwnerSurname:  "Test Owner Surname",
		OwnerEmail:    ownerEmail,
		OwnerPhone:    "+15555555551",
		OwnerPassword: ownerPswd,
	})
	http.ParseResponse(&c.created)
	owner := c.created.Employees[0]
	owner.Password = ownerPswd
	owner.Email = ownerEmail
	c.owner = &Employee{
		created: owner,
	}
}

func (c *Company) GetByName(t *testing.T, status int) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(fmt.Sprintf("/company/name/%s", c.created.Name))
	http.ExpectStatus(status)
	http.Send(nil)
	return http.ResBody
}

func (c *Company) GetById(t *testing.T, status int) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(fmt.Sprintf("/company/%d", c.created.ID))
	http.ExpectStatus(status)
	http.Send(nil)
	t.Log(http.ResBody)
	return http.ResBody
}

func (c *Company) Update(t *testing.T, status int) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("PATCH")
	http.URL(fmt.Sprintf("/company/%d", c.created.ID))
	http.ExpectStatus(status)
	http.Header("Authorization", c.auth_token)
	http.Send(c.created)
	return http.ResBody
}

func (c *Company) Delete(t *testing.T, status int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("DELETE")
	http.URL(fmt.Sprintf("/company/%d", c.created.ID))
	http.ExpectStatus(status)
	http.Header("Authorization", c.auth_token)
	http.Send(nil)
}
