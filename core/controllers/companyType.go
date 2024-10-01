package controllers

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/services"

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
	var companyTypeDTO DTO.CompanyType
	if err := services.ConvertToDTO(companyType, &companyTypeDTO); err != nil {
		return lib.FiberError(500, c, err)
	}
	return c.JSON(companyType)
}

func (ctc *CompanyType) GetOneByName(c fiber.Ctx) error {
	name := c.Params("name")
	var companyType models.CompanyType
	if err := ctc.DB.GetOneByName(&companyType, name, nil); err != nil {
		return lib.FiberError(404, c, err)
	}
	var companyTypeDTO DTO.CompanyType
	if err := services.ConvertToDTO(companyType, &companyTypeDTO); err != nil {
		return lib.FiberError(500, c, err)
	}
	return c.JSON(companyTypeDTO)
}

func (ctc *CompanyType) GetAll(c fiber.Ctx) error {
	var companyTypes []models.CompanyType
	if err := ctc.DB.GetAll(&companyTypes, nil); err != nil {
		return lib.FiberError(404, c, err)
	}
	var companyTypesDTO []DTO.CompanyType
	if err := services.ConvertToDTO(companyTypes, &companyTypesDTO); err != nil {
		return lib.FiberError(500, c, err)
	}
	return c.JSON(companyTypesDTO)
}

func (ctc *CompanyType) Create(c fiber.Ctx) error {
	var companyType models.CompanyType
	if err := lib.BodyParser(c.Body(), &companyType); err != nil {
		return lib.FiberError(400, c, err)
	}
	if err := ctc.DB.Create(&companyType, nil); err != nil {
		return lib.FiberError(400, c, err)
	}
	return c.JSON(companyType)
}
