package controllers

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v3"
)

type Employee struct {
	Request    *handlers.Request
	Middleware *middleware.Employee
	Associations []string
}

func (ce *Employee) getBy(paramKey string, c fiber.Ctx) error {
	var model []models.Employee
	var dto []DTO.Employee
	mdws := []func(fiber.Ctx) (int, error){}
	ce.Request.GetBy(c, paramKey, &model, &dto, ce.Associations, mdws)
	return nil
}

func (ce *Employee) DeleteOneById(c fiber.Ctx) error {
	var model models.Employee
	mdws := []func(fiber.Ctx) (int, error){}
	ce.Request.DeleteOneById(c, &model, mdws)
	return nil
}

func (ce *Employee) ForceDeleteOneById(c fiber.Ctx) error {
	var model models.Employee
	mdws := []func(fiber.Ctx) (int, error){}
	ce.Request.ForceDeleteOneById(c, &model, mdws)
	return nil
}

func (ce *Employee) UpdateOneById(c fiber.Ctx) error {
	var model models.Employee
	var dto DTO.Employee
	var changes map[string]interface{}
	mdws := []func(fiber.Ctx) (int, error){}
	ce.Request.UpdateOneById(c, &model, &dto, changes, ce.Associations, mdws)
	return nil
}

func (ce *Employee) CreateOne(c fiber.Ctx) error {
	var model models.Employee
	var dto DTO.Employee
	mdws := []func(fiber.Ctx) (int, error){ce.Middleware.Create}
	ce.Request.CreateOne(c, &model, &dto, ce.Associations, mdws)
	return nil
}

func (ce *Employee) GetAll(c fiber.Ctx) error {
	return ce.getBy("", c)
}

func (ce *Employee) GetOneById(c fiber.Ctx) error {
	return ce.getBy("id", c)
}

func (ce *Employee) GetOneByEmail(c fiber.Ctx) error {
	return ce.getBy("email", c)
}
