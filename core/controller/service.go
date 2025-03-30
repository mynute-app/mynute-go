package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/service"

	"github.com/gofiber/fiber/v2"
)

// service_controller embeds service.Base in order to extend it with the functions below
type service_controller struct {
	service.Base[model.Service, DTO.Service]
}

// CreateService creates a service
//
//	@Summary		Create service
//	@Description	Create a service
//	@Tags			Service
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization"
//	@Failure		401				{object}	nil
//	@Accept			json
//	@Produce		json
//	@Param			service	body		DTO.CreateService	true	"Service"
//	@Success		201		{object}	DTO.Service
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/service [post]
func (cc *service_controller) CreateService(c *fiber.Ctx) error {
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
func (cc *service_controller) GetServiceById(c *fiber.Ctx) error {
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
func (cc *service_controller) GetServiceByName(c *fiber.Ctx) error {
	return cc.GetBy("name", c)
}

// UpdateServiceById updates a service by ID
//
//	@Summary		Update service by ID
//	@Description	Update a service by its ID
//	@Tags			Service
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Service ID"
//	@Accept			json
//	@Produce		json
//	@Param			service	body		DTO.Service	true	"Service"
//	@Success		200		{object}	DTO.Service
//	@Failure		404		{object}	nil
//	@Router			/service/{id} [patch]
func (cc *service_controller) UpdateServiceById(c *fiber.Ctx) error {
	return cc.UpdateOneById(c)
}

// DeleteServiceById deletes a service by ID
//
//	@Summary		Delete service by ID
//	@Description	Delete a service by its ID
//	@Tags			Service
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Service ID"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		404	{object}	nil
//	@Router			/service/{id} [delete]
func (cc *service_controller) DeleteServiceById(c *fiber.Ctx) error {
	return cc.DeleteOneById(c)
}

// Service returns a service_controller
func Service(Gorm *handler.Gorm) *service_controller {
	sc := &service_controller{
		Base: service.Base[model.Service, DTO.Service]{
			Name:         namespace.ClientKey.Name,
			Request:      handler.Request(Gorm),
			Associations: []string{"Company", "Branches", "Employees"},
		},
	}
	route := &handler.Route{DB: Gorm.DB}
	ServiceResources := []*handler.Resource{
		{
			Path:        "/service",
			Method:      "POST",
			Handler:     sc.CreateService,
			Description: "Create a service",
			Access:      "private",
		},
		{
			Path:        "/service/:id",
			Method:      "GET",
			Handler:     sc.GetServiceById,
			Description: "Get service by ID",
			Access:      "private",
		},
		{
			Path:        "/service/name/:name",
			Method:      "GET",
			Handler:     sc.GetServiceByName,
			Description: "Get service by name",
			Access:      "public",
		},
		{
			Path:        "/service/:id",
			Method:      "PATCH",
			Handler:     sc.UpdateServiceById,
			Description: "Update service by ID",
			Access:      "private",
		},
		{
			Path:        "/service/:id",
			Method:      "DELETE",
			Handler:     sc.DeleteServiceById,
			Description: "Delete service by ID",
			Access:      "private",
		},
	}
	route.BulkRegisterAndSave(ServiceResources)
	return sc
}
