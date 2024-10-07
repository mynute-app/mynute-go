package controllers

import (
	"agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/config/namespace"
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
			FiberCtx(c).
			Model(&modelArr).
			DTO(&dtoArr).
			Assoc(assocs).
			GetOneBy(paramKey)
		return nil
	}

	var model models.CompanyType
	var dto DTO.CompanyType

	ctc.HttpHandler.
		FiberCtx(c).
		Model(&model).
		DTO(&dto).
		Assoc(assocs).
		GetOneBy(paramKey)

	return nil
}

func (ctc *CompanyType) UpdateOneById(c fiber.Ctx) error {
	var model models.CompanyType
	var dto DTO.CompanyType
	assocs := []string{}
	keys := namespace.GeneralKey
	modelParserCtx := middleware.ParseBodyToContext(keys.Model, &model)
	dtoCtx := middleware.AddToContext(keys.Dto, &dto)
	assocsCtx := middleware.AddToContext(keys.Associations, &assocs)

	ctc.HttpHandler.
		FiberCtx(c).
		Middleware(modelParserCtx).
		Middleware(dtoCtx).
		Middleware(assocsCtx).
		Middleware(ctc.Middleware.Update).
		UpdateOneById()

	return nil
}

func (ctc *CompanyType) DeleteOneById(c fiber.Ctx) error {
	var model models.CompanyType
	modelCtx := middleware.AddToContext(namespace.GeneralKey.Model, &model)
	ctc.HttpHandler.
		FiberCtx(c).
		Middleware(modelCtx).
		Middleware(ctc.Middleware.Delete).
		DeleteOneById()

	return nil
}

func (ctc *CompanyType) Create(c fiber.Ctx) error {
	var model models.CompanyType
	var dto DTO.CompanyType
	assocs := []string{}

	modelParserCtx := middleware.ParseBodyToContext(namespace.GeneralKey.Model, &model)
	dtoCtx := middleware.AddToContext(namespace.GeneralKey.Dto, &dto)
	assocsCtx := middleware.AddToContext(namespace.GeneralKey.Associations, &assocs)

	ctc.HttpHandler.
		FiberCtx(c).
		Middleware(modelParserCtx).
		Middleware(dtoCtx).
		Middleware(assocsCtx).
		Middleware(ctc.Middleware.Create).
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

// func (ctc *CompanyType) UpdateById(c fiber.Ctx) error {
// 	return ctc.updateBy("id", c)
// }

// func (ctc *CompanyType) UpdateByName(c fiber.Ctx) error {
// 	return ctc.updateBy("name", c)
// }
