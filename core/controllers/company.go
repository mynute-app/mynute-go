package controllers

import (
	"agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/middleware"
	"agenda-kaki-go/core/services"

	"log"

	"github.com/gofiber/fiber/v3"
)

type Company struct {
	DB    *services.Postgres
	Middleware *middleware.Company
}

// func (cr *Company) updateBy(param string, c fiber.Ctx) error {
// 	var model models.Company
// 	var dto DTO.Company
// 	assocs := []string{"CompanyTypes"}
// 	CtrlService := services.Controller{Ctx: c, DB: cr.DB}
// 	if err := CtrlService.UpdateOneBy(param, &model, &dto, assocs); err != nil {
// 		log.Printf("An internal error occurred! %v", err)
// 		return err
// 	}
// 	return nil
// }

func (cc *Company) updateBy(paramKey string, c fiber.Ctx) error {
	var changes map[string]interface{}

	if err := lib.BodyParser(c.Body(), &changes); err != nil {
		return lib.FiberError(400, c, err)
	}

	if err := cc.Middleware.Update(changes); err != nil {
		return lib.FiberError(400, c, err)
	}

	var model models.Company

	assocs := []string{"CompanyTypes"}
	paramVal := c.Params(paramKey)

	if err := cc.DB.UpdateOneBy(paramKey, paramVal, &model, changes, assocs); err != nil {
		return lib.FiberError(400, c, err)
	}
	
	var dto DTO.Company

	if err := lib.ParseToDTO(model, &dto); err != nil {
		return lib.FiberError(500, c, err)
	}

	if err := c.JSON(dto); err != nil {
		log.Printf("An internal error occurred! %v", err)
		return err
	}

	return nil
}

func (cc *Company) Create(c fiber.Ctx) error {
	var model models.Company

	if err := lib.BodyParser(c.Body(), &model); err != nil {
		return lib.FiberError(400, c, err)
	}

	if err := cc.Middleware.Create(model); err != nil {
		return lib.FiberError(400, c, err)
	}

	assocs := []string{"CompanyTypes"}

	if err := cc.DB.Create(&model, assocs); err != nil {
		return lib.FiberError(400, c, err)
	}

	var dto DTO.Company

	if err := lib.ParseToDTO(model, &dto); err != nil {
		return lib.FiberError(500, c, err)
	}

	if err := c.JSON(dto); err != nil {
		log.Printf("An internal error occurred! %v", err)
		return err
	}

	return nil
}

func (cc *Company) GetAll(c fiber.Ctx) error {
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

func (cc *Company) getBy(param string, c fiber.Ctx) error {
	var model models.Company
	var dto DTO.Company
	assocs := []string{"CompanyTypes"}
	CtrlService := services.Controller{Ctx: c, DB: cc.DB}
	if err := CtrlService.GetOneBy(param, &model, &dto, assocs); err != nil {
		log.Printf("An internal error occurred! %v", err)
		return err
	}
	return nil
}

func (cr *Company) deleteBy(param string, c fiber.Ctx) error {
	var model models.Company
	CtrlService := services.Controller{Ctx: c, DB: cr.DB}
	if err := CtrlService.DeleteOneBy(param, &model); err != nil {
		log.Printf("An internal error occurred! %v", err)
		return err
	}
	return nil
}

func (cc *Company) GetOneById(c fiber.Ctx) error {
	return cc.getBy("id", c)
}

func (cc *Company) GetOneByName(c fiber.Ctx) error {
	return cc.getBy("name", c)
}

func (cc *Company) GetOneByTaxId(c fiber.Ctx) error {
	return cc.getBy("tax_id", c)
}

func (cr *Company) UpdateById(c fiber.Ctx) error {
	return cr.updateBy("id", c)
}

func (cr *Company) UpdateByName(c fiber.Ctx) error {
	return cr.updateBy("name", c)
}

func (cr *Company) UpdateByTaxId(c fiber.Ctx) error {
	return cr.updateBy("tax_id", c)
}

func (cr *Company) DeleteById(c fiber.Ctx) error {
	return cr.deleteBy("id", c)
}

func (cr *Company) DeleteByName(c fiber.Ctx) error {
	return cr.deleteBy("name", c)
}

func (cr *Company) DeleteByTaxId(c fiber.Ctx) error {
	return cr.deleteBy("tax_id", c)
}
