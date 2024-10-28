package controllers

import (
	"agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v3"
)

// CompanyController embeds BaseController in order to extend it with the functions below
type companyController struct {
	BaseController[models.Company, DTO.Company]
}

// Custom extension method to get a company by name
func (cc *companyController) GetOneByName(c fiber.Ctx) error {
	return cc.GetBy("name", c)
}

// Custom extension method to get a company by tax ID
func (cc *companyController) GetOneByTaxId(c fiber.Ctx) error {
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
