package e2e_test

import (
	"agenda-kaki-go/core"
	"agenda-kaki-go/core/config/db/model"
	handler "agenda-kaki-go/core/tests/handlers"
	"testing"
)

type Company struct{
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
	return http.ResBody
}
