package controllers

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"
	"agenda-kaki-go/core/service"

	"github.com/gofiber/fiber/v2"
)

// serviceController embeds service.Base in order to extend it with the functions below
type serviceController struct {
	service.Base[models.Service, DTO.Service]
}

// CreateService creates a service
//
//	@Summary		Create service
//	@Description	Create a service
//	@Tags			Service
//	@Accept			json
//	@Produce		json
//	@Param			service	body		DTO.Service	true	"Service"
//	@Success		200		{object}	DTO.Service
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/service [post]
func (cc *serviceController) CreateService(c *fiber.Ctx) error {
	return cc.CreateOne(c)
}

// GetServiceById retrieves a service by ID
//
//	@Summary		Get service by ID
//	@Description	Retrieve a service by its ID
//	@Tags			Service
//	@Param			id	path	string	true	"Service ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Service
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/service/{id} [get]
func (cc *serviceController) GetServiceById(c *fiber.Ctx) error {
	return cc.GetBy("id", c)
}

// GetServiceByName retrieves a service by name
//
//	@Summary		Get service by name
//	@Description	Retrieve a service by its name
//	@Tags			Service
//	@Param			name	path	string	true	"Service Name"
//	@Produce		json
//	@Success		200	{object}	DTO.Service
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/service/name/{name} [get]
func (cc *serviceController) GetServiceByName(c *fiber.Ctx) error {
	return cc.GetBy("name", c)
}

// UpdateServiceById updates a service by ID
//
//	@Summary		Update service by ID
//	@Description	Update a service by its ID
//	@Tags			Service
//	@Param			id	path	string	true	"Service ID"
//	@Accept			json
//	@Produce		json
//	@Param			service	body		DTO.Service	true	"Service"
//	@Success		200		{object}	DTO.Service
//	@Failure		404		{object}	DTO.ErrorResponse
//	@Router			/service/{id} [patch]
func (cc *serviceController) UpdateServiceById(c *fiber.Ctx) error {
	return cc.UpdateOneById(c)
}

// DeleteServiceById deletes a service by ID
//
//	@Summary		Delete service by ID
//	@Description	Delete a service by its ID
//	@Tags			Service
//	@Param			id	path	string	true	"Service ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Service
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/service/{id} [delete]
func (cc *serviceController) DeleteServiceById(c *fiber.Ctx) error {
	return cc.DeleteOneById(c)
}

// Service returns a serviceController
func Service(Gorm *handlers.Gorm) *serviceController {
	return &serviceController{
		Base: service.Base[models.Service, DTO.Service]{
			Name:         namespace.UserKey.Name,
			Request:      handlers.Request(Gorm),
			Middleware:   middleware.Service(Gorm),
			Associations: []string{"ServiceType"},
		},
	}
}
