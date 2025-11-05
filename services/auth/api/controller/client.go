package controller

import (
	authModel "mynute-go/services/auth/config/db/model"
	DTO "mynute-go/services/auth/config/dto"
	"mynute-go/services/auth/lib"

	"github.com/gofiber/fiber/v2"
)

// =====================
// CLIENT MANAGEMENT
// =====================

// CreateClient creates a client
//
//	@Summary		Create client
//	@Description	Create a client
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Param			client	body		DTO.CreateClient	true	"Client"
//	@Success		200		{object}	DTO.Client
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/users/client [post]
func CreateClient(c *fiber.Ctx) error {
	var user authModel.User
	user.Type = "client"
	if err := CreateUser(c, &user); err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &user, &DTO.Client{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// GetClientByEmail retrieves a client by email
//
//	@Summary		Get client by email
//	@Description	Retrieve a client by its email
//	@Tags			Client
//	@Param			email	path		string	true	"Client Email"
//	@Produce		json
//	@Success		200	{object}	DTO.Client
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/users/client/email/{email} [get]
func GetClientByEmail(c *fiber.Ctx) error {
	var user authModel.User
	if err := GetOneBy("email", c, &user); err != nil {
		return err
	}
	if user.Type != "client" {
		return lib.Error.General.RecordNotFound
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &user, &DTO.Client{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// GetClientById retrieves a client by ID
//
//	@Summary		Get client by ID
//	@Description	Retrieve a client by its ID
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Client ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Client
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/users/client/{id} [get]
func GetClientById(c *fiber.Ctx) error {
	var user authModel.User
	if err := GetOneBy("id", c, &user); err != nil {
		return err
	}
	if user.Type != "client" {
		return lib.Error.General.RecordNotFound
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &user, &DTO.Client{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// UpdateClientById updates a client by ID
//
//	@Summary		Update client
//	@Description	Update a client
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Client ID"
//	@Param			client	body		DTO.UpdateClientRequest		true	"Client"
//	@Success		200		{object}	DTO.Client
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/users/client/{id} [patch]
func UpdateClientById(c *fiber.Ctx) error {
	var user authModel.User
	if err := UpdateOneById(c, &user); err != nil {
		return err
	}
	if user.Type != "client" {
		return lib.Error.General.RecordNotFound
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &user, &DTO.Client{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// DeleteClientById deletes a client by ID
//
//	@Summary		Delete client by ID
//	@Description	Delete a client by its ID
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Client ID"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/users/client/{id} [delete]
func DeleteClientById(c *fiber.Ctx) error {
	return DeleteOneById(c, &authModel.User{})
}
