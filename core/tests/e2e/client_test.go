package e2e_test

import (
	"agenda-kaki-go/core"
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/tests/handlers"
	"fmt"
	"testing"
)

type Client struct {
	created    model.Client
	auth_token string
}

func Test_Client(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	client := &Client{}
	client.Create(t, 200)
	client.VerifyEmail(t, 200)
	client.Login(t, 200)
	client.Update(t, 200, map[string]any{"name": "Updated Client Name"})
	client.GetByEmail(t, 200)
	client.Delete(t, 200)
}

func (u *Client) Set(t *testing.T) {
	u.Create(t, 200)
	u.VerifyEmail(t, 200)
	u.Login(t, 200)
}

func (u *Client) Create(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/client")
	http.ExpectStatus(s)
	email := lib.GenerateRandomEmail("client")
	pswd := "1SecurePswd!"
	http.Send(DTO.CreateClient{
		Email:    email,
		Name:     lib.GenerateRandomName("Client Name"),
		Surname:  lib.GenerateRandomName("Client Surname"),
		Password: pswd,
		Phone:    lib.GenerateRandomPhoneNumber(),
	})
	http.ParseResponse(&u.created)
	u.created.Password = pswd
}

func (u *Client) Update(t *testing.T, s int, changes map[string]any) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("PATCH")
	http.URL("/client/" + fmt.Sprintf("%v", u.created.ID))
	http.ExpectStatus(s)
	http.Header("Authorization", u.auth_token)
	http.Send(changes)
}

func (u *Client) GetByEmail(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL("/client/email/" + u.created.Email)
	http.ExpectStatus(s)
	http.Header("Authorization", u.auth_token)
	http.Send(nil)
}

func (u *Client) Delete(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("DELETE")
	http.URL(fmt.Sprintf("/client/%v", u.created.ID))
	http.ExpectStatus(s)
	http.Header("Authorization", u.auth_token)
	http.Send(nil)
}

func (u *Client) VerifyEmail(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL(fmt.Sprintf("/client/verify-email/%v/%s", u.created.Email, "12345"))
	http.ExpectStatus(s)
	http.Header("Authorization", u.auth_token)
	http.Send(nil)
}

func (u *Client) Login(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/client/login")
	http.ExpectStatus(s)
	http.Send(map[string]any{
		"email":    u.created.Email,
		"password": "1SecurePswd!",
	})
	auth := http.ResHeaders["Authorization"]
	if len(auth) == 0 {
		t.Errorf("Authorization header not found")
		return
	}
	u.auth_token = auth[0]
}

func Test_Client_Create_Success(t *testing.T) {
	server := core.NewServer().Run("test")
	client := &Client{}
	client.Create(t, 200)
	server.Shutdown()
}

func Test_Login_Success(t *testing.T) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/login")
	http.ExpectStatus(200)
	http.Send(map[string]any{
		"email":    "test@email.com",
		"password": "1SecurePswd!",
	})
}
