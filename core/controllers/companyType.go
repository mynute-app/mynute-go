package controllers

import (
	"agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v3"
)

type CompanyType struct {
	Gorm        *handlers.Gorm
	Middleware  *middleware.CompanyType
	HttpHandler *handlers.HTTP
}

func (ctc *CompanyType) getBy(paramKey string, c fiber.Ctx) error {

	assocs := []string{}

	if paramKey == "" {
		var modelArr []models.CompanyType
		var dtoArr []DTO.CompanyType
		ctc.HttpHandler.
			Model(&modelArr).
			DTO(&dtoArr).
			FiberCtx(c).
			Assoc(assocs).
			GetOneBy(paramKey)
		return nil
	}

	var model models.CompanyType
	var dto DTO.CompanyType
	
	ctc.HttpHandler.
		Model(&model).
		DTO(&dto).
		Assoc(assocs).
		FiberCtx(c).
		GetOneBy(paramKey)

	return nil
}

func (ctc *CompanyType) updateBy(paramKey string, c fiber.Ctx) error {
	var model models.CompanyType
	var dto DTO.CompanyType
	assocs := []string{}

	ctc.HttpHandler.
		Model(&model).
		DTO(&dto).
		Assoc(assocs).
		FiberCtx(c).
		UpdateOneBy(paramKey)

	return nil
}

func (ctc *CompanyType) deleteBy(paramKey string, c fiber.Ctx) error {
	var model models.CompanyType
	var dto DTO.CompanyType
	assocs := []string{}

	ctc.HttpHandler.
		Model(&model).
		DTO(&dto).
		Assoc(assocs).
		FiberCtx(c).
		Middleware(ctc.Middleware.Delete).
		DeleteOneBy(paramKey)

	return nil
}

func (ctc *CompanyType) Create(c fiber.Ctx) error {
	var model models.CompanyType
	var dto DTO.CompanyType
	assocs := []string{}

	ctc.HttpHandler.
		Model(&model).
		DTO(&dto).
		Assoc(assocs).
		FiberCtx(c).
		Create()

	return nil
}

func (ctc *CompanyType) GetAll(c fiber.Ctx) error {
	return ctc.getBy("", c)
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

func (ctc *CompanyType) UpdateByName(c fiber.Ctx) error {
	return ctc.updateBy("name", c)
}

func (ctc *CompanyType) DeleteById(c fiber.Ctx) error {
	return ctc.deleteBy("id", c)
}

func (ctc *CompanyType) DeleteByName(c fiber.Ctx) error {
	return ctc.deleteBy("name", c)
}
