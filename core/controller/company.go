package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	database "agenda-kaki-go/core/config/db"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/middleware"
	"fmt"
	"net/url"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	// Must contain only lowercase letters, numbers, and hyphens
	if err := lib.ValidatorV10.Var(body.StartSubdomain, "mySubdomainValidation"); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			BadReq := lib.Error.General.BadRequest
			for _, fieldErr := range validationErrors {
				// You can customize the message
				BadReq.WithError(
					fmt.Errorf("field '%s' failed on the '%s' rule", fieldErr.Field(), fieldErr.Tag()),
				)
			}
			return BadReq
		} else {
			return lib.Error.General.InternalError.WithError(err)
		}
	}

	var company model.Company

	company.LegalName = body.LegalName
	company.TradeName = body.TradeName
	company.TaxID = body.TaxID

	if err := company.Create(tx); err != nil {
		return err
	}

	var domain model.Subdomain
	domain.Name = body.StartSubdomain
	domain.CompanyID = company.ID

	if err := company.AddSubdomain(tx, &domain); err != nil {
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

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	id := c.Params("id")

	if id == "" {
		return lib.Error.General.NotFoundError.WithError(fmt.Errorf("parameter 'id' not found on route parameters"))
	}

	if err := tx.Where("id = ?", id).First(&company).Error; err != nil {
		return lib.Error.Company.NotFound.WithError(err)
	}

	if full_c, err := company.GetFullCompany(tx); err != nil {
		return lib.Error.General.UpdatedError.WithError(err)
	} else {
		if err := lib.ResponseFactory(c).SendDTO(200, full_c, &DTO.Company{}); err != nil {
			return lib.Error.General.InternalError.WithError(err)
		}
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

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	name := c.Params("name")

	if name == "" {
		return lib.Error.General.NotFoundError.WithError(fmt.Errorf("parameter 'name' not found on route parameters"))
	}

	clearName, err := url.QueryUnescape(name)
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := tx.Where("legal_name = ? OR trade_name = ?", clearName, clearName).First(&company).Error; err != nil {
		return lib.Error.Company.NotFound.WithError(err)
	}

	if full_c, err := company.GetFullCompany(tx); err != nil {
		return lib.Error.General.UpdatedError.WithError(err)
	} else {
		if err := lib.ResponseFactory(c).SendDTO(200, full_c, &DTO.Company{}); err != nil {
			return lib.Error.General.InternalError.WithError(err)
		}
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

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	tax_id := c.Params("tax_id")

	if tax_id == "" {
		return lib.Error.General.NotFoundError.WithError(fmt.Errorf("parameter 'tax_id' not found on route parameters"))
	}

	if err := tx.Where("tax_id = ?", tax_id).First(&company).Error; err != nil {
		return lib.Error.Company.NotFound.WithError(err)
	}

	if full_c, err := company.GetFullCompany(tx); err != nil {
		return lib.Error.General.UpdatedError.WithError(err)
	} else {
		if err := lib.ResponseFactory(c).SendDTO(200, full_c, &DTO.Company{}); err != nil {
			return lib.Error.General.InternalError.WithError(err)
		}
	}

	return nil
}

// GetCompanyBySubdomain retrieves a company by subdomain
//
//	@Summary		Get company by subdomain
//	@Description	Retrieve a company by its subdomain
//	@Tags			Company
//	@Param			subdomain_name	path	string	true	"Subdomain Name"
//	@Produce		json
//	@Success		200	{object}	DTO.Company
//	@Failure		404	{object}	DTO.ErrorResponse
//	@Router			/company/subdomain/{subdomain_name} [get]
func GetCompanyBySubdomain(c *fiber.Ctx) error {
	var company model.Company

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	subdomain_name := c.Params("subdomain_name")

	if subdomain_name == "" {
		return lib.Error.General.NotFoundError.WithError(fmt.Errorf("parameter 'subdomain_name' not found on route parameters"))
	}

	var subdomain model.Subdomain

	if err := lib.ChangeToPublicSchema(tx); err != nil {
		return err
	}

	if err := tx.
		Where("name = ?", subdomain_name).
		First(&subdomain).
		Error; err != nil {
		return lib.Error.Company.NotFound.WithError(err)
	}

	company.ID = subdomain.CompanyID

	if full_c, err := company.GetFullCompany(tx); err != nil {
		return lib.Error.General.UpdatedError.WithError(err)
	} else {
		if err := lib.ResponseFactory(c).SendDTO(200, full_c, &DTO.Company{}); err != nil {
			return lib.Error.General.InternalError.WithError(err)
		}
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

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	if full_c, err := company.GetFullCompany(tx); err != nil {
		return lib.Error.General.UpdatedError.WithError(err)
	} else {
		if err := lib.ResponseFactory(c).SendDTO(200, full_c, &DTO.Company{}); err != nil {
			return lib.Error.General.InternalError.WithError(err)
		}
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
	var company model.Company

	company_id := c.Params("id")

	company_uuid, err := uuid.Parse(company_id)
	if err != nil {
		return lib.Error.Company.NotFound.WithError(err)
	}

	tx, end, err := database.ContextTransaction(c)
	defer end()
	if err != nil {
		return err
	}

	company.ID = company_uuid

	if err := company.Delete(tx); err != nil {
		return lib.Error.Company.NotFound.WithError(err)
	}

	return nil
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
