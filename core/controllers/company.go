package controllers

import (
	"agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v3"
)

type Company struct {
	Gorm        *handlers.Gorm
	Middleware  *middleware.Company
	HttpHandler *handlers.HTTP
}

func (cc *Company) getBy(paramKey string, c fiber.Ctx) error {
	var model models.Company
	var dto DTO.Company
	var assocs = []string{"CompanyTypes"}

	cc.HttpHandler.
		FiberCtx(c).
		Model(&model).
		DTO(&dto).
		Assoc(assocs).
		GetBy(paramKey)

	return nil
}

func (cc *Company) UpdateOneById(c fiber.Ctx) error {
	var model models.Company
	var dto DTO.Company
	var assocs = []string{"CompanyTypes"}

	cc.HttpHandler.
		FiberCtx(c).
		Model(&model).
		DTO(&dto).
		Assoc(assocs).
		UpdateOneById()

	return nil
}

func (cc *Company) DeleteOneById(c fiber.Ctx) error {
	var model models.Company
	var dto DTO.Company
	var assocs = []string{"CompanyTypes"}

	cc.HttpHandler.
		FiberCtx(c).
		Model(&model).
		DTO(&dto).
		Assoc(assocs).
		DeleteOneById()

	return nil
}

func (cc *Company) Create(c fiber.Ctx) error {
	var model models.Company
	var dto DTO.Company
	var assocs = []string{"CompanyTypes"}

	cc.HttpHandler.
		FiberCtx(c).
		Model(&model).
		DTO(&dto).
		Assoc(assocs).
		Create()

	return nil
}

func (cc *Company) GetAll(c fiber.Ctx) error {
	return cc.getBy("", c)
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

// func (cr *Company) UpdateById(c fiber.Ctx) error {
// 	return cr.updateBy("id", c)
// }
