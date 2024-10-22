package controllers

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/models"
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v3"
)

var _ IController = (*Service)(nil)

type Service struct {
	Request      *handlers.Request
	Middleware   *middleware.Service
	Associations []string
}

func (cs *Service) getBy(paramKey string, c fiber.Ctx) error {
	var model []models.Service
	var dto []DTO.Service
	mdws := []func(fiber.Ctx) (int, error){}
	cs.Request.GetBy(c, paramKey, &model, &dto, cs.Associations, mdws)
	return nil
}

func (cs *Service) forceGetBy(paramKey string, c fiber.Ctx) error {
	var model []models.Service
	var dto []DTO.Service
	mdws := []func(fiber.Ctx) (int, error){}
	cs.Request.ForceGetBy(c, paramKey, &model, &dto, cs.Associations, mdws)
	return nil
}

func (cs *Service) DeleteOneById(c fiber.Ctx) error {
	var model models.Service
	mdws := []func(fiber.Ctx) (int, error){}
	cs.Request.DeleteOneById(c, &model, mdws)
	return nil
}

func (cs *Service) ForceDeleteOneById(c fiber.Ctx) error {
	var model models.Service
	mdws := []func(fiber.Ctx) (int, error){}
	cs.Request.ForceDeleteOneById(c, &model, mdws)
	return nil
}

func (cs *Service) UpdateOneById(c fiber.Ctx) error {
	var model models.Service
	var dto DTO.Service
	var changes map[string]interface{}
	mdws := []func(fiber.Ctx) (int, error){}
	cs.Request.UpdateOneById(c, &model, &dto, changes, cs.Associations, mdws)
	return nil
}

func (cs *Service) CreateOne(c fiber.Ctx) error {
	var model models.Service
	var dto DTO.Service
	mdws := []func(fiber.Ctx) (int, error){cs.Middleware.Create}
	cs.Request.CreateOne(c, &model, &dto, cs.Associations, mdws)
	return nil
}

func (cs *Service) GetAll(c fiber.Ctx) error {
	return cs.getBy("", c)
}

func (cs *Service) GetOneById(c fiber.Ctx) error {
	return cs.getBy("id", c)
}

func (cs *Service) ForceGetAll(c fiber.Ctx) error {
	return cs.getBy("", c)
}

func (cs *Service) ForceGetOneById(c fiber.Ctx) error {
	return cs.getBy("id", c)
}