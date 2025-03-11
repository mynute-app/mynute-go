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
	user.Create(t, 200)
	user.VerifyEmail(t, 200)
	user.Login(t, 200)
	user.created.Name = "Updated User Name"
	user.created.Surname = "Updated User Surname"
	user.Update(t, 200)
	user.GetByEmail(t, 200)
	user.Delete(t, 200)
}

func (u *User) Create(t *testing.T, s int) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/user")
	http.ExpectStatus(s)
	http.Send(model.CreateUser{
		Email:    "test@email.com",
		Name:     lib.GenerateRandomName("User Name"),
		Surname:  lib.GenerateRandomName("User Surname"),
		Password: "1VerySecurePassword!",
		Phone:    lib.GenerateRandomStrNumber(11), // 55977747309
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

func (u *User) Update(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("PATCH")
	http.URL("/user/" + fmt.Sprintf("%v", u.created.ID))
	http.ExpectStatus(s)
	http.Send(u.created)
}

func (u *User) GetByEmail(t *testing.T, s int) map[string]any {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL("/user/email/" + u.created.Email)
	http.ExpectStatus(s)
	http.Send(nil)
	return http.ResBody
}

func (u *User) Delete(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("DELETE")
	http.URL(fmt.Sprintf("/user/%v", u.created.ID))
	http.ExpectStatus(s)
	http.Send(nil)
}

func (u *User) VerifyEmail(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(fmt.Sprintf("/auth/verify-email/%v/%s", u.created.ID, "12345"))
	http.ExpectStatus(s)
	http.Send(nil)
}

func (u *User) Login(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/auth/login")
	http.ExpectStatus(s)
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
	user.Create(t, 200)
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
