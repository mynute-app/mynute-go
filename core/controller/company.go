package controller

import "fmt"

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/service"

	"github.com/gofiber/fiber/v2"
)

// company_controller embeds service.Base in order to extend it with the functions below
type company_controller struct {
	service.Base[model.Company, DTO.Company]
}

// CreateCompany creates a company
//
//	@Summary		Create company
//	@Description	Create a company
//	@Tags			Company
//	@Accept			json
//	@Produce		json
//	@Param			company	body		DTO.CreateCompany	true	"Company"
//	@Success		200		{object}	DTO.Company
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/company [post]
func (cc *company_controller) CreateCompany(c *fiber.Ctx) error {
	res := &lib.SendResponse{Ctx: c}
	var body DTO.CreateCompany
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	var company model.Company

	company.Name = body.Name
	company.TaxID = body.TaxID

	if err := cc.Request.Gorm.DB.Create(&company).Error; err != nil {
		return err
	}

	var owner model.Employee

	owner.Name = body.OwnerName
	owner.Surname = body.OwnerSurname
	owner.Email = body.OwnerEmail
	owner.Phone = body.OwnerPhone
	owner.Password = body.OwnerPassword
	owner.CompanyID = company.ID

	if err := cc.Request.Gorm.DB.Create(&owner).Error; err != nil {
		return err
	}

	if err := cc.Request.Gorm.GetOneBy("id", fmt.Sprint(company.ID), &company, cc.Associations); err != nil {
		return err
	}

	if err := res.SendDTO(200, &company, &DTO.Company{}); err != nil {
		return err
	}
	return nil
}

// GetOneById retrieves a company by ID
//
//	@Summary		Get company by ID
//	@Description	Retrieve a company by its ID
//	@Tags			Company
//	@Param			id	path	string	true	"Company ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Company
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/company/{id} [get]
func (cc *company_controller) GetCompanyById(c *fiber.Ctx) error {
	return cc.GetBy("id", c)
}

// GetOneByName retrieves a company by name
//
//	@Summary		Get company by name
//	@Description	Retrieve a company by its name
//	@Tags			Company
//	@Param			name	path	string	true	"Company Name"
//	@Produce		json
//	@Success		200	{object}	DTO.Company
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/company/name/{name} [get]
func (cc *company_controller) GetCompanyByName(c *fiber.Ctx) error {
	return cc.GetBy("name", c)
}

// GetOneByTaxId retrieves a company by tax ID
//
//	@Summary		Get company by tax ID
//	@Description	Retrieve a company by its tax identification number
//	@Tags			Company
//	@Param			tax_id	path	string	true	"Company Tax ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Company
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/company/tax_id/{tax_id} [get]
func (cc *company_controller) GetCompanyByTaxId(c *fiber.Ctx) error {
	return cc.GetBy("tax_id", c)
}

// UpdateCompanyById updates a company by ID
//
//	@Summary		Update company by ID
//	@Description	Update a company by its ID
//	@Tags			Company
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Company ID"
//	@Accept			json
//	@Produce		json
//	@Param			company	body		DTO.Company	true	"Company"
//	@Success		201		{object}	DTO.Company
//	@Failure		404		{object}	DTO.ErrorResponse
//	@Router			/company/{id} [patch]
func (cc *company_controller) UpdateCompanyById(c *fiber.Ctx) error {
	return cc.UpdateOneById(c)
}

// DeleteCompanyById deletes a company by ID
//
//	@Summary		Delete company by ID
//	@Description	Delete a company by its ID
//	@Tags			Company
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"Authorization"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Company ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Company
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/company/{id} [delete]
func (cc *company_controller) DeleteCompanyById(c *fiber.Ctx) error {
	return cc.DeleteOneById(c)
}

// Constructor for company_controller
func Company(Gorm *handler.Gorm) *company_controller {
	return &company_controller{
		Base: service.Base[model.Company, DTO.Company]{
			Name:         namespace.CompanyKey.Name,
			Request:      handler.Request(Gorm),
			Associations: []string{"Sector", "Branches", "Employees", "Services"},
		},
	}
}
