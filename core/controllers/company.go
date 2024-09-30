package controllers

import (
	"agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/services"

	"github.com/gofiber/fiber/v3"
)

type Company struct {
	App *fiber.App
	DB  *services.Postgres
}

func (cc *Company) GetOneById(c fiber.Ctx) error {
	id := c.Params("id")
	var company models.Company
	if err := cc.DB.GetOneById(&company, id, []string{"CompanyTypes"}); err != nil {
		return lib.FiberError(404, c, err)
	}
	var companyDTO DTO.Company
	if err := services.ConvertToDTO(company, &companyDTO); err != nil {
		return lib.FiberError(500, c, err)
	}
	return c.JSON(companyDTO)
}

func (cc *Company) GetOneByName(c fiber.Ctx) error {
	name := c.Params("name")
	var company models.Company
	if err := cc.DB.GetOneByName(&company, name, []string{"CompanyTypes"}); err != nil {
		return lib.FiberError(404, c, err)
	}
	var companyDTO DTO.Company
	if err := services.ConvertToDTO(company, &companyDTO); err != nil {
		return lib.FiberError(500, c, err)
	}
	return c.JSON(companyDTO)
}

func (cc *Company) Create(c fiber.Ctx) error {
	var company models.Company
	if err := lib.BodyParser(c.Body(), &company); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := cc.DB.Create(&company, []string{"CompanyTypes"}); err != nil {
		return lib.FiberError(400, c, err)
	}
	var companyDTO DTO.Company
	if err := services.ConvertToDTO(company, &companyDTO); err != nil {
		return lib.FiberError(500, c, err)
	}
	return c.JSON(companyDTO)
}

func (cr *Company) UpdateById(c fiber.Ctx) error {
	id := c.Params("id")
	var company models.Company
	if err := cr.DB.GetOneById(&company, id, []string{"CompanyTypes"}); err != nil {
		return lib.FiberError(404, c, err)
	}
	if err := lib.BodyParser(c.Body(), &company); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := cr.DB.Update(&company); err != nil {
		return lib.FiberError(400, c, err)
	}
	var companyDTO DTO.Company
	if err := services.ConvertToDTO(company, &companyDTO); err != nil {
		return lib.FiberError(500, c, err)
	}
	return c.JSON(companyDTO)
}

func (cr *Company) DeleteById(c fiber.Ctx) error {
	id := c.Params("id")
	var company models.Company
	if err := cr.DB.GetOneById(&company, id, []string{"CompanyTypes"}); err != nil {
		return lib.FiberError(404, c, err)
	}
	if err := cr.DB.Delete(&company); err != nil {
		return lib.FiberError(400, c, err)
	}
	return c.SendStatus(204)
}
