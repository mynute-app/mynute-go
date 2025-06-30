package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	dJSON "agenda-kaki-go/core/config/api/dto/json"
	database "agenda-kaki-go/core/config/db"
	"agenda-kaki-go/core/config/db/model"
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/middleware"
	"fmt"
	"net/url"

	"github.com/go-playground/validator/v10"
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
	tx, end, err := database.ContextTransaction(c)
	defer end(err)
	if err != nil {
		return err
	}
	client.Appointments = &mJSON.ClientAppointments{}
	if err := tx.Model(&model.Client{}).Create(&client).Error; err != nil {
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
	var err error
	var body DTO.LoginClient
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	tx, end, err := database.ContextTransaction(c)
	defer end(err)
	if err != nil {
		return err
	}

	var client model.ClientMeta
	if err := tx.Model(&model.Client{}).Where("email = ?", body.Email).First(&client).Error; err != nil {
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
	var err error
	email := c.Params("email")
	var client model.ClientMeta
	// Parse the email from the URL as it comes in the form of "john.clark%40gmail.com"
	email, err = url.QueryUnescape(email)
	if err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}
	client.Email = email
	if err := lib.ValidatorV10.Var(client.Email, "email"); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("email invalid"))
		} else {
			return lib.Error.General.InternalError.WithError(err)
		}
	}

	tx, end, err := database.ContextTransaction(c)
	defer end(err)
	if err != nil {
		return err
	}
	var clientFull model.Client
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
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
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
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Param			id				path		string	true	"Client ID"
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

	tx, end, err := database.ContextTransaction(c)
	defer end(err)
	if err != nil {
		return err
	}

	var appointments mJSON.ClientAppointments
	if err := tx.Model(&model.Client{}).
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

	if err := UpdateOneById(c, &client); err != nil {
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
//	@Param			client	body		DTO.UpdateClientImages	true	"Client Design Images"
//	@Success		200		{object}	DTO.Client
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/client/{id}/design/images [patch]
func UpdateClientImages(c *fiber.Ctx) error {
	var err error

	tx, end, err := database.ContextTransaction(c)
	defer end(err)
	if err != nil {
		return err
	}

	image_type := c.Params("image_type")
	img_types_allowed := map[string]bool{"picture": true}

	allowed, ok := img_types_allowed[image_type]
	if !ok {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("image_type not allowed: %s", image_type))
	}
	if !allowed {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("image_type not allowed: %s", image_type))
	}

	var client model.Client
	id := c.Params("id")
	if err := tx.First(&client, "id = ?", id).Error; err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("failed to find client (%s): %w", id, err))
	}

	var uploaded_img_types = make([]string, 0)

	defer func() {
		r := recover()
		if r != nil || err != nil {
			for _, img_type := range uploaded_img_types {
				_ = client.Design.Images.Delete(img_type, client.TableName(), client.ID.String())
			}
		}
	}()

	for img_type := range img_types_allowed {
		if c.FormValue(img_type) == "" {
			continue
		}
		file, err := c.FormFile(img_type)
		if err != nil {
			continue
		}
		_, err = client.Design.Images.Save(img_type, client.TableName(), client.ID.String(), file)
		if err != nil {
			return err
		}
		uploaded_img_types = append(uploaded_img_types, img_type)
	}

	if err = tx.Save(&client).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return lib.ResponseFactory(c).SendDTO(200, &client.Design.Images, &dJSON.Images{})
}

// DeleteClientImages deletes the design images of an client
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
//	@Success		200	{object}	nil
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/client/{id}/design/images/{image_type} [delete]
func DeleteClientImages(c *fiber.Ctx) error {
	var err error
	tx, end, err := database.ContextTransaction(c)
	defer end(err)
	if err != nil {
		return err
	}

	image_type := c.Params("image_type")
	img_types_allowed := map[string]bool{"picture": true}

	allowed, ok := img_types_allowed[image_type]
	if !ok {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("image_type not allowed: %s", image_type))
	}
	if !allowed {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("image_type not allowed: %s", image_type))
	}

	var client model.Client
	id := c.Params("id")
	if err := tx.First(&client, "id = ?", id).Error; err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("failed to find client (%s): %w", id, err))
	}

	if err := client.Design.Images.Delete(image_type, client.TableName(), id); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := tx.Save(&client).Error; err != nil {
		return err
	}

	return lib.ResponseFactory(c).SendDTO(200, &client.Design.Images, &dJSON.Images{})
}

func Client(Gorm *handler.Gorm) {
	endpoint := &middleware.Endpoint{DB: Gorm}
	endpoint.BulkRegisterHandler([]fiber.Handler{
		CreateClient,
		LoginClient,
		VerifyClientEmail,
		GetClientByEmail,
		GetClientById,
		GetClientAppointments,
		UpdateClientById,
		DeleteClientById,
	})
}
