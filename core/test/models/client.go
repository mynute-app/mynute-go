package modelT

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/test/handlers"
	"fmt"
)

type Client struct {
	Created    model.ClientFull
	Auth_token string
}

func (u *Client) Set() error {
	if err := u.Create(200); err != nil {
		return err
	}
	if err := u.VerifyEmail(200); err != nil {
		return err
	}
	if err := u.Login(200); err != nil {
		return err
	}
	return nil
}

func (u *Client) Create(s int) error {
	pswd := "1SecurePswd!"
	if err := handler.NewHttpClient().
		Method("POST").
		URL("/client").
		ExpectedStatus(s).
		Send(DTO.CreateClient{
			Email:    lib.GenerateRandomEmail("client"),
			Name:     lib.GenerateRandomName("Client Name"),
			Surname:  lib.GenerateRandomName("Client Surname"),
			Password: pswd,
			Phone:    lib.GenerateRandomPhoneNumber(),
		}).ParseResponse(&u.Created).Error; err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	u.Created.Password = pswd
	return nil
}

func (u *Client) Update(s int, changes map[string]any) error {
	if err := handler.NewHttpClient().
		Method("PATCH").
		URL("/client/" + fmt.Sprintf("%v", u.Created.ID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, u.Auth_token).
		Send(changes).
		Error; err != nil {
		return fmt.Errorf("failed to update client: %w", err)
	}
	return nil
}

func (u *Client) GetByEmail(s int) error {
	if err := handler.NewHttpClient().
		Method("GET").
		URL("/client/email/" + u.Created.Email).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, u.Auth_token).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to get client by email: %w", err)
	}
	return nil
}

func (u *Client) Delete(s int) error {
	if err := handler.NewHttpClient().
		Method("DELETE").
		URL(fmt.Sprintf("/client/%v", u.Created.ID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, u.Auth_token).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}
	return nil
}

func (u *Client) VerifyEmail(s int) error {
	if err := handler.NewHttpClient().
		Method("POST").
		URL(fmt.Sprintf("/client/verify-email/%v/%s", u.Created.Email, "12345")).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, u.Auth_token).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to verify client email: %w", err)
	}
	return nil
}

func (u *Client) Login(s int) error {
	login := DTO.LoginClient{
		Email:    u.Created.Email,
		Password: "1SecurePswd!",
	}
	http := handler.NewHttpClient()
	if err := http.
		Method("POST").
		URL("/client/login").
		ExpectedStatus(s).
		Send(login).Error; err != nil {
		return fmt.Errorf("failed to login client: %w", err)
	}
	auth := http.ResHeaders[namespace.HeadersKey.Auth]
	if len(auth) == 0 {
		return fmt.Errorf("Authorization header not found")
	}
	u.Auth_token = auth[0]
	return nil
}
