package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	database "agenda-kaki-go/core/config/db"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v2"
)

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
func CreateCompany(c *fiber.Ctx) error {
	tx, end, err := database.ContextTransaction(c)
	defer end()
	if err != nil {
		return err
	}
	var body DTO.CreateCompany
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	var company model.Company

	company.Name = body.Name
	company.TaxID = body.TaxID

	if err := company.Create(tx); err != nil {
		return err
	}

	var owner model.Employee

	owner.Name = body.OwnerName
	owner.Surname = body.OwnerSurname
	owner.Email = body.OwnerEmail
	owner.Phone = body.OwnerPhone
	owner.Password = body.OwnerPassword
	owner.CompanyID = company.ID

	if err := company.CreateOwner(tx, &owner); err != nil {
		return err
	}

	if fullCompany, err := company.GetFullCompany(tx); err != nil {
		return err
	} else {
		if err := lib.ResponseFactory(c).SendDTO(200, fullCompany, &DTO.Company{}); err != nil {
			return lib.Error.General.InternalError.WithError(err)
		}
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
func GetCompanyById(c *fiber.Ctx) error {
	var company model.Company
	if err := GetOneBy("id", c, &company); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &company, &DTO.Company{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
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
func GetCompanyByName(c *fiber.Ctx) error {
	var company model.Company

	if err := GetOneBy("name", c, &company); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &company, &DTO.Company{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
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
func GetCompanyByTaxId(c *fiber.Ctx) error {
	var company model.Company

	if err := GetOneBy("tax_id", c, &company); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &company, &DTO.Company{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil

}

// UpdateCompanyById updates a company by ID
//
//	@Summary		Update company by ID
//	@Description	Update a company by its ID
//	@Tags			Company
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Company ID"
//	@Accept			json
//	@Produce		json
//	@Param			company	body		DTO.Company	true	"Company"
//	@Success		201		{object}	DTO.Company
//	@Failure		404		{object}	DTO.ErrorResponse
//	@Router			/company/{id} [patch]
func UpdateCompanyById(c *fiber.Ctx) error {
	var company model.Company
	if err := UpdateOneById(c, &company); err != nil {
		return err
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &company, &DTO.Company{}); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

// DeleteCompanyById deletes a company by ID
//
//	@Summary		Delete company by ID
//	@Description	Delete a company by its ID
//	@Tags			Company
//	@Security		ApiKeyAuth
//	@Param			Authorization	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			id				path		string	true	"Company ID"
//	@Produce		json
//	@Success		200	{object}	DTO.Company
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/company/{id} [delete]
func DeleteCompanyById(c *fiber.Ctx) error {
	return DeleteOneById(c, &model.Company{})
}

// Constructor for company_controller
func Company(Gorm *handler.Gorm) {
	endpoint := &middleware.Endpoint{DB: Gorm}
	endpoint.BulkRegisterHandler([]fiber.Handler{
		CreateCompany,
		GetCompanyById,
		GetCompanyByName,
		GetCompanyByTaxId,
		UpdateCompanyById,
		DeleteCompanyById,
	})
}
