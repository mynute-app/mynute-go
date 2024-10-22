package controllers

import (
	"agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v3"
)

type CompanyTypeController struct {
	BaseController[models.CompanyType, DTO.CompanyType]
}

func NewCompanyTypeController(Req *handlers.Request, Mid middleware.IMiddleware) *CompanyTypeController {
	return &CompanyTypeController{
		BaseController: BaseController[models.CompanyType, DTO.CompanyType]{
			Request:     Req,
			Middleware:  Mid,
			Associations: []string{"Companies"},
		},
	}
}

// Custom extension method to get a company type by name
func (cc *CompanyTypeController) GetOneByName(c fiber.Ctx) error {
	return cc.GetBy("name", c)
}