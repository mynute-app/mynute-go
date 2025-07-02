package controller

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	dJSON "agenda-kaki-go/core/config/api/dto/json"
	database "agenda-kaki-go/core/config/db"
	"agenda-kaki-go/core/config/db/model"
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/middleware"
	"fmt"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateCompany creates a company
//
//	@Summary		Create company
//	@Description	Create a company
//	@Tags			Company
//	@Accept			json
//	@Produce		json
//	@Param			company	body		DTO.CreateCompany	true	"Company"
//	@Success		200		{object}	DTO.CompanyFull
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/company [post]
func CreateCompany(c *fiber.Ctx) error {
	tx, end, err := database.ContextTransaction(c)
	if err != nil {
		return err
	}
	var body DTO.CreateCompany
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	if err := lib.MyCustomStructValidator(body); err != nil {
		return err
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
		end(err)
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
		end(err)
		return err
	}

	if fullCompany, err := company.GetFullCompany(tx); err != nil {
		end(err)
		return err
	} else {
		if err := lib.ResponseFactory(c).SendDTO(200, fullCompany, &DTO.CompanyFull{}); err != nil {
			end(err)
			return lib.Error.General.InternalError.WithError(err)
		}
	}

	end(nil)

	return nil
}

// GetCompanyById retrieves a company by ID
//
//	@Summary		Get company by ID
//	@Description	Retrieve a company by its ID
//	@Tags			Company
//	@Param			id	path	string	true	"Company ID"
//	@Produce		json
//	@Success		200	{object}	DTO.CompanyFull
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/company/{id} [get]
func GetCompanyById(c *fiber.Ctx) error {
	var company model.Company

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	id := c.Params("id")

	if id == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("parameter 'id' not found on route parameters"))
	}

	if err := tx.Where("id = ?", id).First(&company).Error; err != nil {
		return lib.Error.Company.NotFound.WithError(err)
	}

	if full_c, err := company.GetFullCompany(tx); err != nil {
		return lib.Error.General.UpdatedError.WithError(err)
	} else {
		if err := lib.ResponseFactory(c).SendDTO(200, full_c, &DTO.CompanyFull{}); err != nil {
			return lib.Error.General.InternalError.WithError(err)
		}
	}

	return nil
}

// GetCompanyByName retrieves a company by name
//
//	@Summary		Get company by name
//	@Description	Retrieve a company by its name
//	@Tags			Company
//	@Param			name	path	string	true	"Company Name"
//	@Produce		json
//	@Success		200	{object}	DTO.CompanyFull
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/company/name/{name} [get]
func GetCompanyByName(c *fiber.Ctx) error {
	var company model.Company

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	name := c.Params("name")

	if name == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("parameter 'name' not found on route parameters"))
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
		if err := lib.ResponseFactory(c).SendDTO(200, full_c, &DTO.CompanyFull{}); err != nil {
			return lib.Error.General.InternalError.WithError(err)
		}
	}

	return nil
}

// CheckIfCompanyExistsByTaxID checks if a company exists by its tax ID
//
//	@Summary		Check if company exists by tax ID
//	@Description	Check if a company exists by its tax identification number
//	@Tags			Company
//	@Param			tax_id	path	string	true	"Company Tax ID"
//	@Produce		json
//	@Success		200	{object}	bool
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/company/tax_id/{tax_id}/exists [get]
func CheckIfCompanyExistsByTaxID(c *fiber.Ctx) error {
	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	tax_id := c.Params("tax_id")

	if tax_id == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("parameter 'tax_id' not found on route parameters"))
	}

	var company model.Company

	if err := tx.First(&company, "tax_id = ?", tax_id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.ResponseFactory(c).Send(404, nil)
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	return lib.ResponseFactory(c).Send(200, nil)
}

