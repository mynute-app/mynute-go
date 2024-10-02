package controllers

import (
	"agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/services"

	"github.com/gofiber/fiber/v3"
	"log"
)

type Company struct {
	App   *fiber.App
	DB    *services.Postgres
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

func (cr *Company) updateBy(param string, c fiber.Ctx) error {
	var model models.Company
	var dto DTO.Company
	assocs := []string{"CompanyTypes"}
	CtrlService := services.Controller{Ctx: c, DB: cr.DB}
	if err := CtrlService.UpdateOneBy(param, &dto, model, assocs); err != nil {
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
