package controllers

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v2"
)

type companyType struct {
	BaseController[models.CompanyType, DTO.CompanyType]
}

// CreateCompanyType creates a company type
//	@Summary		Create company type
//	@Description	Create a company type
//	@Tags			CompanyType
//	@Accept			json
//	@Produce		json
//	@Param			company_type	body		DTO.CompanyType	true	"Company Type"
//	@Success		200				{object}	DTO.CompanyType
//	@Failure		400				{object}	DTO.ErrorResponse
//	@Router			/company_type [post]
func (cc *companyType) CreateCompanyType(c *fiber.Ctx) error {
	return cc.CreateOne(c)
}

// GetCompanyTypeByName retrieves a company type by ID
//	@Summary		Get company type by ID
//	@Description	Retrieve a company type by its ID
//	@Tags			CompanyType
//	@Param			id	path	string	true	"Company Type ID"
//	@Produce		json
//	@Success		200	{object}	DTO.CompanyType
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/company_type/{id} [get]
func (cc *companyType) GetCompanyTypeByName(c *fiber.Ctx) error {
	return cc.GetBy("name", c)
}

// GetCompanyTypeById retrieves a company type by ID
//	@Summary		Get company type by ID
//	@Description	Retrieve a company type by its ID
//	@Tags			CompanyType
//	@Param			id	path	string	true	"Company Type ID"
//	@Produce		json
//	@Success		200	{object}	DTO.CompanyType
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/company_type/{id} [get]
func (cc *companyType) GetCompanyTypeById(c *fiber.Ctx) error {
	return cc.GetBy("id", c)
}

// UpdateCompanyTypeById updates a company type by ID
//	@Summary		Update company type by ID
//	@Description	Update a company type by its ID
//	@Tags			CompanyType
//	@Param			id	path	string	true	"Company Type ID"
//	@Accept			json
//	@Produce		json
//	@Param			company_type	body		DTO.CompanyType	true	"Company Type"
//	@Success		200				{object}	DTO.CompanyType
//	@Failure		404				{object}	DTO.ErrorResponse
//	@Router			/company_type/{id} [patch]
func (cc *companyType) UpdateCompanyTypeById(c *fiber.Ctx) error {
	return cc.UpdateOneById(c)
}

// DeleteCompanyTypeById deletes a company type by ID
//	@Summary		Delete company type by ID
//	@Description	Delete a company type by its ID
//	@Tags			CompanyType
//	@Param			id	path	string	true	"Company Type ID"
//	@Produce		json
//	@Success		200	{object}	DTO.CompanyType
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/company_type/{id} [delete]
func (cc *companyType) DeleteCompanyTypeById(c *fiber.Ctx) error {
	return cc.DeleteOneById(c)
}

func CompanyType(Gorm *handlers.Gorm) *companyType {
	return &companyType{
		BaseController: BaseController[models.CompanyType, DTO.CompanyType]{
			Name:         namespace.CompanyTypeKey.Name,
			Request:      handlers.Request(Gorm),
			Middleware:   middleware.CompanyType(Gorm),
			Associations: []string{},
		},
	}
}
