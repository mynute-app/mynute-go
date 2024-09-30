package controllers

import (
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/services"
	"log"

	"github.com/gofiber/fiber/v3"
)

type CompanyType struct {
	App *fiber.App
	DB  *services.Postgres
}

func (ctc *CompanyType) GetOneById(c fiber.Ctx) error {
	id := c.Params("id")
	var companyType models.CompanyType
	if err := ctc.DB.GetOneById(&companyType, id, nil); err != nil {
		return lib.FiberError(404, c, err)
	}
	return c.JSON(companyType)
}

func (ctc *CompanyType) Create(c fiber.Ctx) error {
	var companyType models.CompanyType
	if err := lib.BodyParser(c.Body(), &companyType); err != nil {
		return lib.FiberError(400, c, err)
	}
	log.Printf("CompanyType: %+v", companyType)
	if err := ctc.DB.Create(&companyType, nil); err != nil {
		return lib.FiberError(400, c, err)
	}
	return c.JSON(companyType)
}
