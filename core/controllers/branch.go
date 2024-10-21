package controllers

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/middleware"
	"agenda-kaki-go/core/handlers"

	"github.com/gofiber/fiber/v3"
)

type Branch struct {
	Request      *handlers.Request
	Middleware   *middleware.Branch
	Associations []string
}

func (cb *Branch) getBy(paramKey string, c fiber.Ctx) error {
	var model []models.Branch
	var dto []DTO.Branch
	mdws := []func(fiber.Ctx) (int, error){cb.Middleware.CheckCompany}
	cb.Request.GetBy(c, paramKey, &model, &dto, cb.Associations, mdws)
	return nil
}

func (cb *Branch) DeleteOneById(c fiber.Ctx) error {
	var model models.Branch
	mdws := []func(fiber.Ctx) (int, error){cb.Middleware.CheckCompany}
	cb.Request.DeleteOneById(c, &model, mdws)
	return nil
}

func (cb *Branch) ForceDeleteOneById(c fiber.Ctx) error {
	var model models.Branch
	mdws := []func(fiber.Ctx) (int, error){cb.Middleware.CheckCompany}
	cb.Request.ForceDeleteOneById(c, &model, mdws)
	return nil
}

func (cb *Branch) UpdateOneById(c fiber.Ctx) error {
	var model models.Branch
	var dto DTO.Branch
	var changes map[string]interface{}
	mdws := []func(fiber.Ctx) (int, error){cb.Middleware.CheckCompany}
	cb.Request.UpdateOneById(c, &model, &dto, changes, cb.Associations, mdws)
	return nil
}

func (cb *Branch) CreateOne(c fiber.Ctx) error {
	var model models.Branch
	var dto DTO.Branch

	mdws := []func(fiber.Ctx) (int, error){
		cb.Middleware.CheckCompany,
		cb.Middleware.Create,
	}
	cb.Request.CreateOne(c, &model, &dto, cb.Associations, mdws)
	return nil
}

func (cb *Branch) GetAll(c fiber.Ctx) error {
	return cb.getBy("", c)
}

func (cb *Branch) GetOneById(c fiber.Ctx) error {
	return cb.getBy("id", c)
}

