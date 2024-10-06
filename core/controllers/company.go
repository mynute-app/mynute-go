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
		Model(&model).
		DTO(&dto).
		Assoc(assocs).
		FiberCtx(c).
		GetOneBy(paramKey)

	return nil
}

func (cc *Company) updateBy(paramKey string, c fiber.Ctx) error {
	var model models.Company
	var dto DTO.Company
	var assocs = []string{"CompanyTypes"}

	cc.HttpHandler.
		Model(&model).
		DTO(&dto).
		Assoc(assocs).
		FiberCtx(c).
		UpdateOneBy(paramKey)

	return nil
}

func (cc *Company) deleteBy(paramKey string, c fiber.Ctx) error {
	var model models.Company
	var dto DTO.Company
	var assocs = []string{"CompanyTypes"}

	cc.HttpHandler.
		Model(&model).
		DTO(&dto).
		Assoc(assocs).
		FiberCtx(c).
		DeleteOneBy(paramKey)

	return nil
}

func (cc *Company) Create(c fiber.Ctx) error {
	var model models.Company
	var dto DTO.Company
	var assocs = []string{"CompanyTypes"}

	cc.HttpHandler.
		Model(&model).
		DTO(&dto).
		Assoc(assocs).
		FiberCtx(c).
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