// GetCompanyByTaxId retrieves a company by tax ID
//
//	@Summary		Get company by tax ID
//	@Description	Retrieve a company by its tax identification number
//	@Tags			Company
//	@Param			tax_id	path	string	true	"Company Tax ID"
//	@Produce		json
//	@Success		200	{object}	DTO.CompanyFull
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/company/tax_id/{tax_id} [get]
func GetCompanyByTaxId(c *fiber.Ctx) error {
	var company model.Company

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	tax_id := c.Params("tax_id")

	if tax_id == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("parameter 'tax_id' not found on route parameters"))
	}

	if err := tx.First(&company, "tax_id = ?", tax_id).Error; err != nil {
		return lib.Error.Company.NotFound.WithError(err)
	}

	if full_c, err := company.GetFullCompany(tx); err != nil {
		return lib.Error.General.UpdatedError.WithError(err)
	} else {
		if err := lib.ResponseFactory(c).SendDTO(200, full_c, &DTO.CompanyFull{}); err != nil {
			return lib.Error.General.InternalError.WithError(err)
		}
	}

	return nil
}

// GetCompanyBySubdomain retrieves a company by subdomain
//
//	@Summary		Get company ID by subdomain
//	@Description	Retrieve a company by its subdomain
//	@Tags			Company
//	@Param			subdomain_name	path	string	true	"Subdomain Name"
//	@Produce		json
//	@Success		200	{object}	DTO.CompanyBase
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/company/subdomain/{subdomain_name} [get]
func GetCompanyBySubdomain(c *fiber.Ctx) error {
	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	subdomain_name := c.Params("subdomain_name")

	if subdomain_name == "" {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("parameter 'subdomain_name' not found on route parameters"))
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

	var company model.Company

	company.ID = subdomain.CompanyID

	if err := company.Refresh(tx); err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := lib.ResponseFactory(c).SendDTO(200, &company, &DTO.CompanyBase{}); err != nil {
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
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Company ID"
//	@Accept			json
//	@Produce		json
//	@Param			company	body		DTO.CreateCompany	true	"Company"
//	@Success		200		{object}	DTO.CompanyFull
//	@Failure		404		{object}	DTO.ErrorResponse
//	@Router			/company/{id} [patch]
func UpdateCompanyById(c *fiber.Ctx) error {
	var company model.Company

	if err := lib.ChangeToPublicSchemaByContext(c); err != nil {
		return err
	}

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
		if err := lib.ResponseFactory(c).SendDTO(200, full_c, &DTO.CompanyFull{}); err != nil {
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
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Company ID"
//	@Produce		json
//	@Success		200	{object}	DTO.CompanyFull
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Router			/company/{id} [delete]
func DeleteCompanyById(c *fiber.Ctx) error {
	var company model.Company

	company_id := c.Params("id")

	company_uuid, err := uuid.Parse(company_id)
	if err != nil {
		return lib.Error.Company.NotFound.WithError(err)
	}

	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	company.ID = company_uuid

	if err := company.Delete(tx); err != nil {
		return lib.Error.Company.NotFound.WithError(err)
	}

	return nil
}

// @Summary		Update company design images
// @Description	Upload and update design images (logo, banner, etc.)
// @Tags			Company
// @Accept			multipart/form-data
// @Produce		json
// @Security		ApiKeyAuth
// @Param			X-Auth-Token	header		string	true	"X-Auth-Token"
// @Failure		401				{object}	nil
// @Param			X-Company-ID	header		string	true	"X-Company-ID"
// @Param			id				path		string	true	"Company ID"
// @Param			logo			formData	file	false	"Logo image"
// @Param			banner			formData	file	false	"Banner image"
// @Param			favicon			formData	file	false	"Favicon image"
// @Param			background		formData	file	false	"Background image"
// @Success		200				{object}	dJSON.Images
// @Failure		400				{object}	DTO.ErrorResponse
// @Router			/company/{id}/design/images [patch]
func UpdateCompanyImages(c *fiber.Ctx) error {
	img_types_allowed := map[string]bool{"logo": true, "banner": true, "favicon": true, "background": true}

	var company model.Company
	Design, err := UpdateImagesById(c, company.TableName(), &company, img_types_allowed)
	if err != nil {
		return err
	}

	return lib.ResponseFactory(c).SendDTO(200, &Design.Images, &dJSON.Images{})
}

// @Summary		Delete a specific company design image
// @Description	Delete logo, banner, favicon or background
// @Tags			Company
// @Security		ApiKeyAuth
// @Param			X-Auth-Token	header		string	true	"X-Auth-Token"
// @Failure		401				{object}	nil
// @Param			X-Company-ID	header		string	true	"X-Company-ID"
// @Param			id				path		string	true	"Company ID"
// @Param			image_type		path		string	true	"Type of image to delete (logo, banner, favicon, background)"
// @Success		200				{object}	dJSON.Design
// @Failure		400				{object}	DTO.ErrorResponse
// @Router			/company/{id}/design/images/{image_type} [delete]
func DeleteCompanyImage(c *fiber.Ctx) error {
	img_types_allowed := map[string]bool{"logo": true, "banner": true, "favicon": true, "background": true}
	var company model.Company
	Design, err := DeleteImageById(c, company.TableName(), &company, img_types_allowed)
	if err != nil {
		return err
	}

	return lib.ResponseFactory(c).SendDTO(200, &Design.Images, &dJSON.Images{})
}

// UpdateCompanyColors updates the colors of a company
//
//	@Summary		Update company colors
//	@Description	Update the primary, secondary, tertiary, and quaternary colors of a company
//	@Tags			Company
//	@Security		ApiKeyAuth
//	@Param			X-Auth-Token	header		string	true	"X-Auth-Token"
//	@Failure		401				{object}	nil
//	@Param			X-Company-ID	header		string	true	"X-Company-ID"
//	@Param			id				path		string	true	"Company ID"
//	@Accept			json
//	@Produce		json
//	@Param			colors	body		mJSON.Colors	true	"Colors"
//	@Success		200		{object}	dJSON.Colors
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/company/{id}/design/colors [put]
func UpdateCompanyColors(c *fiber.Ctx) error {
	tx, err := lib.Session(c)
	if err != nil {
		return err
	}

	var company model.Company
	id := c.Params("id")
	if err := tx.First(&company, "id = ?", id).Error; err != nil {
		return lib.Error.Company.NotFound.WithError(err)
	}

	var colors mJSON.Colors
	if err := c.BodyParser(&colors); err != nil {
		return lib.Error.General.BadRequest.WithError(err)
	}

	ValidateHexColor := func(color string) error {
		if color == "" {
			return nil // Empty color is allowed
		}
		if len(color) != 7 || color[0] != '#' {
			return fmt.Errorf("invalid hex color: %s", color)
		}
		for _, c := range color[1:] {
			if (c < '0' || c > '9') && (c < 'A' || c > 'F') && (c < 'a' || c > 'f') {
				return fmt.Errorf("invalid hex color: %s", color)
			}
		}
		return nil
	}

	if err := ValidateHexColor(colors.Primary); err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid primary color: %s", colors.Primary))
	}
	company.Design.Colors.Primary = colors.Primary
	if err := ValidateHexColor(colors.Secondary); err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid secondary color: %s", colors.Secondary))
	}
	company.Design.Colors.Secondary = colors.Secondary
	if err := ValidateHexColor(colors.Tertiary); err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid tertiary color: %s", colors.Tertiary))
	}
	company.Design.Colors.Tertiary = colors.Tertiary
	if err := ValidateHexColor(colors.Quaternary); err != nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid quaternary color: %s", colors.Quaternary))
	}
	company.Design.Colors.Quaternary = colors.Quaternary

	if err := tx.Save(&company).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return lib.ResponseFactory(c).SendDTO(200, &company.Design.Colors, &dJSON.Colors{})
}

// Constructor for company_controller
func Company(Gorm *handler.Gorm) {
	endpoint := &middleware.Endpoint{DB: Gorm}
	endpoint.BulkRegisterHandler([]fiber.Handler{
		CreateCompany,
		GetCompanyById,
		GetCompanyByName,
		CheckIfCompanyExistsByTaxID,
		GetCompanyByTaxId,
		GetCompanyBySubdomain,
		UpdateCompanyImages,
		UpdateCompanyColors,
		DeleteCompanyImage,
		UpdateCompanyById,
		DeleteCompanyById,
	})
}
