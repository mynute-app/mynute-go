package controllers

import (
	"agenda-kaki-company-go/api/lib"
	"agenda-kaki-company-go/api/services"
	"agenda-kaki-company-go/api/models"

	"github.com/gofiber/fiber/v3"
)

type Company struct {
	App *fiber.App
	DB  *services.Postgres
}

func (cc *Company) GetOneById(c fiber.Ctx) error {
	id := c.Params("id")
	var company models.Company
	result := cc.DB.GetOneById(&company, id)
	if result.Error() != "" {
		return c.Status(404).JSON(fiber.Map{"error": result.Error()})
	}
	return c.JSON(company)
}

func (cc *Company) Create(c fiber.Ctx) error {
	var company models.Company
	if err := lib.BodyParser(c.Body(), &company); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	result := cc.DB.Create(&company)
	if result.Error() != "" {
		return c.Status(400).JSON(fiber.Map{"error": result.Error()})
	}
	return c.JSON(company)
}

func (cr *Company) UpdateById(c fiber.Ctx) error {
	id := c.Params("id")
	var company models.Company
	result := cr.DB.GetOneById(&company, id)
	if result.Error() != "" {
		return c.Status(404).JSON(fiber.Map{"error": result.Error()})
	}
	if err := lib.BodyParser(c.Body(), &company); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	result = cr.DB.Update(&company)
	if result.Error() != "" {
		return c.Status(400).JSON(fiber.Map{"error": result.Error()})
	}
	return c.JSON(company)
}

func (cr *Company) DeleteById(c fiber.Ctx) error {
	id := c.Params("id")
	var company models.Company
	result := cr.DB.GetOneById(&company, id)
	if result.Error() != "" {
		return c.Status(404).JSON(fiber.Map{"error": result.Error()})
	}
	result = cr.DB.Delete(&company)
	if result.Error() != "" {
		return c.Status(400).JSON(fiber.Map{"error": result.Error()})
	}
	return c.SendStatus(204)
}