package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/middleware"
	"agenda-kaki-go/core/service"

	"github.com/gofiber/fiber/v2"
)

type sector_controller struct {
	service.Base[model.Sector, DTO.Sector]
}

// CreateCompanyType creates a company type
//
//	@Summary		Create company type
//	@Description	Create a company type
//	@Tags			Sector
//	@Accept			json
//	@Produce		json
//	@Param			sector	body		DTO.Sector	true	"Company Type"
//	@Success		200				{object}	DTO.Sector
//	@Failure		400				{object}	DTO.ErrorResponse
//	@Router			/sector [post]
func (cc *sector_controller) CreateCompanyType(c *fiber.Ctx) error {
	return cc.CreateOne(c)
}

// GetCompanyTypeByName retrieves a company type by ID
//
//	@Summary		Get company type by ID
//	@Description	Retrieve a company type by its ID
//	@Tags			Sector
//	@Param			id	path	string	true	"Company Type ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Sector
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/sector/{id} [get]
func (cc *sector_controller) GetCompanyTypeByName(c *fiber.Ctx) error {
	return cc.GetBy("name", c)
}

// GetCompanyTypeById retrieves a company type by ID
//
//	@Summary		Get company type by ID
//	@Description	Retrieve a company type by its ID
//	@Tags			Sector
//	@Param			id	path	string	true	"Company Type ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Sector
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/sector/{id} [get]
func (cc *sector_controller) GetCompanyTypeById(c *fiber.Ctx) error {
	return cc.GetBy("id", c)
}

// UpdateCompanyTypeById updates a company type by ID
//
//	@Summary		Update company type by ID
//	@Description	Update a company type by its ID
//	@Tags			Sector
//	@Param			id	path	string	true	"Company Type ID"
//	@Accept			json
//	@Produce		json
//	@Param			sector	body		DTO.Sector	true	"Company Type"
//	@Success		200				{object}	DTO.Sector
//	@Failure		404				{object}	DTO.ErrorResponse
//	@Router			/sector/{id} [patch]
func (cc *sector_controller) UpdateCompanyTypeById(c *fiber.Ctx) error {
	return cc.UpdateOneById(c)
}

// DeleteCompanyTypeById deletes a company type by ID
//
//	@Summary		Delete company type by ID
//	@Description	Delete a company type by its ID
//	@Tags			Sector
//	@Param			id	path	string	true	"Company Type ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Sector
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/sector/{id} [delete]
func (cc *sector_controller) DeleteCompanyTypeById(c *fiber.Ctx) error {
	return cc.DeleteOneById(c)
}

func Sector(Gorm *handler.Gorm) *sector_controller {
	return &sector_controller{
		Base: service.Base[model.Sector, DTO.Sector]{
			Name:         namespace.CompanyTypeKey.Name,
			Request:      handler.Request(Gorm),
			Middleware:   middleware.Sector(Gorm),
			Associations: []string{},
		},
	}
}
