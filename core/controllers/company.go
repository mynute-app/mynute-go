package controllers

import (
	"agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v3"
)

// CompanyController embeds BaseController in order to extend it with the functions below
type CompanyController struct {
	BaseController[models.Company, DTO.Company]
}

// Custom extension method to get a company by name
func (cc *CompanyController) GetOneByName(c fiber.Ctx) error {
	return cc.GetBy("name", c)
}

// Custom extension method to get a company by tax ID
func (cc *CompanyController) GetOneByTaxId(c fiber.Ctx) error {
	return cc.GetBy("tax_id", c)
}

// Constructor for CompanyController
func NewCompanyController(Req *handlers.Request, Mid middleware.IMiddleware) *CompanyController {
	return &CompanyController{
		BaseController: BaseController[models.Company, DTO.Company]{
			Request:     Req,
			Middleware:  Mid,
			Associations: []string{"CompanyTypes", "Branches"},
		},
	}
}


