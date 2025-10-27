package controller

import (
	"fmt"
	DTO "mynute-go/core/src/config/api/dto"
	dJSON "mynute-go/core/src/config/api/dto/json"
	"mynute-go/core/src/config/db/model"
	"mynute-go/core/src/config/namespace"
	"mynute-go/core/src/handler"
	"mynute-go/core/src/lib"
	"mynute-go/core/src/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
	var err error
	var client model.Client
	if err := c.BodyParser(&client); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	tx, err := lib.Session(c)
	if err != nil {
		return err
	}
	if err := tx.Model(&model.Client{}).Create(&client).Error; err != nil {
		return err
	}
	if err := lib.ResponseFactory(c).SendDTO(200, &client, &DTO.Client{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

// LoginClientByPassword logs an client in
//
//	@Summary		Login
//	@Description	Log in an client using password
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Param			client	body	DTO.LoginClient	true	"Client"
//	@Success		200
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/login [post]
func LoginClientByPassword(c *fiber.Ctx) error {
	token, err := LoginByPassword(namespace.ClientKey.Name, &model.Client{}, c)
	if err != nil {
		return err
	}
	c.Response().Header.Set(namespace.HeadersKey.Auth, token)
	return nil
}

// LoginClientByEmailCode logs in a client using email and validation code
//
//	@Summary		Login client by email code
//	@Description	Login client using email and validation code
//	@Tags			client/auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body	DTO.LoginByEmailCode	true	"Login credentials"
//	@Success		200
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/login-with-code [post]
func LoginClientByEmailCode(c *fiber.Ctx) error {
	token, err := LoginByEmailCode(namespace.ClientKey.Name, &model.Client{}, c)
	if err != nil {
		return err
	}
	c.Response().Header.Set(namespace.HeadersKey.Auth, token)
	return nil
}

// SendClientLoginValidationCodeByEmail sends a login validation code to a client's email
//
//	@Summary		Send client login validation code by email
//	@Description	Send a login validation code to a client's email
//	@Tags			Client
//	@Param			email	path	string	true	"Client Email"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/send-login-code/email/{email} [post]
func SendClientLoginValidationCodeByEmail(c *fiber.Ctx) error {
	if err := SendLoginValidationCodeByEmail(c, &model.Client{}); err != nil {
		return err
	}
	return nil
}

// GetClientByEmail retrieves an client by email
//
//	@Summary		Get client by email
//	@Description	Retrieve an client by its email
//	@Tags			Client
//	@Param			email	path	string	true	"Client Email"
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
	if err := tx.Model(&model.Client{}).Where("email = ?", email).First(&client).Error; err != nil {
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

// GetClientById retrieves an client by ID
//
//	@Summary		Get client by ID
//	@Description	Retrieve an client by its ID
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			Authorization	header	string	true	"X-Auth-Token"
//	@Param			id				path	string	true	"Client ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Client
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/{id} [get]
func GetClientById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing 'id' at params route"))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var client model.ClientMeta
	if err := tx.Model(&model.Client{}).Where("id = ?", id).First(&client).Error; err != nil {
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
//	@Param			X-Auth-Token	header	string	true	"X-Auth-Token"
//	@Param			id				path	string	true	"Client ID"
//	@Produce		json
//	@Success		200	{object}	dJSON.ClientAppointments
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/{id}/appointments [get]
func GetClientAppointments(c *fiber.Ctx) error {
	var err error
	id := c.Params("id")
	if id == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing ID on params route"))
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var Appointments []model.ClientAppointment
	if err := tx.Model(&model.ClientAppointment{}).
		Where("client_id = ?", id).
		Find(&Appointments).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.Client.NotFound
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	return lib.ResponseFactory(c).Send(200, fiber.Map{
		"appointments": Appointments,
	})
}

// UpdateClientById updates an client by ID
//
//	@Summary		Update client
//	@Description	Update an client
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string		true	"Client ID"
//	@Param			client	body		DTO.Client	true	"Client"
//	@Success		200		{object}	DTO.Client
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/client/{id} [patch]
func UpdateClientById(c *fiber.Ctx) error {
	var client model.Client

	if err := UpdateOneById(c, &client, nil); err != nil {
		return err
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	if err := client.GetFullClient(tx); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &client, &DTO.Client{}); err != nil {
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
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Client ID"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		404	{object}	nil
//	@Router			/client/{id} [delete]
func DeleteClientById(c *fiber.Ctx) error {
	return DeleteOneById(c, &model.Client{})
}

// UpdateClientImages updates the design images of an client
//
//	@Summary		Update client design images
//	@Description	Update the design images of an client
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Client ID"
//	@Accept			json
//	@Produce		json
//	@Param			profile	formData	file	false	"Profile image"
//	@Success		200		{object}	DTO.Client
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/client/{id}/design/images [patch]
func UpdateClientImages(c *fiber.Ctx) error {
	img_types_allowed := map[string]bool{"profile": true}

	var client model.Client
	Design, err := UpdateImagesById(c, client.TableName(), &client, img_types_allowed)
	if err != nil {
		return err
	}

	return lib.ResponseFactory(c).SendDTO(200, &Design.Images, &dJSON.Images{})
}

// DeleteClientImage deletes the design images of an client
//
//	@Summary		Delete client design images
//	@Description	Delete the design images of an client
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Client ID"
//	@Param			image_type		path		string	true	"Image Type"
//	@Produce		json
//	@Success		200	{object}	dJSON.Images
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/{id}/design/images/{image_type} [delete]
func DeleteClientImage(c *fiber.Ctx) error {
	img_types_allowed := map[string]bool{"profile": true}
	var client model.Client
	Design, err := DeleteImageById(c, client.TableName(), &client, img_types_allowed)
	if err != nil {
		return err
	}

	return lib.ResponseFactory(c).SendDTO(200, &Design.Images, &dJSON.Images{})
}

// ResetClientPasswordByEmail resets the password of a client by email
//
//	@Summary		Reset client password by email
//	@Description	Reset the password of a client by its email
//	@Tags			Client
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			email			path		string	true	"Client Email"
//	@Query			lang								query		string	false	"Language code (default: en)"
//	@Produce		json
//	@Success		200	{object}	DTO.PasswordReseted
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/reset-password/{email} [post]
func ResetClientPasswordByEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("missing 'email' at params route"))
	}
	if err := SendNewPasswordByEmail(c, email, &model.Client{}); err != nil {
		return err
	}
	return lib.ResponseFactory(c).Http200(nil)
}

// SendClientVerificationCodeByEmail sends a verification code to a client's email
//
//	@Summary		Send client verification code by email
//	@Description	Send a verification code to a client's email
//	@Tags			Client
//	@Param			email	path	string	true	"Client Email"
//	@Query			lang	string	false	"Language for the verification email"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/send-verification-code/email/{email} [post]
func SendClientVerificationCodeByEmail(c *fiber.Ctx) error {
	return SendVerificationCodeByEmail(c, &model.Client{})
}

// VerifyClientEmail verifies a client's email
//
//	@Summary		Verify client email
//	@Description	Verify a client's email
//	@Tags			Client
//	@Param			email	path	string	true	"Client Email"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/verify-email/{email}/{code} [get]
func VerifyClientEmail(c *fiber.Ctx) error {
	return VerifyEmail(c, &model.Client{})
}

func Client(Gorm *handler.Gorm) {
	endpoint := &middleware.Endpoint{DB: Gorm}
	endpoint.BulkRegisterHandler([]fiber.Handler{
		CreateClient,
		LoginClientByPassword,
		LoginClientByEmailCode,
		SendClientLoginValidationCodeByEmail,
		ResetClientPasswordByEmail,
		GetClientByEmail,
		GetClientById,
		GetClientAppointments,
		UpdateClientById,
		DeleteClientById,
		UpdateClientImages,
		DeleteClientImage,
		SendClientVerificationCodeByEmail,
		VerifyClientEmail,
	})
}
