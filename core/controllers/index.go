package controllers

import (
	"agenda-kaki-go/core/handlers"
	"agenda-kaki-go/core/middleware"

	"github.com/gofiber/fiber/v3"
)

var _ IController = (*BaseController[IController, IController])(nil)

type BaseController[MODEL any, DTO any] struct {
	Request      *handlers.Request
	Middleware   middleware.IMiddleware
	Associations []string
}

type IController interface {
	GetBy(paramKey string, c fiber.Ctx) error
	ForceGetBy(paramKey string, c fiber.Ctx) error
	CreateOne(c fiber.Ctx) error
	GetAll(c fiber.Ctx) error
	GetOneById(c fiber.Ctx) error
	UpdateOneById(c fiber.Ctx) error
	DeleteOneById(c fiber.Ctx) error
	ForceDeleteOneById(c fiber.Ctx) error
	ForceGetOneById(c fiber.Ctx) error
	ForceGetAll(c fiber.Ctx) error
}

func CreateRoutes(r fiber.Router, ci IController) {
	r.Post("/", ci.CreateOne)       // ok
	r.Get("/", ci.GetAll)           // ok
	r.Get("/force", ci.ForceGetAll) // ok
	id := r.Group("/:id")
	id.Get("/", ci.GetOneById)                 // ok
	id.Patch("/", ci.UpdateOneById)            // ok
	id.Delete("/", ci.DeleteOneById)           // ok
	id.Delete("/force", ci.ForceDeleteOneById) // ok
	id.Get("/force", ci.ForceGetOneById)       // ok
}

func (bc *BaseController[MODEL, DTO]) GetBy(paramKey string, c fiber.Ctx) error {
	var model []MODEL
	var dto []DTO
	bc.Request.GetBy(c, paramKey, &model, &dto, bc.Associations, bc.Middleware.GET())
	return nil
}

func (bc *BaseController[MODEL, DTO]) ForceGetBy(paramKey string, c fiber.Ctx) error {
	var model []MODEL
	var dto []DTO
	bc.Request.ForceGetBy(c, paramKey, &model, &dto, bc.Associations, bc.Middleware.ForceGET())
	return nil
}

func (bc *BaseController[MODEL, DTO]) DeleteOneById(c fiber.Ctx) error {
	var model MODEL
	bc.Request.DeleteOneById(c, &model, bc.Middleware.DELETE())
	return nil
}

func (bc *BaseController[MODEL, DTO]) ForceDeleteOneById(c fiber.Ctx) error {
	var model MODEL
	bc.Request.ForceDeleteOneById(c, &model, bc.Middleware.ForceDELETE())
	return nil
}

func (bc *BaseController[MODEL, DTO]) UpdateOneById(c fiber.Ctx) error {
	var model MODEL
	var dto DTO
	var changes map[string]interface{}
	bc.Request.UpdateOneById(c, &model, &dto, changes, bc.Associations, bc.Middleware.PATCH())
	return nil
}

func (bc *BaseController[MODEL, DTO]) CreateOne(c fiber.Ctx) error {
	var model MODEL
	var dto DTO
	bc.Request.CreateOne(c, &model, &dto, bc.Associations, bc.Middleware.POST())
	return nil
}

func (bc *BaseController[MODEL, DTO]) GetAll(c fiber.Ctx) error {
	return bc.GetBy("", c)
}

func (bc *BaseController[MODEL, DTO]) GetOneById(c fiber.Ctx) error {
	return bc.GetBy("id", c)
}

func (bc *BaseController[MODEL, DTO]) ForceGetOneById(c fiber.Ctx) error {
	return bc.ForceGetBy("id", c)
}

func (bc *BaseController[MODEL, DTO]) ForceGetAll(c fiber.Ctx) error {
	return bc.ForceGetBy("", c)
}
