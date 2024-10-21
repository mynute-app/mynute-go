package controllers

import (
	"agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v3"
)

var _ IController = (*Company)(nil) // Ensures that the Struct implements the IController interface correctly

type Company struct {
	Request    *handlers.Request
	Middleware *middleware.Company
	Associations []string
}

func (cc *Company) getBy(paramKey string, c fiber.Ctx) error {
	var model []models.Company
	var dto []DTO.Company
	mdws := []func(fiber.Ctx) (int, error){}
	cc.Request.GetBy(c, paramKey, &model, &dto, cc.Associations, mdws)
	return nil
}

func (cc *Company) forceGetBy(paramKey string, c fiber.Ctx) error {
	var model []models.Company
	var dto []DTO.Company
	mdws := []func(fiber.Ctx) (int, error){}
	cc.Request.ForceGetBy(c, paramKey, &model, &dto, cc.Associations, mdws)
	return nil
}

func (cc *Company) UpdateOneById(c fiber.Ctx) error {
	var model models.Company
	var dto DTO.Company
	var changes map[string]interface{}
	mdws := []func(fiber.Ctx) (int, error){cc.Middleware.Update}
	cc.Request.UpdateOneById(c, &model, &dto, changes, cc.Associations, mdws)
	return nil
}

func (cc *Company) DeleteOneById(c fiber.Ctx) error {
	var model models.Company
	mdws := []func(fiber.Ctx) (int, error){}
	cc.Request.DeleteOneById(c, &model, mdws)
	return nil
}

func (cc *Company) ForceDeleteOneById(c fiber.Ctx) error {
	var model models.Company
	mdws := []func(fiber.Ctx) (int, error){}
	cc.Request.ForceDeleteOneById(c, &model, mdws)
	return nil
}

func (cc *Company) CreateOne(c fiber.Ctx) error {
	var model models.Company
	var dto DTO.Company
	mdws := []func(fiber.Ctx) (int, error){cc.Middleware.Create}
	cc.Request.CreateOne(c, &model, &dto, cc.Associations, mdws)
	return nil
}

func (cc *Company) GetAll(c fiber.Ctx) error {
	return cc.getBy("", c)
}

func (cc *Company) ForceGetAll(c fiber.Ctx) error {
	return cc.forceGetBy("", c)
}

func (cc *Company) ForceGetOneById(c fiber.Ctx) error {
	return cc.forceGetBy("id", c)
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
