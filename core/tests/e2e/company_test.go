package e2e_test

import (
	"agenda-kaki-go/core"
	"agenda-kaki-go/core/config/db/model"
	handler "agenda-kaki-go/core/tests/handlers"
	"fmt"
	"testing"
)

type Company struct {
	created model.Company
}

func Test_Company(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	user := &User{}
	user.Create(t)
	user.VerifyEmail(t)
	user.Login(t)
	company := &Company{}
	company.Create(t, user)
	company.Update(t, map[string]any{
		"name": "Updated Company Name",
	})
	company.Get(t)
	company.Delete(t)
}

func (c *Company) Create(t *testing.T, user *User) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/company")
	http.ExpectStatus(200)
	http.Header("Authorization", user.auth_token)
	http.Send(model.CreateCompany{
		Name:  "Test Company",
		TaxID: "41915230000168",
	})
	c.created = model.Company{
		Name:  http.ResBody["name"].(string),
		TaxID: http.ResBody["tax_id"].(string),
	}
	c.created.ID = uint(http.ResBody["id"].(float64))
	return http.ResBody
}

func (c *Company) Get(t *testing.T) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(fmt.Sprintf("/company/%d", c.created.ID))
	http.ExpectStatus(200)
	http.Send(nil)
	return http.ResBody
}

func (c *Company) Update(t *testing.T, data map[string]any) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("PATCH")
	http.URL(fmt.Sprintf("/company/%d", c.created.ID))
	http.ExpectStatus(200)
	http.Send(data)
	return http.ResBody
}

func (c *Company) Delete(t *testing.T) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("DELETE")
	http.URL(fmt.Sprintf("/company/%d", c.created.ID))
	http.ExpectStatus(200)
	http.Send(nil)
}
