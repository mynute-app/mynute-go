package controllers

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v2"
)

type companyTypeController struct {
	BaseController[models.CompanyType, DTO.CompanyType]
}

func CompanyType(Gorm *handlers.Gorm) *companyTypeController {
	return &companyTypeController{
		BaseController: BaseController[models.CompanyType, DTO.CompanyType]{
			Name:         namespace.CompanyTypeKey.Name,
			Request:      handlers.Request(Gorm),
			Middleware:   middleware.CompanyType(Gorm),
			Associations: []string{},
		},
	}
}

// Custom extension method to get a company type by name
func (cc *companyTypeController) GetOneByName(c *fiber.Ctx) error {
	return cc.GetBy("name", c)
}
