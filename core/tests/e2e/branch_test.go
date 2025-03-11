package e2e_test

import (
	"agenda-kaki-go/core"
	"agenda-kaki-go/core/config/db/model"
	handler "agenda-kaki-go/core/tests/handlers"
	"fmt"
	"testing"
)

type Branch struct {
	created model.Branch
	auth_token string
	company *Company
}

func Test_Branch(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	user := &User{}
	user.Create(t, 200)
	user.VerifyEmail(t, 200)
	user.Login(t, 200)
	company := &Company{}
	company.auth_token = user.auth_token
	company.Create(t, 200)
	branch := &Branch{}
	branch.auth_token = user.auth_token
	branch.company = company
	branch.Create(t, 200)
	branch.created.Name = "Updated Branch Name"
	branch.Update(t, 200)
	branch.GetById(t, 200)
	branch.GetByName(t, 200)
	branch.Delete(t, 200)
}

func (b *Branch) Create(t *testing.T, status int) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/branch")
	http.ExpectStatus(status)
	http.Header("Authorization", b.auth_token)
	http.Send(model.CreateBranch{
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
	b.created = model.Branch{
		Name:         http.ResBody["name"].(string),
		CompanyID:    uint(http.ResBody["company_id"].(float64)),
		Street:       http.ResBody["street"].(string),
		Number:       http.ResBody["number"].(string),
		Neighborhood: http.ResBody["neighborhood"].(string),
		ZipCode:      http.ResBody["zip_code"].(string),
		City:         http.ResBody["city"].(string),
		State:        http.ResBody["state"].(string),
		Country:      http.ResBody["country"].(string),
	}
	b.created.ID = uint(http.ResBody["id"].(float64))
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
	return http.ResBody
}

func (b *Branch) GetById(t *testing.T, status int) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(fmt.Sprintf("/branch/%d", b.created.ID))
	http.ExpectStatus(status)
	http.Send(nil)
	return http.ResBody
}

func (b *Branch) Delete(t *testing.T, status int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("DELETE")
	http.URL(fmt.Sprintf("/branch/%d", b.created.ID))
	http.ExpectStatus(status)
	http.Send(nil)
}


