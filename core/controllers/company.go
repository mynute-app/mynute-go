package controllers

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v2"
)

// CompanyController embeds BaseController in order to extend it with the functions below
type companyController struct {
	BaseController[models.Company, DTO.Company]
}

// GetOneByName retrieves a company by name
// @Summary Get company by name
// @Description Retrieve a company by its name
// @Tags Company
// @Param name path string true "Company Name"
// @Produce json
// @Success 200 {object} DTO.Company
// @Failure 404 {object} DTO.ErrorResponse
// @Router /company/name/{name} [get]
func (cc *companyController) GetOneByName(c *fiber.Ctx) error {
	return cc.GetBy("name", c)
}

// GetOneByTaxId retrieves a company by tax ID
// @Summary Get company by tax ID
// @Description Retrieve a company by its tax identification number
// @Tags Company
// @Param tax_id path string true "Company Tax ID"
// @Produce json
// @Success 200 {object} DTO.Company
// @Failure 404 {object} DTO.ErrorResponse
// @Router /company/tax_id/{tax_id} [get]
func (cc *companyController) GetOneByTaxId(c *fiber.Ctx) error {
	return cc.GetBy("tax_id", c)
}

// Constructor for CompanyController
func Company(Gorm *handlers.Gorm) *companyController {
	return &companyController{
		BaseController: BaseController[models.Company, DTO.Company]{
			Name:         namespace.CompanyKey.Name,
			Request:      handlers.Request(Gorm),
			Middleware:   middleware.Company(Gorm),
			Associations: []string{"CompanyTypes", "Branches", "Employees", "Services"},
		},
	}
}
