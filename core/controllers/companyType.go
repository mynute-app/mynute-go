package controllers

import (
	"agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v3"
)

type CompanyType struct {
	Request    *handlers.Request
	Middleware *middleware.CompanyType
}

func (ctc *CompanyType) getBy(paramKey string, c fiber.Ctx) error {
	var model []models.CompanyType
	var dto []DTO.CompanyType
	assocs := []string{}
	mdws := []func(fiber.Ctx) (int, error){}

	ctc.Request.GetBy(c, paramKey, &model, &dto, assocs, mdws)

	return nil
}

func (ctc *CompanyType) DeleteOneById(c fiber.Ctx) error {
	var model models.CompanyType
	mdws := []func(fiber.Ctx) (int, error){ctc.Middleware.DeleteOneById}

	ctc.Request.DeleteOneById(c, &model, mdws)

	return nil
}

func (ctc *CompanyType) UpdateOneById(c fiber.Ctx) error {
	var model models.CompanyType
	var dto DTO.CompanyType
	var changes map[string]interface{}
	assocs := []string{}
	mdws := []func(fiber.Ctx) (int, error){ctc.Middleware.Update}

	ctc.Request.UpdateOneById(c, &model, &dto, changes, assocs, mdws)

	return nil
}

func (ctc *CompanyType) CreateOne(c fiber.Ctx) error {
	var model models.CompanyType
	var dto DTO.CompanyType
	assocs := []string{}
	mdws := []func(fiber.Ctx) (int, error){ctc.Middleware.Create}

	ctc.Request.CreateOne(c, &model, &dto, assocs, mdws)

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
