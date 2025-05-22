package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	database "agenda-kaki-go/core/config/db"
	"agenda-kaki-go/core/config/db/model"
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/middleware"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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
func CreateClient(c *fiber.Ctx) error {
	var client model.ClientFull
	if err := c.BodyParser(&client); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	tx, end, err := database.ContextTransaction(c)
	defer end()
	if err != nil {
		return err
	}
	client.Appointments = mJSON.ClientAppointments{}
	if err := tx.Model(&model.ClientFull{}).Create(&client).Error; err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &client, &DTO.Client{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
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
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/login [post]
func LoginClient(c *fiber.Ctx) error {
	var body DTO.LoginClient
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	tx, end, err := database.ContextTransaction(c)
	defer end()
	if err != nil {
		return err
	}

	var client model.ClientMeta
	if err := tx.Model(&model.ClientFull{}).Where("email = ?", body.Email).First(&client).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.Client.NotFound
		}
		return lib.Error.General.InternalError.WithError(err)
	}
	if !client.Verified {
		return lib.Error.Client.NotVerified
	}
	if !handler.ComparePassword(client.Password, body.Password) {
		return lib.Error.Auth.InvalidLogin
	}
	token, err := handler.JWT(c).Encode(client)
	if err != nil {
		return err
	}
	c.Response().Header.Set(namespace.HeadersKey.Auth, token)
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
func VerifyClientEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	var client model.ClientMeta
	client.Email = email

	if err := lib.ValidatorV10.Var(client.Email, "email"); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			BadReq := lib.Error.General.BadRequest
			for _, fieldErr := range validationErrors {
				// You can customize the message
				BadReq.WithError(
					fmt.Errorf("field '%s' failed on the '%s' rule", fieldErr.Field(), fieldErr.Tag()),
				)
			}
			return BadReq
		} else {
			return lib.Error.General.InternalError.WithError(err)
		}
	}

	tx, end, err := database.ContextTransaction(c)
	defer end()
	if err != nil {
		return err
	}
	var clientFull model.ClientFull
	var exists string
	if err := tx.Model(&clientFull).
		Where("email = ?", email).
		Pluck("email", &exists).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if exists == "" {
		return lib.Error.Client.NotFound
	}

	sqlQuery := fmt.Sprintf("UPDATE %s SET verified = true WHERE email = ?", clientFull.TableName())

	if err := tx.Exec(sqlQuery, email).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// GetClientByEmail retrieves an client by email
//
//	@Summary		Get client by email
//	@Description	Retrieve an client by its email
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			email			path		string	true	"Client Email"
//	@Produce		json
//	@Success		200	{object}	DTO.Client
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/email/{email} [get]
func GetClientByEmail(c *fiber.Ctx) error {
	email := c.Params("email")

	if email == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing 'email' at params route"))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var client model.ClientMeta
	if err := tx.Model(&model.ClientFull{}).Where("email = ?", email).First(&client).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.Client.NotFound
		}
		return lib.Error.General.InternalError.WithError(err)
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &client, &DTO.Client{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// GetClientAppointments returns only the appointments of a client
//
//	@Summary		Get client appointments
//	@Description	Get only the appointments field from a client
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			Authorization	header	string	true	"X-Auth-Token"
//	@Param			id				path	string	true	"Client ID"
//	@Produce		json
//	@Success		200	{object}	dJSON.ClientAppointments
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/{id}/appointments [get]
func GetClientAppointments(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing ID on params route"))
	}

	tx, end, err := database.ContextTransaction(c)
	defer end()
	if err != nil {
		return err
	}

	var appointments mJSON.ClientAppointments
	if err := tx.Model(&model.ClientFull{}).
		Where("id = ?", id).
		Pluck("appointments", &appointments).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.Client.NotFound
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	return lib.ResponseFactory(c).Send(200, fiber.Map{
		"appointments": appointments,
	})
}

// UpdateClientById updates an client by ID
//
//	@Summary		Update client
//	@Description	Update an client
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string		true	"Client ID"
//	@Param			client	body		DTO.Client	true	"Client"
//	@Success		200		{object}	DTO.Client
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/client/{id} [patch]
func UpdateClientById(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing 'id' at params"))
	}

	changes := make(map[string]any)
	if err := c.BodyParser(&changes); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	tx, end, err := database.ContextTransaction(c)
	defer end()
	if err != nil {
		return err
	}

	var cFull model.ClientFull

	if err := tx.
		Model(&cFull).
		Where("id = ?", id).
		Omit(clause.Associations).
		Updates(changes).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &cFull, &DTO.Client{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// DeleteClientById deletes an client by ID
//
//	@Summary		Delete client
//	@Description	Delete an client
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Client ID"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		404	{object}	nil
//	@Router			/client/{id} [delete]
func DeleteClientById(c *fiber.Ctx) error {
	return DeleteOneById(c, &model.ClientFull{})
}

func Client(Gorm *handler.Gorm) {
	endpoint := &middleware.Endpoint{DB: Gorm}
	endpoint.BulkRegisterHandler([]fiber.Handler{
		CreateClient,
		LoginClient,
		VerifyClientEmail,
		GetClientByEmail,
		GetClientAppointments,
		UpdateClientById,
		DeleteClientById,
	})
}
