package controllers

import (
	"agenda-kaki-company-go/api/lib"
	"agenda-kaki-company-go/api/models"
	"agenda-kaki-company-go/api/services"

	"github.com/gofiber/fiber/v3"
)

type Company struct {
	App *fiber.App
	DB  *services.Postgres
}

func (cc *Company) GetOneById(c fiber.Ctx) error {
	id := c.Params("id")
	var company models.Company
	if err := cc.DB.GetOneById(&company, id); err != nil {
		return lib.FiberError(404, c, err)
	}
	return c.JSON(company)
}

func (cc *Company) Create(c fiber.Ctx) error {
	var company models.Company
	if err := lib.BodyParser(c.Body(), &company); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := cc.DB.Create(&company); err != nil {
		return lib.FiberError(400, c, err)
	}
	return c.JSON(company)
}

func (cr *Company) UpdateById(c fiber.Ctx) error {
	id := c.Params("id")
	var company models.Company
	if err := cr.DB.GetOneById(&company, id); err != nil {
		return lib.FiberError(404, c, err)
	}
	if err := lib.BodyParser(c.Body(), &company); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := cr.DB.Update(&company); err != nil {
		return lib.FiberError(400, c, err)
	}
	return c.JSON(company)
}

func (cr *Company) DeleteById(c fiber.Ctx) error {
	id := c.Params("id")
	var company models.Company
	if err := cr.DB.GetOneById(&company, id); err != nil {
		return lib.FiberError(404, c, err)
	}
	if err := cr.DB.Delete(&company); err != nil {
		return lib.FiberError(400, c, err)
	}
	return c.SendStatus(204)
}
