package controllers

import (
	"agenda-kaki-company-go/api/lib"
	"agenda-kaki-company-go/api/models"
	"agenda-kaki-company-go/api/services"

	"github.com/gofiber/fiber/v3"
)

type CompanyType struct {
	App *fiber.App
	DB  *services.Postgres
}

func (ctc *CompanyType) GetOneById(c fiber.Ctx) error {
	id := c.Params("id")
	var companyType models.CompanyType
	result := ctc.DB.GetOneById(&companyType, id)
	if result.Error() != "" {
		return c.Status(404).JSON(fiber.Map{"error": result.Error()})
	}
	return c.JSON(companyType)
}

func (ctc *CompanyType) Create(c fiber.Ctx) error {
	var companyType models.CompanyType
	if err := lib.ParseBody(c.Request(), &companyType); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	result := ctc.DB.Create(&companyType)
	if result.Error() != "" {
		return c.Status(400).JSON(fiber.Map{"error": result.Error()})
	}
	return c.JSON(companyType)
}
