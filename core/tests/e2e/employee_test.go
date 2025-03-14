package e2e_test

// import (
// 	"agenda-kaki-go/core"
// 	"agenda-kaki-go/core/config/db/model"
// 	handler "agenda-kaki-go/core/tests/handlers"
// 	"fmt"
// 	"testing"

// 	"github.com/prometheus/common/server"
// )

// type Employee struct {
// 	user       *User
// 	company    *Company
// 	created    model.Employee
// }

// func (e *Employee) Test_User(t *testing.T) *Employee {
// 	server := core.NewServer().Run("test")
// 	defer server.Shutdown()
// 	user := &User{}
// 	user.Create(t, 200)
// 	user.VerifyEmail(t, 200)
// 	user.Login(t, 200)
// 	company := &Company{}
// 	company.auth_token = user.auth_token
// 	company.Create(t, 200)
// 	e.user = user
// 	e.company = company
// 	e.Create(t, 200)
// }

// func (e *Employee) Create(t *testing.T, status int) map[string]any {
// 	http := (&handler.HttpClient{}).SetTest(t)
// 	http.Method("POST")
// 	http.URL("/employee")
// 	http.ExpectStatus(status)
// 	http.Header("Authorization", e.user.auth_token)
// 	http.Send(model.CreateEmployee{
// 		CompanyID: e.company.created.ID,
// 		Name:      "Test Employee Name",
// 		Surname:   "Test Surname",
// 		Email:     "test_new_employee@gmail.com",
// 		Phone:     "+15555555552",
// 	})
// 	http.ParseResponse(&e.created)
// 	return http.ResBody
// }

// func (e *Employee) Update(t *testing.T, status int) map[string]any {
// 	http := (&handler.HttpClient{}).SetTest(t)
// 	http.Method("PUT")
// 	http.URL(fmt.Sprintf("/employee/%d", e.created.ID))
// 	http.ExpectStatus(status)
// 	http.Header("Authorization", e.user.auth_token)
// 	http.Send(e.created)
// 	http.ParseResponse(&e.created)
// 	return http.ResBody
// }

// func (e *Employee) Get(t *testing.T, status int) map[string]any {
// 	http := (&handler.HttpClient{}).SetTest(t)
// 	http.Method("GET")
// 	http.URL("/employee/" + string(e.created.ID))
// 	http.ExpectStatus(status)
// 	http.Header("Authorization", e.user.auth_token)
// 	http.ParseResponse(&e.created)
// 	return http.ResBody
// }

// func (e *Employee) Delete(t *testing.T, status int) map[string]any {
// 	http := (&handler.HttpClient{}).SetTest(t)
// 	http.Method("DELETE")
// 	http.URL("/employee/" + string(e.created.ID))
// 	http.ExpectStatus(status)
// 	http.Header("Authorization", e.user.auth_token)
// 	http.ParseResponse(&e.created)
// 	return http.ResBody
// }
