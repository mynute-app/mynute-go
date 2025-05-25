package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v2"
)

// CreateService creates a service
//
//	@Summary		Create service
//	@Description	Create a service
//	@Tags			Service
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Accept			json
//	@Produce		json
//	@Param			service	body		DTO.CreateService	true	"Service"
//	@Success		201		{object}	DTO.Service
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
//	@Param			Authorization	header		string	true	"X-Auth-Token"
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
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Service ID"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		404	{object}	nil
//	@Router			/service/{id} [delete]
func DeleteServiceById(c *fiber.Ctx) error {
	return DeleteOneById(c, &model.Service{})
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
	})
}
