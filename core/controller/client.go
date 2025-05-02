package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/middleware"
	"agenda-kaki-go/core/service"

	"github.com/gofiber/fiber/v2"
)

// EmployeeController embeds service.Base in order to extend it with the functions below
type client_controller struct {
	service.Base[model.Client, DTO.Client]
}

// CreateClient creates an client
//
//	@Summary		Create client
//	@Description	Create an client
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Param			client	body		DTO.CreateClient	true	"Client"
//	@Success		200		{object}	DTO.Client
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/client [post]
func (cc *client_controller) CreateClient(c *fiber.Ctx) error {
	var client model.Client
	if err := c.BodyParser(&client); err != nil {
		return err
	}
	if err := cc.Request.Gorm.DB.Create(&client).Error; err != nil {
		return err
	}
	res := &lib.SendResponse{Ctx: c}
	if err := res.SendDTO(200, &client, &DTO.Client{}); err != nil {
		return err
	}
	return nil
}

// LoginClient logs an client in
//
//	@Summary		Login
//	@Description	Log in an client
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Param			client	body	DTO.LoginClient	true	"Client"
//	@Success		200
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/client/login [post]
func (cc *client_controller) LoginClient(c *fiber.Ctx) error {
	var body model.Client
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	var client model.Client
	if err := cc.Request.Gorm.GetOneBy("email", body.Email, &client); err != nil {
		return err
	}
	if !client.Verified {
		return lib.Error.Client.NotVerified.SendToClient(c)
	}
	if !handler.ComparePassword(client.Password, body.Password) {
		return lib.Error.Auth.InvalidLogin.SendToClient(c)
	}
	token, err := handler.JWT(c).Encode(client)
	if err != nil {
		return err
	}
	c.Response().Header.Set("Authorization", token)
	return nil
}

// VerifyClientEmail Does the email verification for an client
//
//	@Summary		Verify email
//	@Description	Verify an client's email
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Param			email	path		string	true	"Client Email"
//	@Param			code	path		string	true	"Verification Code"
//	@Success		200		{object}	nil
//	@Failure		404		{object}	nil
//	@Router			/client/verify-email/{email}/{code} [post]
func (cc *client_controller) VerifyClientEmail(c *fiber.Ctx) error {
	res := &lib.SendResponse{Ctx: c}
	email := c.Params("email")
	var client model.Client
	client.Email = email
	if err := lib.ValidatorV10.Var(client.Email, "email"); err != nil {
		return res.Send(400, err)
	}
	if err := cc.Request.Gorm.GetOneBy("email", email, &client); err != nil {
		return err
	}
	// code := c.Params("code")
	// }
	// if client.VerificationCode != code {
	// 	return lib.Error.Auth.EmailCodeInvalid.SendToClient(c)
	// }
	client.Verified = true
	if err := cc.Request.Gorm.DB.Save(&client).Error; err != nil {
		return err
	}
	return nil
}

// GetClientByEmail retrieves an client by email
//
//	@Summary		Get client by email
//	@Description	Retrieve an client by its email
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization"
//	@Failure		401				{object}	nil
//	@Param			email			path		string	true	"Client Email"
//	@Produce		json
//	@Success		200	{object}	DTO.Client
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/client/email/{email} [get]
func (cc *client_controller) GetClientByEmail(c *fiber.Ctx) error {
	return cc.GetBy("email", c)
}

// UpdateClientById updates an client by ID
//
//	@Summary		Update client
//	@Description	Update an client
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string		true	"Client ID"
//	@Param			client	body		DTO.Client	true	"Client"
//	@Success		200		{object}	DTO.Client
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/client/{id} [patch]
func (cc *client_controller) UpdateClientById(c *fiber.Ctx) error {
	return cc.UpdateOneById(c)
}

// DeleteClientById deletes an client by ID
//
//	@Summary		Delete client
//	@Description	Delete an client
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Client ID"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		404	{object}	nil
//	@Router			/client/{id} [delete]
func (cc *client_controller) DeleteClientById(c *fiber.Ctx) error {
	return cc.DeleteOneById(c)
}

func Client(Gorm *handler.Gorm) *client_controller {
	cc := &client_controller{
		Base: service.Base[model.Client, DTO.Client]{
			Name:    namespace.ClientKey.Name,
			Request: handler.Request(Gorm),
		},
	}
	endpoint := &middleware.Endpoint{DB: Gorm}
	endpoint.BulkRegisterHandler([]fiber.Handler{
		cc.CreateClient,
		cc.LoginClient,
		cc.VerifyClientEmail,
		cc.GetClientByEmail,
		cc.UpdateClientById,
		cc.DeleteClientById,
	})
	return cc
}
