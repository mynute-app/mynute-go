package controllers

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/services"

	"github.com/gofiber/fiber/v3"
)

type CompanyType struct {
	App *fiber.App
	DB  *services.Postgres
}

func (ctc *CompanyType) Create(c fiber.Ctx) error {
	var model models.CompanyType
	var dto DTO.CompanyType
	assocs := []string{}

	CtrlService := services.Controller{Ctx: c, DB: ctc.DB}
	if err := CtrlService.Create(&model, &dto, assocs); err != nil {
		return err
	}
	return nil
}

func (ctc *CompanyType) GetAll(c fiber.Ctx) error {
	var model []models.CompanyType
	var dto []DTO.CompanyType
	assocs := []string{}

	CtrlService := services.Controller{Ctx: c, DB: ctc.DB}
	if err := CtrlService.GetAll(&model, &dto, assocs); err != nil {
		return err
	}
	return nil
}

func (ctc *CompanyType) getBy(param string, c fiber.Ctx) error {
	var model models.CompanyType
	var dto DTO.CompanyType
	assocs := []string{}

	CtrlService := services.Controller{Ctx: c, DB: ctc.DB}
	if err := CtrlService.GetOneBy(param, &model, &dto, assocs); err != nil {
		return err
	}
	return nil
}

func (ctc *CompanyType) updateBy(param string, c fiber.Ctx) error {
	var model models.CompanyType
	var dto DTO.CompanyType
	assocs := []string{}

	CtrlService := services.Controller{Ctx: c, DB: ctc.DB}
	if err := CtrlService.UpdateOneBy(param, &model, &dto, assocs); err != nil {
		return err
	}
	return nil
}

func (ctc *CompanyType) deleteBy(param string, c fiber.Ctx) error {
	var model models.CompanyType

	CtrlService := services.Controller{Ctx: c, DB: ctc.DB}
	if err := CtrlService.DeleteOneBy(param, &model); err != nil {
		return err
	}
	return nil
}

func (ctc *CompanyType) GetOneById(c fiber.Ctx) error {
	return ctc.getBy("id", c)
}

func (ctc *CompanyType) GetOneByName(c fiber.Ctx) error {
	return ctc.getBy("name", c)
}

func (ctc *CompanyType) UpdateById(c fiber.Ctx) error {
	return ctc.updateBy("id", c)
}

func (ctc *CompanyType) DeleteById(c fiber.Ctx) error {
	return ctc.deleteBy("id", c)
}