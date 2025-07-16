package controller

import (
	DTO "mynute-go/core/config/api/dto"
	dJSON "mynute-go/core/config/api/dto/json"
	"mynute-go/core/config/db/model"
	"mynute-go/core/handler"
	"mynute-go/core/lib"
	"mynute-go/core/middleware"

	"github.com/gofiber/fiber/v2"
)

// CreateService creates a service
//
//	@Summary		Create service
//	@Description	Create a service
//	@Tags			Service
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Accept			json
//	@Produce		json
//	@Param			service	body		DTO.CreateService	true	"Service"
//	@Success		200		{object}	DTO.Service
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/service [post]
func CreateService(c *fiber.Ctx) error {
	var service model.Service
	if err := Create(c, &service); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &service, &DTO.Service{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// GetServiceById retrieves a service by ID
//
//	@Summary		Get service by ID
//	@Description	Retrieve a service by its ID
//	@Tags			Service
//	@Param			id	path	string	true	"Service ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Service
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/service/{id} [get]
func GetServiceById(c *fiber.Ctx) error {
	var service model.Service

	if err := GetOneBy("id", c, &service); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &service, &DTO.Service{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// GetServiceByName retrieves a service by name
//
//	@Summary		Get service by name
//	@Description	Retrieve a service by its name
//	@Tags			Service
//	@Param			name	path	string	true	"Service Name"
//	@Produce		json
//	@Success		200	{object}	DTO.Service
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/service/name/{name} [get]
func GetServiceByName(c *fiber.Ctx) error {
	var service model.Service

	if err := GetOneBy("name", c, &service); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &service, &DTO.Service{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// UpdateServiceById updates a service by ID
//
//	@Summary		Update service by ID
//	@Description	Update a service by its ID
//	@Tags			Service
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Service ID"
//	@Accept			json
//	@Produce		json
//	@Param			service	body		DTO.Service	true	"Service"
//	@Success		200		{object}	DTO.Service
//	@Failure		404		{object}	nil
//	@Router			/service/{id} [patch]
func UpdateServiceById(c *fiber.Ctx) error {
	var service model.Service

	if err := UpdateOneById(c, &service); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &service, &DTO.Service{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// DeleteServiceById deletes a service by ID
//
//	@Summary		Delete service by ID
//	@Description	Delete a service by its ID
//	@Tags			Service
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Service ID"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		404	{object}	nil
//	@Router			/service/{id} [delete]
func DeleteServiceById(c *fiber.Ctx) error {
	return DeleteOneById(c, &model.Service{})
}

// UpdateServiceImages updates images of a service
//
//	@Summary		Update service images
//	@Description	Update images of a service
//	@Tags			Service
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Service ID"
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			images	formData	dJSON.Images	true	"Images"
//	@Success		200		{object}	dJSON.Images
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/service/{id}/design/images [patch]
func UpdateServiceImages(c *fiber.Ctx) error {
	img_types_allowed := map[string]bool{"profile": true}

	var service model.Service
	Design, err := UpdateImagesById(c, service.TableName(), &service, img_types_allowed)
	if err != nil {
		return err
	}

	return lib.ResponseFactory(c).SendDTO(200, &Design.Images, &dJSON.Images{})
}

// DeleteServiceImage deletes images of a service
//
//	@Summary		Delete service images
//	@Description	Delete images of a service
//	@Tags			Service
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Service ID"
//	@Produce		json
//	@Success		200	{object}	dJSON.Images
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/service/{id}/design/images/{image_type} [delete]
func DeleteServiceImage(c *fiber.Ctx) error {
	img_types_allowed := map[string]bool{"profile": true}

	var service model.Service
	Design, err := DeleteImageById(c, service.TableName(), &service, img_types_allowed)
	if err != nil {
		return err
	}

	return lib.ResponseFactory(c).SendDTO(200, &Design.Images, &dJSON.Images{})
}

// Service returns a service_controller
func Service(Gorm *handler.Gorm) {
	endpoint := &middleware.Endpoint{DB: Gorm}
	endpoint.BulkRegisterHandler([]fiber.Handler{
		CreateService,
		GetServiceById,
		GetServiceByName,
		UpdateServiceById,
		DeleteServiceById,
		UpdateServiceImages,
		DeleteServiceImage,
	})
}
