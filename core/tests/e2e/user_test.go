package e2e_test

import (
	"agenda-kaki-go/core"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/tests/handlers"
	"fmt"
	"testing"
)

type User struct {
}

func (u *User) Create(t *testing.T) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL(namespace.QueryKey.BaseURL + "/user")
	http.ExpectStatus(200)
	http.Send(map[string]any{
		"email":    "test@email.com",
		"name":     lib.GenerateRandomName("User Name"),
		"surname":  lib.GenerateRandomName("User Surname"),
		"password": "1VerySecurePassword!",
		"phone":    lib.GenerateRandomStrNumber(11), // 55977747309
	})
	return http.ResBody
}

func (u *User) Update(t *testing.T, body map[string]any) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("PATCH")
	user_id := body["id"].(float64)
	http.URL(namespace.QueryKey.BaseURL + "/user/" + fmt.Sprintf("%v",user_id))
	http.ExpectStatus(200)
	http.Send(body)
}

func (u *User) GetByEmail(t *testing.T, email string) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(namespace.QueryKey.BaseURL + "/user/email/" + email)
	http.ExpectStatus(200)
	http.Send(nil)
	return http.ResBody
}

func (u *User) VerifyEmail(t *testing.T, id string) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(namespace.QueryKey.BaseURL + fmt.Sprintf("/auth/verify-email/%v/%s", id, "12345"))
	http.ExpectStatus(200)
	http.Send(nil)
}

func (u *User) Delete(t *testing.T, id string) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("DELETE")
	http.URL(namespace.QueryKey.BaseURL + "/user/" + id)
	http.ExpectStatus(200)
	http.Send(nil)
}

func Test_User(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	user := &User{}
	user_created := user.Create(t)
	id := fmt.Sprintf("%v", user_created["id"].(float64))
	email, ok := user_created["email"].(string)
	if !ok {
		t.Errorf("Email not found in user_created")
	}
	if email == "" {
		t.Errorf("Email is empty")
	}
	user.VerifyEmail(t, id)
	user.Update(t, user_created)
	user.GetByEmail(t, email)
	user.Delete(t, id)
}

func Test_User_Create_Success(t *testing.T) {
	server := core.NewServer().Run("test")
	user := &User{}
	user.Create(t)
	server.Shutdown()
}

func Test_Login_Success(t *testing.T) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL(namespace.QueryKey.BaseURL + "/login")
	http.ExpectStatus(200)
	http.Send(map[string]any{
		"email":    "test@email.com",
		"password": "1VerySecurePassword!",
	})
}
