package controllers

import (
	"agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/services"

	"github.com/gofiber/fiber/v3"
	"log"
)

type Company struct {
	App   *fiber.App
	DB    *services.Postgres
}

func (cc *Company) GetAll(c fiber.Ctx) error {
	// var companies []models.Company
	// if err := cc.DB.GetAll(&companies, []string{"CompanyTypes"}); err != nil {
	// 	return lib.FiberError(404, c, err)
	// }
	// var companiesDTO []DTO.Company
	// if err := services.ConvertToDTO(companies, &companiesDTO); err != nil {
	// 	return lib.FiberError(500, c, err)
	// }
	// return c.JSON(companiesDTO)
	var model []models.Company
	var dto []DTO.Company
	assocs := []string{"CompanyTypes"}
	CtrlService := services.Controller{Ctx: c, DB: cc.DB}
	if err := CtrlService.GetAll(&model, &dto, assocs); err != nil {
		log.Printf("An internal error occurred! %v", err)
		return err
	}
	return nil
}

func (cc *Company) GetOneById(c fiber.Ctx) error {
	// id := c.Params("id")
	// var company models.Company
	// if err := cc.DB.GetOneById(&company, id, []string{"CompanyTypes"}); err != nil {
	// 	return lib.FiberError(404, c, err)
	// }
	// var companyDTO DTO.Company
	// if err := services.ConvertToDTO(company, &companyDTO); err != nil {
	// 	return lib.FiberError(500, c, err)
	// }
	// return c.JSON(companyDTO)
	var model models.Company
	var dto DTO.Company
	assocs := []string{"CompanyTypes"}
	CtrlService := services.Controller{Ctx: c, DB: cc.DB}
	if err := CtrlService.GetOneById(&model, &dto, assocs); err != nil {
		log.Printf("An internal error occurred! %v", err)
		return err
	}
	return nil
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
	var model models.Company
	var dto DTO.Company
	assocs := []string{"CompanyTypes"}
	CtrlService := services.Controller{Ctx: c, DB: cc.DB}
	if err := CtrlService.Create(&model, &dto, assocs); err != nil {
		log.Printf("An internal error occurred! %v", err)
		return err
	}
	return nil
}

func (cr *Company) UpdateById(c fiber.Ctx) error {
	// id := c.Params("id")
	// var company models.Company
	// if err := cr.DB.GetOneById(&company, id, []string{"CompanyTypes"}); err != nil {
	// 	return lib.FiberError(404, c, err)
	// }
	// if err := lib.BodyParser(c.Body(), &company); err != nil {
	// 	return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	// }
	// if err := cr.DB.UpdateOne(&company); err != nil {
	// 	return lib.FiberError(400, c, err)
	// }
	// var companyDTO DTO.Company
	// if err := services.ConvertToDTO(company, &companyDTO); err != nil {
	// 	return lib.FiberError(500, c, err)
	// }
	// return c.JSON(companyDTO)
	var model models.Company
	var dto DTO.Company
	assocs := []string{"CompanyTypes"}
	CtrlService := services.Controller{Ctx: c, DB: cr.DB}
	if err := CtrlService.UpdateById(&model, &dto, assocs); err != nil {
		log.Printf("An internal error occurred! %v", err)
		return err
	}
	return nil
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
