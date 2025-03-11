package e2e_test

import (
	"agenda-kaki-go/core"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/tests/handlers"
	"fmt"
	"testing"
)

type User struct {
	created model.User
	auth_token 	string
}


func Test_User(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	user := &User{}
	user.Create(t)
	user.VerifyEmail(t)
	user.Login(t)
	user.Update(t, map[string]any{})
	user.GetByEmail(t)
	user.Delete(t)
}

func (u *User) Create(t *testing.T) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/user")
	http.ExpectStatus(200)
	http.Send(map[string]any{
		"email":    "test@email.com",
		"name":     lib.GenerateRandomName("User Name"),
		"surname":  lib.GenerateRandomName("User Surname"),
		"password": "1VerySecurePassword!",
		"phone":    lib.GenerateRandomStrNumber(11), // 55977747309
	})

	id := fmt.Sprintf("%v", http.ResBody["id"].(float64))
	email, ok := http.ResBody["email"].(string)
	if id == "" {
		t.Errorf("User ID is empty")
	} else if !ok {
		t.Errorf("Email not found in user_created")
	} else if email == "" {
		t.Errorf("Email is empty")
	}
	u.created = model.User{
		Email:   email,
		Name:    http.ResBody["name"].(string),
		Surname: http.ResBody["surname"].(string),
		Phone:   http.ResBody["phone"].(string),
	}
	u.created.ID = uint(http.ResBody["id"].(float64))
	return http.ResBody
}

func (u *User) Update(t *testing.T, body map[string]any) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("PATCH")
	http.URL("/user/" + fmt.Sprintf("%v", u.created.ID))
	http.ExpectStatus(200)
	http.Send(body)
}

func (u *User) GetByEmail(t *testing.T) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL("/user/email/" + u.created.Email)
	http.ExpectStatus(200)
	http.Send(nil)
	return http.ResBody
}

func (u *User) Delete(t *testing.T) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("DELETE")
	http.URL(fmt.Sprintf("/user/%v", u.created.ID))
	http.ExpectStatus(200)
	http.Send(nil)
}

func (u *User) VerifyEmail(t *testing.T) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(fmt.Sprintf("/auth/verify-email/%v/%s", u.created.ID, "12345"))
	http.ExpectStatus(200)
	http.Send(nil)
}

func (u *User) Login(t *testing.T) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/auth/login")
	http.ExpectStatus(200)
	http.Send(map[string]any{
		"email":    u.created.Email,
		"password": "1VerySecurePassword!",
	})
	auth := http.ResHeaders["Authorization"]
	if len(auth) == 0 {
		t.Errorf("Authorization header not found")
		return
	}
	u.auth_token = auth[0]
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
	http.URL("/login")
	http.ExpectStatus(200)
	http.Send(map[string]any{
		"email":    "test@email.com",
		"password": "1VerySecurePassword!",
	})
}
