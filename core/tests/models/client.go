package models_test

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/tests/handlers"
	"fmt"
	"testing"
)

type Client struct {
	Created    model.ClientFull
	Auth_token string
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
	http.ParseResponse(&u.Created)
	u.Created.Password = pswd
}

func (u *Client) Update(t *testing.T, s int, changes map[string]any) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("PATCH")
	http.URL("/client/" + fmt.Sprintf("%v", u.Created.ID.String()))
	http.ExpectStatus(s)
	http.Header(namespace.HeadersKey.Auth, u.Auth_token)
	http.Send(changes)
}

func (u *Client) GetByEmail(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL("/client/email/" + u.Created.Email)
	http.ExpectStatus(s)
	http.Header(namespace.HeadersKey.Auth, u.Auth_token)
	http.Send(nil)
}

func (u *Client) Delete(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("DELETE")
	http.URL(fmt.Sprintf("/client/%v", u.Created.ID.String()))
	http.ExpectStatus(s)
	http.Header(namespace.HeadersKey.Auth, u.Auth_token)
	http.Send(nil)
}

func (u *Client) VerifyEmail(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL(fmt.Sprintf("/client/verify-email/%v/%s", u.Created.Email, "12345"))
	http.ExpectStatus(s)
	http.Header(namespace.HeadersKey.Auth, u.Auth_token)
	http.Send(nil)
}

func (u *Client) Login(t *testing.T, s int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/client/login")
	http.ExpectStatus(s)
	http.Send(map[string]any{
		"email":    u.Created.Email,
		"password": "1SecurePswd!",
	})
	auth := http.ResHeaders[namespace.HeadersKey.Auth]
	if len(auth) == 0 {
		t.Errorf("Authorization header not found")
		return
	}
	u.Auth_token = auth[0]
}

